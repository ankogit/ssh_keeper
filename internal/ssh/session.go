package ssh

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"ssh-keeper/internal/models"
)

const (
	// AskpassEnv включает ранний выход бинарника: напечатать пароль и завершиться.
	AskpassEnv = "SSH_KEEPER_ASKPASS"
	// PassFileEnv — путь к файлу 0600, из которого askpass читает пароль.
	PassFileEnv = "SSH_KEEPER_PASS_FILE"
)

// BuildArgs собирает аргументы ssh. Пароль в аргументы не попадает.
func BuildArgs(conn *models.Connection, knownHostsPath string) ([]string, error) {
	if conn == nil {
		return nil, fmt.Errorf("подключение не задано")
	}
	if strings.TrimSpace(conn.Host) == "" || strings.TrimSpace(conn.User) == "" {
		return nil, fmt.Errorf("не заданы хост или пользователь")
	}

	port := conn.Port
	if port == 0 {
		port = 22
	}

	args := []string{
		"-o", "StrictHostKeyChecking=yes",
		"-o", "UserKnownHostsFile=" + knownHostsPath,
		"-o", "GlobalKnownHostsFile=" + nullDevice(),
	}
	if port != 22 {
		args = append(args, "-p", strconv.Itoa(port))
	}

	if conn.UseSSHKey {
		if strings.TrimSpace(conn.KeyPath) != "" {
			keyPath, err := filepath.Abs(conn.KeyPath)
			if err != nil {
				return nil, fmt.Errorf("путь к ключу: %w", err)
			}
			args = append(args, "-i", keyPath, "-o", "IdentitiesOnly=yes")
		}
		args = append(args,
			"-o", "PreferredAuthentications=publickey",
			"-o", "PubkeyAuthentication=yes",
			"-o", "PasswordAuthentication=no",
		)
	} else {
		args = append(args,
			"-o", "PreferredAuthentications=password,keyboard-interactive",
			"-o", "PubkeyAuthentication=no",
			"-o", "PasswordAuthentication=yes",
		)
	}

	args = append(args, fmt.Sprintf("%s@%s", conn.User, conn.Host))
	return args, nil
}

// Prepare готовит интерактивную команду ssh. Если пароль сохранён, возвращает путь к файлу,
// который нужно удалить после завершения сессии.
func Prepare(conn *models.Connection) (*exec.Cmd, string, error) {
	return startCommand(context.Background(), conn, "")
}

// RunCommand выполняет одну команду на сервере и возвращает объединённый вывод.
// trust записывает новый или изменившийся ключ хоста перед запуском.
func RunCommand(conn *models.Connection, command string, timeout time.Duration, trust bool) (string, error) {
	command = strings.TrimSpace(command)
	if command == "" {
		return "", fmt.Errorf("команда не задана")
	}
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	if timeout > 120*time.Second {
		timeout = 120 * time.Second
	}

	check := Check(conn.Host, conn.Port)
	switch check.Status {
	case StatusTrusted:
	case StatusUnknown, StatusChanged:
		if !trust {
			return "", &HostKeyError{Result: check}
		}
		if check.Presented == nil {
			return "", fmt.Errorf("нет ключа сервера для записи")
		}
		path, err := KnownHostsPath()
		if err != nil {
			return "", err
		}
		port := conn.Port
		if port == 0 {
			port = 22
		}
		if err := SaveKey(path, conn.Host, port, check.Presented); err != nil {
			return "", err
		}
	default:
		if check.Err != nil {
			return "", check.Err
		}
		return "", fmt.Errorf("сервер недоступен")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd, passFile, err := startCommand(ctx, conn, command)
	if passFile != "" {
		defer os.Remove(passFile)
	}
	if err != nil {
		return "", err
	}

	out, err := cmd.CombinedOutput()
	text := trimOutput(out)
	if ctx.Err() == context.DeadlineExceeded {
		return text, fmt.Errorf("команда превысила время ожидания")
	}
	if err != nil {
		if text == "" {
			return "", err
		}
		return text, fmt.Errorf("код выхода: %w", err)
	}
	return text, nil
}

func startCommand(ctx context.Context, conn *models.Connection, remote string) (*exec.Cmd, string, error) {
	if _, err := exec.LookPath("ssh"); err != nil {
		return nil, "", fmt.Errorf("команда ssh не найдена в PATH")
	}

	knownHosts, err := EnsureKnownHosts()
	if err != nil {
		return nil, "", err
	}

	args, err := BuildArgs(conn, knownHosts)
	if err != nil {
		return nil, "", err
	}
	if strings.TrimSpace(remote) != "" {
		destination := args[len(args)-1]
		args = append(args[:len(args)-1], "-o", "ConnectTimeout=10", destination, remote)
	}

	cmd := exec.CommandContext(ctx, "ssh", args...)

	if conn.UseSSHKey || conn.Password == "" {
		return cmd, "", nil
	}

	passFile, err := writePassFile(conn.Password)
	if err != nil {
		return nil, "", err
	}

	exe, err := os.Executable()
	if err != nil {
		os.Remove(passFile)
		return nil, "", fmt.Errorf("путь к программе: %w", err)
	}

	cmd.Env = append(os.Environ(),
		AskpassEnv+"=1",
		PassFileEnv+"="+passFile,
		"SSH_ASKPASS="+exe,
		"SSH_ASKPASS_REQUIRE=force",
	)
	return cmd, passFile, nil
}

func writePassFile(password string) (string, error) {
	file, err := os.CreateTemp("", "ssh-keeper-askpass-*")
	if err != nil {
		return "", fmt.Errorf("файл пароля: %w", err)
	}
	name := file.Name()
	if err := file.Chmod(0o600); err != nil {
		file.Close()
		os.Remove(name)
		return "", err
	}
	if !strings.HasSuffix(password, "\n") {
		password += "\n"
	}
	if _, err := file.WriteString(password); err != nil {
		file.Close()
		os.Remove(name)
		return "", err
	}
	if err := file.Close(); err != nil {
		os.Remove(name)
		return "", err
	}
	return name, nil
}

func trimOutput(out []byte) string {
	const limit = 100_000
	if len(out) <= limit {
		return string(out)
	}
	return string(out[:limit]) + "\n…вывод обрезан"
}

func nullDevice() string {
	if runtime.GOOS == "windows" {
		return "NUL"
	}
	return "/dev/null"
}
