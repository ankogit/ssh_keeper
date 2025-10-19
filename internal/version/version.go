package version

import (
	"os"
	"strings"
)

var (
	// Version is set via ldflags during build
	Version   string = "dev"
	BuildTime string = "dev"
)

// GetVersion возвращает версию приложения
func GetVersion() string {
	// Используем переменную из ldflags (установлена во время сборки)
	if Version != "" && Version != "dev" {
		// Убираем префикс 'v' если есть
		return strings.TrimPrefix(Version, "v")
	}
	// Fallback для разработки
	return "dev"
}

// GetBuildTime возвращает время сборки
func GetBuildTime() string {
	// Сначала пробуем получить из переменной окружения (для CI)
	if envBuildTime := os.Getenv("BUILD_TIME"); envBuildTime != "" {
		return envBuildTime
	}
	// Затем используем переменную из ldflags
	if BuildTime != "" && BuildTime != "dev" {
		return BuildTime
	}
	// Fallback
	return "dev"
}
