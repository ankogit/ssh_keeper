package version

import "os"

var (
	// Version is set via ldflags during build
	Version   string = "dev"
	BuildTime string = "dev"
)

// GetVersion возвращает версию приложения
func GetVersion() string {
	// Сначала пробуем получить из переменной окружения (для CI)
	if envVersion := os.Getenv("APP_VERSION"); envVersion != "" {
		return envVersion
	}
	// Затем используем переменную из ldflags
	if Version != "" && Version != "dev" {
		return Version
	}
	// Fallback
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
