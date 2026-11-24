package mcp

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"ssh-keeper/internal/services"
	"ssh-keeper/internal/ssh"
	"ssh-keeper/internal/version"
)

// Options настраивает MCP-сервер.
type Options struct {
	// EnsureReady готовит сервисы перед tools/call. Handshake на него не ждёт.
	EnsureReady func() error
}

// Serve читает MCP-сообщения из stdin и отвечает в stdout.
func Serve(opts Options) error {
	reader := bufio.NewReader(os.Stdin)
	for {
		body, err := readFrame(reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		response := dispatch(body, opts)
		if response == nil {
			continue
		}
		if err := writeFrame(os.Stdout, response); err != nil {
			return err
		}
	}
}

type rpcMessage struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func dispatch(body []byte, opts Options) []byte {
	var msg rpcMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		return mustJSON(errorResponse(nil, -32700, "parse error"))
	}
	if len(msg.ID) == 0 || string(msg.ID) == "null" {
		return nil
	}

	switch msg.Method {
	case "initialize":
		return mustJSON(resultResponse(msg.ID, initializeResult(msg.Params)))
	case "ping":
		return mustJSON(resultResponse(msg.ID, map[string]any{}))
	case "tools/list":
		return mustJSON(resultResponse(msg.ID, map[string]any{"tools": toolDefinitions()}))
	case "tools/call":
		if opts.EnsureReady != nil {
			if err := opts.EnsureReady(); err != nil {
				return mustJSON(resultResponse(msg.ID, toolError("инициализация: "+err.Error())))
			}
		}
		return mustJSON(resultResponse(msg.ID, callTool(msg.Params)))
	default:
		return mustJSON(errorResponse(msg.ID, -32601, "method not found"))
	}
}

func initializeResult(params json.RawMessage) map[string]any {
	protocol := "2024-11-05"
	var requested struct {
		ProtocolVersion string `json:"protocolVersion"`
	}
	if len(params) > 0 {
		_ = json.Unmarshal(params, &requested)
	}
	switch requested.ProtocolVersion {
	case "2024-11-05", "2025-03-26", "2025-06-18":
		protocol = requested.ProtocolVersion
	}
	return map[string]any{
		"protocolVersion": protocol,
		"capabilities": map[string]any{
			"tools": map[string]any{},
		},
		"serverInfo": map[string]any{
			"name":    "ssh-keeper",
			"version": version.GetVersion(),
		},
		"instructions": "Сначала вызови list_connections и найди сервер по имени. Затем connect выполнит команду на этом сервере через сохранённые ключ или пароль.",
	}
}

func toolDefinitions() []map[string]any {
	return []map[string]any{
		{
			"name":        "list_connections",
			"description": "Список сохранённых SSH-подключений: имя, пользователь, хост, порт и способ входа. Пароли не возвращаются. Вызови это, чтобы понять, какой сервер имеется в виду.",
			"inputSchema": map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
		},
		{
			"name":        "connect",
			"description": "Подключиться к сохранённому серверу и выполнить команду. Имя можно передать частично, например «dev toto». Если ключ хоста новый или изменился, вызов вернёт отпечаток и ничего не выполнит, пока пользователь не подтвердит trust_host.",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"connection": map[string]any{
						"type":        "string",
						"description": "Имя подключения или его часть",
					},
					"command": map[string]any{
						"type":        "string",
						"description": "Команда, которую нужно выполнить на сервере",
					},
					"trust_host": map[string]any{
						"type":        "boolean",
						"description": "Записать новый или изменившийся ключ хоста. Ставь true только после согласия пользователя.",
					},
					"timeout_seconds": map[string]any{
						"type":        "integer",
						"description": "Лимит ожидания в секундах, от 1 до 120. По умолчанию 30.",
					},
				},
				"required": []string{"connection", "command"},
			},
		},
	}
}

