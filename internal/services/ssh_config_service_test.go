package services

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadConfigKeepsIDAndTimestamps(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config")
	created := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	updated := created.Add(time.Hour)
	content := "Host example\n" +
		"    HostName example.com\n" +
		"    User root\n" +
		"    ID abc-123\n" +
		"    CreatedAt " + created.Format(time.RFC3339) + "\n" +
		"    UpdatedAt " + updated.Format(time.RFC3339) + "\n"

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	config, err := NewSSHConfigService(path).LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if len(config.Hosts) != 1 {
		t.Fatalf("хостов %d, ожидался 1", len(config.Hosts))
	}
	host := config.Hosts[0]
	if host.ID != "abc-123" {
		t.Fatalf("ID %q, ожидался abc-123", host.ID)
	}
	if !host.CreatedAt.Equal(created) {
		t.Fatalf("CreatedAt %s, ожидался %s", host.CreatedAt, created)
	}
	if !host.UpdatedAt.Equal(updated) {
		t.Fatalf("UpdatedAt %s, ожидался %s", host.UpdatedAt, updated)
	}
}
