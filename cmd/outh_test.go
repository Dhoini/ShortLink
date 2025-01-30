package main

import (
	auth "Lessons/internal/outh" // Модуль для работы с аутентификацией.
	"Lessons/internal/user"      // Модуль для работы с пользователями.
	"bytes"                      // Пакет для работы с буферами данных.
	"encoding/json"              // Пакет для кодирования и декодирования JSON.
	"fmt"                        // Стандартный пакет для форматированного вывода.
	"github.com/joho/godotenv"   // Пакет для загрузки переменных окружения из .env файла.
	"gorm.io/driver/postgres"    // Драйвер PostgreSQL для GORM.
	"gorm.io/gorm"               // ORM для работы с базой данных.
	"gorm.io/gorm/logger"        // Логирование запросов к базе данных.
	"io"                         // Пакет для работы с вводом-выводом.
	"net/http"                   // Стандартный пакет для работы с HTTP.
	"net/http/httptest"          // Пакет для тестирования HTTP-серверов.
	"os"                         // Пакет для работы с операционной системой (например, переменные окружения).
	"testing"                    // Пакет для написания тестов.
)

// initDb инициализирует подключение к базе данных PostgreSQL.
func initDb() *gorm.DB {
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
	return db
}

// initData создает тестовые данные в базе данных.
func initData(db *gorm.DB) {
	// Создаем тестового пользователя в базе данных.
	db.Create(&user.User{
		Email:    "test@gmail.com",                                               // Email пользователя.
		Password: "$2a$10$kZVz97tThenmUzHVYxPG2e8KzhnZVG/8yHsbT4PKxwktbrO8JLsE2", // Хэшированный пароль.
		Name:     "test",                                                         // Имя пользователя.
	})
}

// removeData удаляет тестовые данные из базы данных.
func removeData(db *gorm.DB) {
	// Удаляем тестового пользователя из базы данных.
	db.Unscoped().
		Where("email = ?", "test@gmail.com").
		Delete(&user.User{})
}

// TestLoginSuccess тестирует успешный процесс входа пользователя в систему.
func TestLoginSuccess(t *testing.T) {
	// Подготовка тестовой среды:
	// 1. Инициализируем подключение к базе данных.
	db := initDb()

	// 2. Создаем тестовые данные в базе данных.
	initData(db)

	// 3. Создаем тестовый HTTP-сервер с использованием маршрутизатора приложения.
	ts := httptest.NewServer(App())
	defer ts.Close() // Убедимся, что сервер будет закрыт после завершения теста.

	// 4. Формируем JSON-данные для запроса на вход.
	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "test@gmail.com", // Email тестового пользователя.
		Password: "1",              // Пароль тестового пользователя.
	})

	// 5. Отправляем POST-запрос на эндпоинт входа.
	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		// Если произошла ошибка при отправке запроса, завершаем тест.
		t.Fatal(err)
		return
	}

	// 6. Проверяем статус ответа.
	if res.StatusCode != 200 {
		// Если статус не равен 200 (OK), завершаем тест с ошибкой.
		t.Fatalf("Expecting 200, got %d", res.StatusCode)
		return
	}

	// 7. Читаем тело ответа.
	body, err := io.ReadAll(res.Body)
	if err != nil {
		// Если произошла ошибка при чтении тела ответа, завершаем тест.
		t.Fatal(err)
		return
	}

	// 8. Декодируем JSON-ответ в структуру LoginResponse.
	var resData auth.LoginResponse
	err = json.Unmarshal(body, &resData)
	if err != nil {
		// Если произошла ошибка при декодировании JSON, завершаем тест.
		t.Fatal(err)
		return
	}

	// 9. Проверяем, что в ответе присутствует токен.
	if resData.Token == "" {
		// Если токен отсутствует, завершаем тест с ошибкой.
		t.Fatalf("Expecting token, got %s", resData.Token)
		return
	}

	// 10. Удаляем тестовые данные из базы данных.
	removeData(db)
}
