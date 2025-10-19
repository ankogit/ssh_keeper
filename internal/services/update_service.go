package services

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// UpdateInfo содержит информацию об обновлении
type UpdateInfo struct {
	Version     string    `json:"version"`
	ReleaseDate time.Time `json:"release_date"`
	DownloadURL string    `json:"download_url"`
	Size        int64     `json:"size"`
	Changelog   string    `json:"changelog"`
	IsAvailable bool      `json:"is_available"`
}

// UpdateService управляет проверкой и загрузкой обновлений
type UpdateService struct {
	currentVersion string
	githubOwner    string
	githubRepo     string
	httpClient     *http.Client
}

// GitHubRelease представляет структуру релиза с GitHub API
type GitHubRelease struct {
	TagName     string    `json:"tag_name"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
		Size               int64  `json:"size"`
	} `json:"assets"`
	Body string `json:"body"`
}

// NewUpdateService создает новый сервис обновлений
func NewUpdateService(currentVersion string) *UpdateService {
	return &UpdateService{
		currentVersion: currentVersion,
		githubOwner:    "ankogit", // Замените на ваш GitHub username
		githubRepo:     "ssh_keeper",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CheckForUpdates проверяет наличие обновлений
func (us *UpdateService) CheckForUpdates() (*UpdateInfo, error) {
	// Получаем последний релиз с GitHub
	latestRelease, err := us.getLatestRelease()
	if err != nil {
		return nil, fmt.Errorf("failed to get latest release: %w", err)
	}

	// Проверяем, есть ли обновление
	isUpdateAvailable := us.isUpdateAvailable(latestRelease.TagName)

	// Находим подходящий файл для текущей платформы
	downloadURL, size := us.findPlatformAsset(latestRelease.Assets)
	if downloadURL == "" {
		return nil, fmt.Errorf("no suitable release asset found for platform %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	return &UpdateInfo{
		Version:     latestRelease.TagName,
		ReleaseDate: latestRelease.PublishedAt,
		DownloadURL: downloadURL,
		Size:        size,
		Changelog:   latestRelease.Body,
		IsAvailable: isUpdateAvailable,
	}, nil
}

// DownloadAndInstallUpdate загружает и устанавливает обновление
func (us *UpdateService) DownloadAndInstallUpdate(updateInfo *UpdateInfo) error {
	// Загружаем файл
	data, err := us.downloadFile(updateInfo.DownloadURL)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}

	// Валидируем загруженный файл
	if err := us.validateDownloadedFile(data, updateInfo.DownloadURL); err != nil {
		return fmt.Errorf("downloaded file validation failed: %w", err)
	}

	// Определяем путь к текущему исполняемому файлу
	currentExecutable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable path: %w", err)
	}

	// Пробуем разные методы установки по порядку
	execDir := filepath.Dir(currentExecutable)

	// Метод 1: Прямая установка (если есть права)
	if us.canWriteToDirectory(execDir) {
		return us.installDirectly(data, updateInfo.DownloadURL, currentExecutable)
	}

	// Метод 2: Установка в пользовательскую директорию (предпочтительный)
	if err := us.installToUserDirectory(data, updateInfo.DownloadURL, currentExecutable); err == nil {
		return nil
	}

	// Метод 3: Установка через sudo
	return us.installUpdateWithSudo(data, updateInfo.DownloadURL, currentExecutable)
}

// getLatestRelease получает информацию о последнем релизе с GitHub
func (us *UpdateService) getLatestRelease() (*GitHubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", us.githubOwner, us.githubRepo)

	resp, err := us.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

// isUpdateAvailable проверяет, доступно ли обновление
func (us *UpdateService) isUpdateAvailable(latestVersion string) bool {
	// Убираем префикс 'v' если есть (костыль для правильного сравнения)
	current := strings.TrimPrefix(us.currentVersion, "v")
	latest := strings.TrimPrefix(latestVersion, "v")

	// Отладочная информация
	fmt.Printf("DEBUG: Current version: \"%s\" (clean: \"%s\")\n", us.currentVersion, current)
	fmt.Printf("DEBUG: Latest version: \"%s\" (clean: \"%s\")\n", latestVersion, latest)
	fmt.Printf("DEBUG: Are equal: %t\n", current == latest)
	fmt.Printf("DEBUG: Update available: %t\n", current != latest)

	// Дополнительная проверка: если версии одинаковые, обновление не нужно
	if current == latest {
		return false
	}

	// Простое сравнение версий (можно улучшить для семантического версионирования)
	return current != latest
}

// findPlatformAsset находит подходящий файл для текущей платформы
func (us *UpdateService) findPlatformAsset(assets []struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}) (string, int64) {
	platform := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)

	for _, asset := range assets {
		if strings.Contains(asset.Name, platform) {
			return asset.BrowserDownloadURL, asset.Size
		}
	}

	return "", 0
}

// downloadFile загружает файл по URL
func (us *UpdateService) downloadFile(url string) ([]byte, error) {
	resp, err := us.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// extractExecutable извлекает исполняемый файл из архива
func (us *UpdateService) extractExecutable(data []byte, url string) (string, error) {
	// Определяем тип архива по URL
	if strings.HasSuffix(url, ".zip") {
		return us.extractFromZip(data)
	} else if strings.HasSuffix(url, ".tar.gz") {
		return us.extractFromTarGz(data)
	}

	return "", fmt.Errorf("unsupported archive format")
}

// extractFromZip извлекает файл из ZIP архива
func (us *UpdateService) extractFromZip(data []byte) (string, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", err
	}

	for _, file := range reader.File {
		if !file.FileInfo().IsDir() && us.isExecutableFile(file.Name) {
			rc, err := file.Open()
			if err != nil {
				return "", err
			}
			defer rc.Close()

			// Создаем временный файл
			tempFile, err := os.CreateTemp("", "ssh-keeper-update-*")
			if err != nil {
				return "", err
			}
			defer tempFile.Close()

			// Копируем содержимое
			if _, err := io.Copy(tempFile, rc); err != nil {
				os.Remove(tempFile.Name())
				return "", err
			}

			return tempFile.Name(), nil
		}
	}

	return "", fmt.Errorf("no executable found in archive")
}

// extractFromTarGz извлекает файл из TAR.GZ архива
func (us *UpdateService) extractFromTarGz(data []byte) (string, error) {
	gzReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		if header.Typeflag == tar.TypeReg && us.isExecutableFile(header.Name) {
			// Создаем временный файл
			tempFile, err := os.CreateTemp("", "ssh-keeper-update-*")
			if err != nil {
				return "", err
			}
			defer tempFile.Close()

			// Копируем содержимое
			if _, err := io.Copy(tempFile, tarReader); err != nil {
				os.Remove(tempFile.Name())
				return "", err
			}

			return tempFile.Name(), nil
		}
	}

	return "", fmt.Errorf("no executable found in archive")
}

// isExecutableFile проверяет, является ли файл исполняемым
func (us *UpdateService) isExecutableFile(filename string) bool {
	// Проверяем расширение файла
	if runtime.GOOS == "windows" {
		return strings.HasSuffix(filename, ".exe")
	}

	// Для Unix-систем файл должен быть исполняемым
	// Здесь мы просто проверяем, что это не директория и не скрытый файл
	return !strings.Contains(filename, "/") && !strings.HasPrefix(filepath.Base(filename), ".")
}

// installDirectly устанавливает обновление напрямую
func (us *UpdateService) installDirectly(data []byte, downloadURL, currentExecutable string) error {
	// Извлекаем новый исполняемый файл
	newExecutable, err := us.extractExecutable(data, downloadURL)
	if err != nil {
		return fmt.Errorf("failed to extract executable: %w", err)
	}

	// Создаем временную копию текущего файла для отката
	backupFile := currentExecutable + ".backup"
	if err := us.copyFile(currentExecutable, backupFile); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Заменяем текущий файл новым
	if err := us.copyFile(newExecutable, currentExecutable); err != nil {
		// При ошибке восстанавливаем из backup
		us.copyFile(backupFile, currentExecutable)
		os.Remove(backupFile)
		return fmt.Errorf("failed to install update: %w", err)
	}

	// Устанавливаем права на выполнение
	if err := os.Chmod(currentExecutable, 0755); err != nil {
		// При ошибке восстанавливаем из backup
		us.copyFile(backupFile, currentExecutable)
		os.Remove(backupFile)
		return fmt.Errorf("failed to set executable permissions: %w", err)
	}

	// Удаляем временный файл и backup после успешной установки
	os.Remove(newExecutable)
	os.Remove(backupFile)

	return nil
}

// installToUserDirectory устанавливает в пользовательскую директорию
func (us *UpdateService) installToUserDirectory(data []byte, downloadURL, currentExecutable string) error {
	// Получаем домашнюю директорию
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Создаем директорию ~/bin
	userBin := filepath.Join(homeDir, "bin")
	if err := os.MkdirAll(userBin, 0755); err != nil {
		return fmt.Errorf("failed to create user bin directory: %w", err)
	}

	// Извлекаем новый исполняемый файл
	newExecutable, err := us.extractExecutable(data, downloadURL)
	if err != nil {
		return fmt.Errorf("failed to extract executable: %w", err)
	}
	defer os.Remove(newExecutable)

	// Устанавливаем в ~/bin/ssh-keeper
	userExecutable := filepath.Join(userBin, "ssh-keeper")
	if err := us.copyFile(newExecutable, userExecutable); err != nil {
		return fmt.Errorf("failed to install to user directory: %w", err)
	}

	// Устанавливаем права на выполнение
	if err := os.Chmod(userExecutable, 0755); err != nil {
		return fmt.Errorf("failed to set executable permissions: %w", err)
	}

	// Возвращаем специальную ошибку с инструкциями
	return fmt.Errorf("installed to %s. Please add ~/bin to your PATH and restart the application", userExecutable)
}

// installUpdateWithSudo устанавливает обновление через sudo
func (us *UpdateService) installUpdateWithSudo(data []byte, downloadURL, currentExecutable string) error {
	// Извлекаем новый исполняемый файл
	newExecutable, err := us.extractExecutable(data, downloadURL)
	if err != nil {
		return fmt.Errorf("failed to extract executable: %w", err)
	}
	defer os.Remove(newExecutable)

	// Создаем временный скрипт для обновления
	scriptContent := fmt.Sprintf(`#!/bin/bash
# Обновление SSH Keeper
echo "Installing SSH Keeper update..."
cp "%s" "%s"
chmod 755 "%s"
echo "Update completed successfully!"
`, newExecutable, currentExecutable, currentExecutable)

	// Создаем временный скрипт
	scriptFile, err := os.CreateTemp("", "ssh-keeper-update-*.sh")
	if err != nil {
		return fmt.Errorf("failed to create update script: %w", err)
	}
	defer os.Remove(scriptFile.Name())

	// Записываем скрипт
	if _, err := scriptFile.WriteString(scriptContent); err != nil {
		return fmt.Errorf("failed to write update script: %w", err)
	}
	scriptFile.Close()

	// Устанавливаем права на выполнение для скрипта
	if err := os.Chmod(scriptFile.Name(), 0755); err != nil {
		return fmt.Errorf("failed to set script permissions: %w", err)
	}

	// Запускаем скрипт через sudo
	cmd := exec.Command("sudo", scriptFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run update script with sudo: %w", err)
	}

	return nil
}

// validateDownloadedFile проверяет загруженный файл
func (us *UpdateService) validateDownloadedFile(data []byte, downloadURL string) error {
	// Проверяем минимальный размер файла (должен быть больше 1MB)
	if len(data) < 1024*1024 {
		return fmt.Errorf("downloaded file too small (%d bytes), possible corruption", len(data))
	}

	// Проверяем, что файл начинается с правильного заголовка ELF/Mach-O/PE
	if len(data) >= 4 {
		// ELF header (Linux)
		if data[0] == 0x7f && data[1] == 'E' && data[2] == 'L' && data[3] == 'F' {
			return nil
		}
		// Mach-O header (macOS)
		if data[0] == 0xfe && data[1] == 0xed && data[2] == 0xfa && data[3] == 0xce {
			return nil
		}
		if data[0] == 0xce && data[1] == 0xfa && data[2] == 0xed && data[3] == 0xfe {
			return nil
		}
		// PE header (Windows)
		if len(data) >= 64 && data[0] == 'M' && data[1] == 'Z' {
			return nil
		}
	}

	// Если не распознали заголовок, но файл достаточно большой, считаем валидным
	// (возможно, это архив tar.gz)
	return nil
}

// canWriteToDirectory проверяет, можно ли записывать в директорию
func (us *UpdateService) canWriteToDirectory(dir string) bool {
	// Пробуем создать временный файл в директории
	tempFile := filepath.Join(dir, ".ssh-keeper-test-write")
	file, err := os.Create(tempFile)
	if err != nil {
		return false
	}
	file.Close()
	os.Remove(tempFile)
	return true
}

// copyFile копирует файл
func (us *UpdateService) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// RestartApplication перезапускает приложение
func (us *UpdateService) RestartApplication() error {
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Запускаем новую версию приложения
	cmd := exec.Command(executable)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start new version: %w", err)
	}

	// Завершаем текущий процесс
	os.Exit(0)
	return nil
}

// GetUpdateInfo возвращает информацию о последней проверке обновлений
func (us *UpdateService) GetUpdateInfo() *UpdateInfo {
	// Этот метод должен возвращать кэшированную информацию об обновлениях
	// Пока что возвращаем nil, так как информация хранится в экране
	return nil
}
