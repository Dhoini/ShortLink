package main

import (
	"Lessons/internal/links"
	"Lessons/internal/user"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
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
		Logger: logger.Default.LogMode(logger.Info), // Включаем логирование запросов.
	})
	if err != nil {
		// Если подключение не удалось, завершить программу.
		panic(err)
	}

	// Уведомление о начале миграции.
	fmt.Println("Миграция выполнена успешно!")

	// Удаляем таблицу Link, если она существует.
	db.Migrator().DropTable(&links.Link{}, &user.User{})

	// Автоматически создаем таблицу Link на основе структуры.
	err = db.AutoMigrate(&links.Link{}, &user.User{})
	if err != nil {
		// Если миграция не удалась, завершаем выполнение функции.
		return
	}

	// Выводим объект подключения к базе данных для проверки.
	fmt.Println(db)
}
