package mcp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestDispatchInitializeAndTools(t *testing.T) {
	initBody := dispatch([]byte(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05"}}`), Options{})
	var initMsg struct {
		Result struct {
			ProtocolVersion string `json:"protocolVersion"`
			Instructions    string `json:"instructions"`
		} `json:"result"`
	}
	if err := json.Unmarshal(initBody, &initMsg); err != nil {
		t.Fatal(err)
	}
	if initMsg.Result.ProtocolVersion != "2024-11-05" {
		t.Fatalf("протокол %s", initMsg.Result.ProtocolVersion)
	}
	if !strings.Contains(initMsg.Result.Instructions, "list_connections") {
		t.Fatalf("нет инструкции: %s", initMsg.Result.Instructions)
	}

	toolsBody := dispatch([]byte(`{"jsonrpc":"2.0","id":2,"method":"tools/list"}`), Options{})
	if !strings.Contains(string(toolsBody), `"connect"`) || !strings.Contains(string(toolsBody), `"list_connections"`) {
		t.Fatalf("список инструментов: %s", toolsBody)
	}

	if dispatch([]byte(`{"jsonrpc":"2.0","method":"notifications/initialized"}`), Options{}) != nil {
		t.Fatal("на уведомление не нужен ответ")
	}
}

func TestNewlineDelimitedRoundTrip(t *testing.T) {
	payload := []byte(`{"jsonrpc":"2.0","id":1,"method":"ping"}`)
	var buf bytes.Buffer
	if err := writeMessage(&buf, payload); err != nil {
		t.Fatal(err)
	}
	got, err := readMessage(bufio.NewReader(&buf))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(payload) {
		t.Fatalf("сообщение %s", got)
	}
}

func TestContentLengthStillSupported(t *testing.T) {
	payload := []byte(`{"jsonrpc":"2.0","id":1,"method":"ping"}`)
	var raw bytes.Buffer
	raw.WriteString("Content-Length: 40\r\n\r\n")
	raw.Write(payload)
	got, err := readMessage(bufio.NewReader(&raw))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(payload) {
		t.Fatalf("сообщение %s", got)
	}
}
