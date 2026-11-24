package ssh

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// Status — результат сравнения ключа сервера с known_hosts.
type Status int

const (
	StatusTrusted Status = iota
	StatusUnknown
	StatusChanged
	StatusUnreachable
)

// CheckResult описывает ключ, который прислал сервер.
type CheckResult struct {
	Status               Status
	Presented            ssh.PublicKey
	PresentedFingerprint string
	StoredFingerprint    string
	Err                  error
}

var errHostKeyCaptured = errors.New("host key captured")

// HostKeyError means the server key is new or no longer matches known_hosts.
type HostKeyError struct {
	Result CheckResult
}

func (e *HostKeyError) Error() string {
	if e == nil {
		return "ключ сервера не проверен"
	}
	switch e.Result.Status {
	case StatusChanged:
		return fmt.Sprintf("ключ сервера изменился: было %s, стало %s", e.Result.StoredFingerprint, e.Result.PresentedFingerprint)
	default:
		return fmt.Sprintf("новый ключ сервера %s", e.Result.PresentedFingerprint)
	}
}

// KnownHostsPath возвращает путь к файлу known_hosts приложения.
func KnownHostsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("домашний каталог: %w", err)
	}
	return filepath.Join(home, ".ssh-keeper", "known_hosts"), nil
}

// EnsureKnownHosts создаёт каталог и пустой файл, если их ещё нет.
func EnsureKnownHosts() (string, error) {
	path, err := KnownHostsPath()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return "", err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return "", err
	}
	if err := file.Close(); err != nil {
		return "", err
	}
	return path, nil
}

// Check подключается к серверу только чтобы прочитать ключ хоста и сравнить его с файлом.
func Check(host string, port int) CheckResult {
	port = normalizePort(port)
	path, err := EnsureKnownHosts()
	if err != nil {
		return CheckResult{Status: StatusUnreachable, Err: err}
	}

	presented, err := Probe(host, port)
	if err != nil || presented == nil {
		if err == nil {
			err = errors.New("сервер не прислал ключ хоста")
		}
		return CheckResult{Status: StatusUnreachable, Err: err}
	}

	status, stored, err := Classify(path, host, port, presented)
	if err != nil {
		return CheckResult{
			Status:               StatusUnreachable,
			Presented:            presented,
			PresentedFingerprint: ssh.FingerprintSHA256(presented),
			Err:                  err,
		}
	}

	return CheckResult{
		Status:               status,
		Presented:            presented,
		PresentedFingerprint: ssh.FingerprintSHA256(presented),
		StoredFingerprint:    stored,
	}
}

// Probe читает ключ хоста и обрывает рукопожатие до аутентификации.
func Probe(host string, port int) (ssh.PublicKey, error) {
	port = normalizePort(port)
	var presented ssh.PublicKey
	config := &ssh.ClientConfig{
		User: "ssh-keeper",
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			presented = key
			return errHostKeyCaptured
		},
		Timeout: 8 * time.Second,
	}

	client, err := ssh.Dial("tcp", net.JoinHostPort(host, strconv.Itoa(port)), config)
	if client != nil {
		client.Close()
	}
	if presented != nil {
		return presented, nil
	}
	if err == nil {
		return nil, errors.New("сервер не прислал ключ хоста")
	}
	return nil, err
}

// Classify сравнивает уже полученный ключ с файлом known_hosts.
func Classify(path, host string, port int, presented ssh.PublicKey) (Status, string, error) {
	if presented == nil {
		return StatusUnreachable, "", errors.New("нет ключа хоста")
	}
	port = normalizePort(port)

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return StatusUnknown, "", nil
		}
		return StatusUnreachable, "", err
	}
	if info.Size() == 0 {
		return StatusUnknown, "", nil
	}

	callback, err := knownhosts.New(path)
	if err != nil {
		return StatusUnreachable, "", err
	}

	err = callback(dialHostname(host, port), &net.TCPAddr{Port: port}, presented)
	if err == nil {
		return StatusTrusted, ssh.FingerprintSHA256(presented), nil
	}

	var keyErr *knownhosts.KeyError
	if errors.As(err, &keyErr) {
		if len(keyErr.Want) == 0 {
			return StatusUnknown, "", nil
		}
		stored := ""
		if keyErr.Want[0].Key != nil {
			stored = ssh.FingerprintSHA256(keyErr.Want[0].Key)
		}
		return StatusChanged, stored, nil
	}
	return StatusUnreachable, "", err
}

// SaveKey записывает ключ хоста, заменяя прежние строки этого адреса.
func SaveKey(path, host string, port int, key ssh.PublicKey) error {
	if key == nil {
		return errors.New("нет ключа хоста")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}

	pattern := hostPattern(host, port)
	var kept []string
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if lineMatchesHost(line, pattern) {
			continue
		}
		kept = append(kept, line)
	}
	kept = append(kept, knownhosts.Line([]string{pattern}, key))
	content := strings.Join(kept, "\n") + "\n"
	return os.WriteFile(path, []byte(content), 0o600)
}

func hostPattern(host string, port int) string {
	host = strings.TrimSpace(host)
	port = normalizePort(port)
	needsBrackets := strings.Contains(host, ":") && !strings.HasPrefix(host, "[")
	if port != 22 || needsBrackets {
		if !strings.HasPrefix(host, "[") {
			host = "[" + host + "]"
		}
		return fmt.Sprintf("%s:%d", host, port)
	}
	return host
}

func dialHostname(host string, port int) string {
	return net.JoinHostPort(host, strconv.Itoa(normalizePort(port)))
}

func lineMatchesHost(line, pattern string) bool {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return false
	}
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return false
	}
	hostsField := fields[0]
	if strings.HasPrefix(hostsField, "@") && len(fields) > 1 {
		hostsField = fields[1]
	}
	for _, host := range strings.Split(hostsField, ",") {
		if host == pattern {
			return true
		}
	}
	return false
}

func normalizePort(port int) int {
	if port == 0 {
		return 22
	}
	return port
}
