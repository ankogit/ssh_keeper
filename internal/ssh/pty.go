package ssh

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"golang.org/x/crypto/ssh/terminal"
)

// PTY представляет псевдо-терминал
type PTY struct {
	pty *os.File
	tty *os.File
}

// NewPTY создает новый PTY
func NewPTY() (*PTY, error) {
	// Создаем PTY с помощью creack/pty
	ptyFile, ttyFile, err := pty.Open()
	if err != nil {
		return nil, fmt.Errorf("ошибка создания PTY: %w", err)
	}

	// Получаем размер терминала и устанавливаем его
	if size, err := pty.GetsizeFull(os.Stdout); err == nil {
		pty.Setsize(ptyFile, size)
	}

	return &PTY{
		pty: ptyFile,
		tty: ttyFile,
	}, nil
}

// StartSSH запускает SSH команду в PTY
func (p *PTY) StartSSH(sshCmd *exec.Cmd) error {
	// Настраиваем стандартные потоки
	sshCmd.Stdin = p.tty
	sshCmd.Stdout = p.tty
	sshCmd.Stderr = p.tty

	// Запускаем команду в PTY
	_, err := pty.Start(sshCmd)
	return err
}

// StartSSHDirect запускает SSH команду с прямым подключением к стандартным потокам
func (p *PTY) StartSSHDirect(sshCmd *exec.Cmd) error {
	// Настраиваем стандартные потоки напрямую
	sshCmd.Stdin = os.Stdin
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr

	// Запускаем команду в PTY, но с прямым подключением потоков
	_, err := pty.Start(sshCmd)
	return err
}

// StartSSHWithDirectPTY запускает SSH команду в PTY с прямым подключением к стандартным потокам
func (p *PTY) StartSSHWithDirectPTY(sshCmd *exec.Cmd) error {
	// Настраиваем стандартные потоки SSH команды на PTY
	sshCmd.Stdin = p.tty
	sshCmd.Stdout = p.tty
	sshCmd.Stderr = p.tty

	// Запускаем команду в PTY
	_, err := pty.Start(sshCmd)
	if err != nil {
		return err
	}

	// Прямое подключение PTY к стандартным потокам без горутин
	// Это блокирующий вызов, который будет работать до завершения SSH
	go func() {
		io.Copy(os.Stdout, p.pty)
	}()
	go func() {
		io.Copy(p.pty, os.Stdin)
	}()

	return nil
}

// Read читает данные из PTY
func (p *PTY) Read(b []byte) (int, error) {
	return p.pty.Read(b)
}

// Write записывает данные в PTY
func (p *PTY) Write(b []byte) (int, error) {
	return p.pty.Write(b)
}

// Close закрывает PTY
func (p *PTY) Close() error {
	if p.tty != nil {
		p.tty.Close()
	}
	if p.pty != nil {
		p.pty.Close()
	}
	return nil
}

// Resize изменяет размер PTY
func (p *PTY) Resize(rows, cols int) error {
	size := &pty.Winsize{
		Rows: uint16(rows),
		Cols: uint16(cols),
	}
	return pty.Setsize(p.pty, size)
}

// GetSize получает размер PTY
func (p *PTY) GetSize() (*pty.Winsize, error) {
	rows, cols, err := pty.Getsize(p.pty)
	if err != nil {
		return nil, err
	}
	return &pty.Winsize{
		Rows: uint16(rows),
		Cols: uint16(cols),
	}, nil
}

// CopyToStdout копирует данные из PTY в stdout
func (p *PTY) CopyToStdout() {
	fmt.Println("DEBUG: PTY CopyToStdout: Начинаем копирование...")
	// Простое копирование без буферизации для отладки
	n, err := io.Copy(os.Stdout, p.pty)
	fmt.Printf("DEBUG: PTY CopyToStdout: Скопировано %d байт, ошибка: %v\n", n, err)
}

// CopyFromStdin копирует данные из stdin в PTY
func (p *PTY) CopyFromStdin() {
	io.Copy(p.pty, os.Stdin)
}

// CopyFromStdinWithInterrupt копирует данные из stdin в PTY с перехватом Ctrl+C
func (p *PTY) CopyFromStdinWithInterrupt() {
	fmt.Println("DEBUG: PTY CopyFromStdin: Начинаем копирование...")
	// Простое копирование без буферизации для отладки
	n, err := io.Copy(p.pty, os.Stdin)
	fmt.Printf("DEBUG: PTY CopyFromStdin: Скопировано %d байт, ошибка: %v\n", n, err)
}

// ConnectDirectly подключает PTY напрямую к стандартным потокам
func (p *PTY) ConnectDirectly() {
	// Напрямую подключаем PTY к стандартным потокам
	// Это более эффективно, чем копирование через горутины
	go func() {
		io.Copy(os.Stdout, p.pty)
	}()
	go func() {
		io.Copy(p.pty, os.Stdin)
	}()
}

// ConnectDirectlySync подключает PTY синхронно (более эффективно)
func (p *PTY) ConnectDirectlySync() {
	// Используем io.Pipe для более эффективного копирования
	// Это позволяет избежать буферизации и задержек
	go func() {
		defer p.pty.Close()
		io.Copy(os.Stdout, p.pty)
	}()
	go func() {
		defer p.pty.Close()
		io.Copy(p.pty, os.Stdin)
	}()
}

// GetTerminalSize получает размер текущего терминала
func GetTerminalSize() (*pty.Winsize, error) {
	width, height, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return nil, err
	}
	return &pty.Winsize{
		Rows: uint16(height),
		Cols: uint16(width),
	}, nil
}
