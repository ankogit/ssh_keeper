package mcp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestDispatchInitializeAndTools(t *testing.T) {
	initBody := dispatch([]byte(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05"}}`))
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

	toolsBody := dispatch([]byte(`{"jsonrpc":"2.0","id":2,"method":"tools/list"}`))
	if !strings.Contains(string(toolsBody), `"connect"`) || !strings.Contains(string(toolsBody), `"list_connections"`) {
		t.Fatalf("список инструментов: %s", toolsBody)
	}

	if dispatch([]byte(`{"jsonrpc":"2.0","method":"notifications/initialized"}`)) != nil {
		t.Fatal("на уведомление не нужен ответ")
	}
}

func TestFrameRoundTrip(t *testing.T) {
	var buf bytes.Buffer
	payload := []byte(`{"jsonrpc":"2.0","id":1,"method":"ping"}`)
	if err := writeFrame(&buf, payload); err != nil {
		t.Fatal(err)
	}
	got, err := readFrame(bufio.NewReader(&buf))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(payload) {
		t.Fatalf("кадр %s", got)
	}
}
