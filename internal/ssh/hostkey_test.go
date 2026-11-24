package ssh

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"path/filepath"
	"testing"

	"golang.org/x/crypto/ssh"
)

func testPublicKey(t *testing.T) ssh.PublicKey {
	t.Helper()
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	signer, err := ssh.NewSignerFromKey(priv)
	if err != nil {
		t.Fatal(err)
	}
	return signer.PublicKey()
}

func TestClassifyAndSaveKey(t *testing.T) {
	path := filepath.Join(t.TempDir(), "known_hosts")
	host := "example.com"
	port := 22
	first := testPublicKey(t)
	second := testPublicKey(t)

	status, stored, err := Classify(path, host, port, first)
	if err != nil {
		t.Fatal(err)
	}
	if status != StatusUnknown {
		t.Fatalf("новый хост: статус %v, ожидался unknown", status)
	}
	if stored != "" {
		t.Fatalf("у нового хоста не должно быть сохранённого отпечатка, получено %q", stored)
	}

	if err := SaveKey(path, host, port, first); err != nil {
		t.Fatal(err)
	}

	status, stored, err = Classify(path, host, port, first)
	if err != nil {
		t.Fatal(err)
	}
	if status != StatusTrusted {
		t.Fatalf("тот же ключ: статус %v, ожидался trusted", status)
	}
	if stored != ssh.FingerprintSHA256(first) {
		t.Fatalf("отпечаток %q, ожидался %q", stored, ssh.FingerprintSHA256(first))
	}

	status, stored, err = Classify(path, host, port, second)
	if err != nil {
		t.Fatal(err)
	}
	if status != StatusChanged {
		t.Fatalf("другой ключ: статус %v, ожидался changed", status)
	}
	if stored != ssh.FingerprintSHA256(first) {
		t.Fatalf("старый отпечаток %q, ожидался %q", stored, ssh.FingerprintSHA256(first))
	}

	if err := SaveKey(path, host, port, second); err != nil {
		t.Fatal(err)
	}
	status, _, err = Classify(path, host, port, second)
	if err != nil {
		t.Fatal(err)
	}
	if status != StatusTrusted {
		t.Fatalf("после перезаписи: статус %v, ожидался trusted", status)
	}
	status, _, err = Classify(path, host, port, first)
	if err != nil {
		t.Fatal(err)
	}
	if status != StatusChanged {
		t.Fatalf("старый ключ после перезаписи: статус %v, ожидался changed", status)
	}
}

func TestClassifyCustomPort(t *testing.T) {
	path := filepath.Join(t.TempDir(), "known_hosts")
	key := testPublicKey(t)

	if err := SaveKey(path, "example.com", 2222, key); err != nil {
		t.Fatal(err)
	}

	status, _, err := Classify(path, "example.com", 2222, key)
	if err != nil {
		t.Fatal(err)
	}
	if status != StatusTrusted {
		t.Fatalf("порт 2222: статус %v, ожидался trusted", status)
	}

	status, _, err = Classify(path, "example.com", 22, key)
	if err != nil {
		t.Fatal(err)
	}
	if status != StatusUnknown {
		t.Fatalf("другой порт не должен совпасть: статус %v", status)
	}
}
