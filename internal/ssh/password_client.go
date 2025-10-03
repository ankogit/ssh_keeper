package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"ssh-keeper/internal/models"
)

// PasswordClient представляет SSH клиент для аутентификации по паролю
type PasswordClient struct {
	connection *models.Connection
	password   string
}

// NewPasswordClient создает новый SSH клиент для аутентификации по паролю
func NewPasswordClient(conn *models.Connection) *PasswordClient {
	return &PasswordClient{
		connection: conn,
	}
}

// SetPassword устанавливает пароль для подключения
func (pc *PasswordClient) SetPassword(password string) {
	pc.password = password
}

// Connect устанавливает SSH подключение с использованием пароля
func (pc *PasswordClient) Connect() error {
	// Всегда используем PTY с передачей пароля
	return pc.ConnectWithPTY()
}

// ConnectWithPTY использует PTY для подключения с передачей пароля
func (pc *PasswordClient) ConnectWithPTY() error {
	// Восстанавливаем терминал перед запуском SSH
	pc.restoreTerminal()

	// Создаем PTY для интерактивной сессии
	pty, err := NewPTY()
	if err != nil {
		return fmt.Errorf("ошибка создания PTY: %w", err)
	}
	defer pty.Close()

	// Строим команду SSH
	args := pc.buildSSHArgs()
	cmd := exec.Command("ssh", args...)

	// Запускаем SSH в PTY с прямым подключением к стандартным потокам
	if err := pty.StartSSHWithDirectPTY(cmd); err != nil {
		return fmt.Errorf("ошибка запуска SSH: %w", err)
	}

	// Если есть пароль, передаем его в отдельной горутине
	if pc.password != "" {
		go func() {
			// Увеличиваем задержку, чтобы SSH успел запросить пароль
			time.Sleep(2 * time.Second)
			pty.Write([]byte(pc.password + "\n"))
		}()
	}

	// Ждем завершения SSH команды
	err = cmd.Wait()

	// Восстанавливаем терминал после завершения SSH
	pc.restoreTerminal()

	return err
}

// restoreTerminal восстанавливает терминал
func (pc *PasswordClient) restoreTerminal() {
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

// buildSSHArgs строит аргументы для команды ssh с аутентификацией по паролю
func (pc *PasswordClient) buildSSHArgs() []string {
	args := []string{}

	// Отключаем проверку host key для тестирования
	args = append(args, "-o", "StrictHostKeyChecking=no")
	args = append(args, "-o", "UserKnownHostsFile=/dev/null")

	// Порт
	if pc.connection.Port != 22 {
		args = append(args, "-p", fmt.Sprintf("%d", pc.connection.Port))
	}

	// Настройки аутентификации - только пароль
	args = append(args, "-o", "PreferredAuthentications=password")
	args = append(args, "-o", "PubkeyAuthentication=no")
	args = append(args, "-o", "PasswordAuthentication=yes")

	// Адрес подключения
	address := fmt.Sprintf("%s@%s", pc.connection.User, pc.connection.Host)
	args = append(args, address)

	return args
}

// GetConnectionString возвращает строку подключения
func (pc *PasswordClient) GetConnectionString() string {
	port := ""
	if pc.connection.Port != 22 {
		port = fmt.Sprintf(" -p %d", pc.connection.Port)
	}

	return fmt.Sprintf("ssh%s %s@%s (password auth)", port, pc.connection.User, pc.connection.Host)
}
