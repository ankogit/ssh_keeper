package ssh

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"

	"ssh-keeper/internal/models"
)

func TestBuildArgsUsesKeyWhenRequested(t *testing.T) {
	conn := &models.Connection{
		Host:      "example.com",
		Port:      2222,
		User:      "root",
		UseSSHKey: true,
		KeyPath:   "/tmp/id_ed25519",
		Password:  "should-not-appear",
	}

	args, err := BuildArgs(conn, "/tmp/known_hosts")
	if err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(args, " ")
	for _, want := range []string{
		"StrictHostKeyChecking=yes",
		"UserKnownHostsFile=/tmp/known_hosts",
		"-i /tmp/id_ed25519",
		"IdentitiesOnly=yes",
		"PubkeyAuthentication=yes",
		"PasswordAuthentication=no",
		"-p 2222",
		"root@example.com",
	} {
		if !strings.Contains(joined, want) {
			t.Fatalf("в аргументах нет %q: %s", want, joined)
		}
	}
	if strings.Contains(joined, "should-not-appear") {
		t.Fatalf("пароль попал в аргументы: %s", joined)
	}
}

func TestBuildArgsPasswordAuthWithoutStoredSecret(t *testing.T) {
	conn := &models.Connection{
		Host:      "example.com",
		Port:      22,
		User:      "ubuntu",
		UseSSHKey: false,
	}

	args, err := BuildArgs(conn, "/tmp/known_hosts")
	if err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(args, " ")
	if strings.Contains(joined, "IdentitiesOnly") || strings.Contains(joined, "-i ") {
		t.Fatalf("парольный вход не должен выбирать ключ: %s", joined)
	}
	if !strings.Contains(joined, "PasswordAuthentication=yes") {
		t.Fatalf("нет парольной аутентификации: %s", joined)
	}
	if strings.Contains(joined, "-p ") {
		t.Fatalf("порт 22 не нужно передавать: %s", joined)
	}
}

func TestStartCommandKeepsOptionsBeforeRemoteCommand(t *testing.T) {
	if _, err := exec.LookPath("ssh"); err != nil {
		t.Skip("ssh не установлен")
	}
	t.Setenv("HOME", t.TempDir())

	conn := &models.Connection{
		Host:      "example.com",
		User:      "root",
		Port:      22,
		UseSSHKey: true,
	}
	cmd, passFile, err := startCommand(context.Background(), conn, "uname -a")
	if passFile != "" {
		defer os.Remove(passFile)
	}
	if err != nil {
		t.Fatal(err)
	}
	if len(cmd.Args) < 2 {
		t.Fatalf("аргументы %v", cmd.Args)
	}
	if cmd.Args[len(cmd.Args)-1] != "uname -a" || cmd.Args[len(cmd.Args)-2] != "root@example.com" {
		t.Fatalf("команда должна идти после адреса: %v", cmd.Args)
	}
	joined := strings.Join(cmd.Args, " ")
	if !strings.Contains(joined, "ConnectTimeout=10") {
		t.Fatalf("нет таймаута: %s", joined)
	}
	if strings.Index(joined, "ConnectTimeout=10") > strings.Index(joined, "root@example.com") {
		t.Fatalf("опция оказалась после адреса: %s", joined)
	}
}

func TestPrepareStoredPasswordUsesAskpassFile(t *testing.T) {
	if _, err := exec.LookPath("ssh"); err != nil {
		t.Skip("ssh не установлен")
	}
	t.Setenv("HOME", t.TempDir())

	conn := &models.Connection{
		Host:      "example.com",
		Port:      22,
		User:      "root",
		UseSSHKey: false,
		Password:  "s3cret-value",
	}

	cmd, passFile, err := Prepare(conn)
	if err != nil {
		t.Fatal(err)
	}
	if passFile == "" {
		t.Fatal("ожидался файл пароля")
	}
	defer os.Remove(passFile)

	for _, arg := range cmd.Args {
		if strings.Contains(arg, "s3cret-value") {
			t.Fatalf("пароль попал в аргументы: %v", cmd.Args)
		}
	}

	info, err := os.Stat(passFile)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("права файла %v, ожидались 0600", info.Mode().Perm())
	}

	body, err := os.ReadFile(passFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "s3cret-value\n" {
		t.Fatalf("содержимое файла %q", body)
	}

	env := strings.Join(cmd.Env, "\n")
	if !strings.Contains(env, AskpassEnv+"=1") || !strings.Contains(env, "SSH_ASKPASS_REQUIRE=force") {
		t.Fatalf("нет окружения askpass: %s", env)
	}
	if !strings.Contains(env, PassFileEnv+"="+passFile) {
		t.Fatalf("нет пути к файлу пароля в окружении")
	}
}
