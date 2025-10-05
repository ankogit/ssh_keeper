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
		githubOwner:    "ankogit",
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

	// Определяем путь к текущему исполняемому файлу
	currentExecutable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable path: %w", err)
	}

	// Извлекаем новый исполняемый файл
	newExecutable, err := us.extractExecutable(data, updateInfo.DownloadURL)
	if err != nil {
		return fmt.Errorf("failed to extract executable: %w", err)
	}

	// Создаем резервную копию текущего файла
	backupPath := currentExecutable + ".backup"
	if err := us.copyFile(currentExecutable, backupPath); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Заменяем текущий файл новым
	if err := us.copyFile(newExecutable, currentExecutable); err != nil {
		// Восстанавливаем из резервной копии при ошибке
		us.copyFile(backupPath, currentExecutable)
		return fmt.Errorf("failed to install update: %w", err)
	}

	// Устанавливаем права на выполнение
	if err := os.Chmod(currentExecutable, 0755); err != nil {
		return fmt.Errorf("failed to set executable permissions: %w", err)
	}

	// Удаляем временные файлы
	os.Remove(newExecutable)
	os.Remove(backupPath)

	return nil
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
	// Убираем префикс 'v' если есть
	current := strings.TrimPrefix(us.currentVersion, "v")
	latest := strings.TrimPrefix(latestVersion, "v")

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
