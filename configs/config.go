package configs

import (
	"github.com/joho/godotenv" // Библиотека для работы с .env файлами.
	"log"
	"os"
)

// Config представляет общую конфигурацию приложения.
type Config struct {
	Db   DbConfig   // Конфигурация для базы данных.
	Auth AuthConfig // Конфигурация для аутентификации.
}

// DbConfig содержит параметры подключения к базе данных.
type DbConfig struct {
	Dsn string // Data Source Name (строка подключения).
}

// AuthConfig содержит параметры для аутентификации.
type AuthConfig struct {
	Secret string // Секретный ключ для токенов.
}

// LoadConfig загружает конфигурацию из .env файла и окружения.
// Возвращает указатель на заполненную структуру Config.
func LoadConfig() *Config {
	// Загружаем .env файл.
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file") // Логируем ошибку, если файл не найден.
	}

	// Возвращаем конфигурацию, заполненную из переменных окружения.
	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"), // Читаем строку подключения к базе данных.
		},
		Auth: AuthConfig{
			Secret: os.Getenv("TOKEN"), // Читаем секретный ключ для аутентификации.
		},
	}
}
