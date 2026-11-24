package mcp

import (
	"strings"
	"testing"

	"ssh-keeper/internal/models"
)

func TestMatchConnection(t *testing.T) {
	connections := []models.Connection{
		{Name: "dev toto", Host: "10.0.0.2", User: "root", Port: 22, UseSSHKey: true},
		{Name: "prod", Host: "prod.example", User: "ubuntu", Port: 22},
		{Name: "dev logs", Host: "10.0.0.3", User: "root", Port: 2222, UseSSHKey: true},
	}

	got, err := MatchConnection(connections, "dev toto")
	if err != nil {
		t.Fatal(err)
	}
	if got.Host != "10.0.0.2" {
		t.Fatalf("хост %s", got.Host)
	}

	got, err = MatchConnection(connections, "DEV TOTO")
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "dev toto" {
		t.Fatalf("имя %s", got.Name)
	}

	_, err = MatchConnection(connections, "dev")
	if err == nil || !strings.Contains(err.Error(), "нескольким") {
		t.Fatalf("ожидалась неоднозначность, получено %v", err)
	}

	_, err = MatchConnection(connections, "missing")
	if err == nil {
		t.Fatal("ожидалась ошибка")
	}
}
