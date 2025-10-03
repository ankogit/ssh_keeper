package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"ssh-keeper/internal/models"
)

// KeyClient представляет SSH клиент для аутентификации по ключу
type KeyClient struct {
	connection *models.Connection
}

// NewKeyClient создает новый SSH клиент для аутентификации по ключу
func NewKeyClient(conn *models.Connection) *KeyClient {
	return &KeyClient{
		connection: conn,
	}
}

// Connect устанавливает SSH подключение с использованием ключа
func (kc *KeyClient) Connect() error {
	// Восстанавливаем терминал перед запуском SSH
	kc.restoreTerminal()

	// Создаем PTY для интерактивной сессии
	pty, err := NewPTY()
	if err != nil {
		return fmt.Errorf("ошибка создания PTY: %w", err)
	}
	defer pty.Close()

	// Строим команду SSH
	args := kc.buildSSHArgs()
	cmd := exec.Command("ssh", args...)

	// Запускаем SSH в PTY с прямым подключением к стандартным потокам
	if err := pty.StartSSHWithDirectPTY(cmd); err != nil {
		return fmt.Errorf("ошибка запуска SSH: %w", err)
	}

	// Ждем завершения SSH команды
	err = cmd.Wait()

	// Восстанавливаем терминал после завершения SSH
	kc.restoreTerminal()

	return err
}

// restoreTerminal восстанавливает терминал
func (kc *KeyClient) restoreTerminal() {
	// Радикальное восстановление терминала
	os.Stdout.WriteString("\033[?1049l") // Выход из альтернативного буфера
	os.Stdout.WriteString("\033[?25h")   // Показать курсор
	os.Stdout.WriteString("\033[?2004l") // Отключаем bracketed paste mode
	os.Stdout.WriteString("\033[?1l")    // Отключаем application cursor keys
	os.Stdout.WriteString("\033[?7h")    // Включаем auto wrap mode
	os.Stdout.WriteString("\033[?12l")   // Отключаем local echo
	os.Stdout.WriteString("\033[?1000l") // Отключаем mouse reporting
	os.Stdout.WriteString("\033[?1001l") // Отключаем mouse reporting
	os.Stdout.WriteString("\033[?1002l") // Отключаем mouse reporting
	os.Stdout.WriteString("\033[?1003l") // Отключаем mouse reporting
	os.Stdout.WriteString("\033[?1005l") // Отключаем mouse reporting
	os.Stdout.WriteString("\033[?1006l") // Отключаем mouse reporting
	os.Stdout.WriteString("\033[?1015l") // Отключаем mouse reporting
	os.Stdout.WriteString("\033[?25h")   // Показать курсор
	os.Stdout.WriteString("\033[0m")     // Сбрасываем все атрибуты
	os.Stdout.WriteString("\033c")       // Полный сброс терминала
	os.Stdout.WriteString("\033[2J")     // Очищаем экран
	os.Stdout.WriteString("\033[H")      // Перемещаем курсор в начало
	os.Stdout.WriteString("\033[?25h")   // Показать курсор
	os.Stdout.WriteString("\033[0m")     // Сбрасываем все атрибуты
	os.Stdout.Sync()                     // Принудительно сбрасываем буфер

	// Радикальная защита через reset и stty
	exec.Command("reset").Run()
	exec.Command("stty", "sane").Run()
	exec.Command("tput", "reset").Run()
}

// buildSSHArgs строит аргументы для команды ssh с аутентификацией по ключу
func (kc *KeyClient) buildSSHArgs() []string {
	args := []string{}

	// Отключаем проверку host key для тестирования
	args = append(args, "-o", "StrictHostKeyChecking=no")
	args = append(args, "-o", "UserKnownHostsFile=/dev/null")

	// Порт
	if kc.connection.Port != 22 {
		args = append(args, "-p", fmt.Sprintf("%d", kc.connection.Port))
	}

	// SSH ключ
	if kc.connection.KeyPath != "" {
		// Получаем абсолютный путь к ключу
		keyPath, err := filepath.Abs(kc.connection.KeyPath)
		if err == nil {
			args = append(args, "-i", keyPath)
		}
	} else {
		// Если нет явного ключа, пробуем дефолтные ключи
		defaultKeys := kc.findDefaultSSHKeys()
		if len(defaultKeys) > 0 {
			// Используем первый найденный дефолтный ключ
			args = append(args, "-i", defaultKeys[0])
		}
	}

	// Настройки аутентификации - только ключи
	args = append(args, "-o", "PreferredAuthentications=publickey")
	args = append(args, "-o", "PubkeyAuthentication=yes")
	args = append(args, "-o", "PasswordAuthentication=no")

	// Адрес подключения
	address := fmt.Sprintf("%s@%s", kc.connection.User, kc.connection.Host)
	args = append(args, address)

	return args
}

// findDefaultSSHKeys ищет дефолтные SSH ключи в ~/.ssh/
func (kc *KeyClient) findDefaultSSHKeys() []string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	sshDir := filepath.Join(homeDir, ".ssh")

	// Стандартные имена ключей
	keyNames := []string{
		"id_rsa",
		"id_ed25519",
		"id_ecdsa",
		"id_dsa",
	}

	var foundKeys []string

	for _, keyName := range keyNames {
		keyPath := filepath.Join(sshDir, keyName)
		if _, err := os.Stat(keyPath); err == nil {
			// Проверяем, что это файл ключа (не директория)
			if info, err := os.Stat(keyPath); err == nil && !info.IsDir() {
				foundKeys = append(foundKeys, keyPath)
			}
		}
	}

	return foundKeys
}

// GetConnectionString возвращает строку подключения
func (kc *KeyClient) GetConnectionString() string {
	port := ""
	if kc.connection.Port != 22 {
		port = fmt.Sprintf(" -p %d", kc.connection.Port)
	}

	key := ""
	if kc.connection.KeyPath != "" {
		key = fmt.Sprintf(" -i %s", kc.connection.KeyPath)
	} else {
		// Показываем, что будут использоваться дефолтные ключи
		defaultKeys := kc.findDefaultSSHKeys()
		if len(defaultKeys) > 0 {
			key = fmt.Sprintf(" -i %s", defaultKeys[0])
		}
	}

	return fmt.Sprintf("ssh%s%s %s@%s", key, port, kc.connection.User, kc.connection.Host)
}

// GetAvailableKeys возвращает список доступных SSH ключей
func (kc *KeyClient) GetAvailableKeys() []string {
	if kc.connection.KeyPath != "" {
		// Если указан конкретный ключ, возвращаем его
		keyPath, err := filepath.Abs(kc.connection.KeyPath)
		if err == nil {
			return []string{keyPath}
		}
		return []string{}
	}

	// Иначе возвращаем дефолтные ключи
	return kc.findDefaultSSHKeys()
}