func callTool(params json.RawMessage) map[string]any {
	var call struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(params, &call); err != nil {
		return toolError("не удалось прочитать аргументы: " + err.Error())
	}

	switch call.Name {
	case "list_connections":
		if err := services.ReloadConnections(); err != nil {
			return toolError(err.Error())
		}
		return toolText(FormatConnections(services.GetConnections()))
	case "connect":
		return connectTool(call.Arguments)
	default:
		return toolError("неизвестный инструмент " + call.Name)
	}
}

func connectTool(raw json.RawMessage) map[string]any {
	var args struct {
		Connection string `json:"connection"`
		Command    string `json:"command"`
		TrustHost  bool   `json:"trust_host"`
		TimeoutSec int    `json:"timeout_seconds"`
	}
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &args); err != nil {
			return toolError("не удалось прочитать аргументы: " + err.Error())
		}
	}
	if err := services.ReloadConnections(); err != nil {
		return toolError(err.Error())
	}
	conn, err := MatchConnection(services.GetConnections(), args.Connection)
	if err != nil {
		return toolError(err.Error())
	}

	timeout := 30 * time.Second
	if args.TimeoutSec > 0 {
		timeout = time.Duration(args.TimeoutSec) * time.Second
	}
	output, err := ssh.RunCommand(&conn, args.Command, timeout, args.TrustHost)
	if err != nil {
		var keyErr *ssh.HostKeyError
		if errors.As(err, &keyErr) {
			return toolError(hostKeyMessage(conn.Name, keyErr))
		}
		if strings.TrimSpace(output) == "" {
			return toolError(err.Error())
		}
		return toolError(err.Error() + "\n" + output)
	}
	if strings.TrimSpace(output) == "" {
		output = "Команда выполнена, вывод пуст."
	}
	return toolText(output)
}

func hostKeyMessage(name string, err *ssh.HostKeyError) string {
	status := "новый"
	if err.Result.Status == ssh.StatusChanged {
		status = "изменился"
	}
	stored := err.Result.StoredFingerprint
	if stored == "" {
		stored = "нет"
	}
	return fmt.Sprintf("Ключ сервера %q %s. Было: %s. Стало: %s. Повтори connect с trust_host=true только если пользователь согласен записать этот ключ.",
		name, status, stored, err.Result.PresentedFingerprint)
}

func toolText(text string) map[string]any {
	return map[string]any{
		"content": []map[string]any{{"type": "text", "text": text}},
		"isError": false,
	}
}

func toolError(text string) map[string]any {
	return map[string]any{
		"content": []map[string]any{{"type": "text", "text": text}},
		"isError": true,
	}
}

func resultResponse(id json.RawMessage, result any) map[string]any {
	return map[string]any{"jsonrpc": "2.0", "id": json.RawMessage(id), "result": result}
}

func errorResponse(id json.RawMessage, code int, message string) map[string]any {
	response := map[string]any{
		"jsonrpc": "2.0",
		"error":   rpcError{Code: code, Message: message},
	}
	if len(id) > 0 {
		response["id"] = json.RawMessage(id)
	} else {
		response["id"] = nil
	}
	return response
}

func mustJSON(value any) []byte {
	body, err := json.Marshal(value)
	if err != nil {
		body = []byte(`{"jsonrpc":"2.0","id":null,"error":{"code":-32603,"message":"encode error"}}`)
	}
	return body
}

func readFrame(reader *bufio.Reader) ([]byte, error) {
	length := -1
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break
		}
		if len(line) >= 15 && strings.EqualFold(line[:15], "content-length:") {
			var n int
			if _, err := fmt.Sscanf(line[15:], "%d", &n); err == nil {
				length = n
			}
		}
	}
	if length < 0 {
		return nil, fmt.Errorf("в сообщении нет Content-Length")
	}
	body := make([]byte, length)
	if _, err := io.ReadFull(reader, body); err != nil {
		return nil, err
	}
	return body, nil
}

func writeFrame(writer io.Writer, body []byte) error {
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(body))
	if _, err := io.WriteString(writer, header); err != nil {
		return err
	}
	if _, err := writer.Write(body); err != nil {
		return err
	}
	// На пайпе Sync может вернуть EINVAL — это не ошибка протокола.
	if f, ok := writer.(*os.File); ok {
		_ = f.Sync()
	}
	return nil
}
