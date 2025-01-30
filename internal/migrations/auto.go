package main

import (
	"Lessons/internal/links"
	"Lessons/internal/stats"
	"Lessons/internal/user"
	"fmt"
	"github.com/joho/godotenv" // Библиотека для загрузки переменных окружения из .env файла.
	"gorm.io/driver/postgres"  // Драйвер PostgreSQL для GORM.
	"gorm.io/gorm"             // ORM для работы с базой данных.
	"gorm.io/gorm/logger"      // Логирование запросов к базе данных.
	"os"                       // Пакет для работы с операционной системой (например, переменные окружения).
)

func main() {
	// Загружаем переменные окружения из указанного .env файла.
	err := godotenv.Load(".env")
	if err != nil {
		// Если файл не найден или произошла другая ошибка, завершить программу.
		panic(err)
	}

	fmt.Println("Подключение к базе данных успешно!")

	// Подключаемся к базе данных PostgreSQL.
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Включаем логирование запросов в консоль.
	})
	if err != nil {
		// Если подключение не удалось, завершить программу.
		panic(err)
	}

	// Автоматически создаем таблицы в базе данных на основе структур Link, User и Stats.
	// Метод AutoMigrate проверяет, существуют ли таблицы, и создает их при необходимости.
	err = db.AutoMigrate(&links.Link{}, &user.User{}, &stats.Stats{})
	if err != nil {
		// Если миграция не удалась, завершаем выполнение функции.
		fmt.Println("Ошибка при выполнении миграции:", err)
		return
	}

	// Выводим объект подключения к базе данных для проверки.
	fmt.Println("Объект подключения к базе данных:", db)
}
