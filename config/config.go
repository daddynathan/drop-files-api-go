package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type FileConfig struct {
	MaxFileSizeBytes int64
	MaxTTLHours      int
	DefaultTTLHours  int
	MinFilenameLen   int
	MaxFilenameLen   int
	UploadDir        string
}

// настройки сервера
type ServerConfig struct {
	Host string
	Port string
}

// настройки БД
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// Конфигурация всего приложения
type App struct {
	Server ServerConfig
	DB     DBConfig
	File   FileConfig
}

func Load() *App {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system environment")
	}
	return &App{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "dropfiles"),
			Password: getEnv("DB_PASSWORD", "dropfiles"),
			Name:     getEnv("DB_NAME", "dropfiles"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		File: loadFileConfig(),
	}
}

func loadFileConfig() FileConfig {
	maxFileSizeMB := getEnvInt("MAX_FILE_SIZE_MB", 10)
	maxTTLHours := getEnvInt("MAX_TTL_HOURS", 24)
	defaultTTLHours := getEnvInt("DEFAULT_TTL_HOURS", 24)
	minFilenameLen := getEnvInt("MIN_FILENAME_LEN", 1)
	maxFilenameLen := getEnvInt("MAX_FILENAME_LEN", 255)

	// Защита от глупости
	if maxFileSizeMB <= 0 {
		maxFileSizeMB = 10
	}
	if maxFileSizeMB > (1024 * 1000) { // до 1тб
		maxFileSizeMB = 1024 * 10 // 10гб
	}
	if maxTTLHours <= 0 {
		maxTTLHours = 24
	}
	if maxTTLHours > 168 { // до 7 дней
		maxTTLHours = 168
	}
	if defaultTTLHours <= 0 || defaultTTLHours > maxTTLHours {
		defaultTTLHours = maxTTLHours
	}
	return FileConfig{
		MaxFileSizeBytes: int64(maxFileSizeMB) << 20,
		MaxTTLHours:      maxTTLHours,
		DefaultTTLHours:  defaultTTLHours,
		MinFilenameLen:   minFilenameLen,
		MaxFilenameLen:   maxFilenameLen,
		UploadDir:        getEnv("UPLOAD_DIR", "./uploads"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if valStr := os.Getenv(key); valStr != "" {
		if val, err := strconv.Atoi(valStr); err == nil {
			return val
		}
	}
	return fallback
}
