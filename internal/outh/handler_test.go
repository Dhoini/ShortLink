package auth

import (
	"Lessons/configs"
	"Lessons/internal/user"
	"Lessons/pkg/db"
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

// bootstrap инициализирует тестовую среду, создавая обработчик аутентификации и мок базы данных.
func bootstrap() (*AuthHandler, sqlmock.Sqlmock, error) {
	// Создаем мок базы данных с помощью sqlmock.
	database, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err // Возвращаем ошибку, если не удалось создать мок.
	}

	// Инициализируем GORM с использованием мока базы данных.
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: database, // Передаем мок в качестве подключения к базе данных.
	}))
	if err != nil {
		return nil, nil, err // Возвращаем ошибку, если не удалось инициализировать GORM.
	}

	// Создаем репозиторий пользователей, используя мок базы данных.
	userRepo := user.NewUserRepository(&db.Db{
		DB: gormDB,
	})

	// Инициализируем обработчик аутентификации.
	handler := AuthHandler{
		Config: &configs.Config{
			Auth: configs.AuthConfig{
				Secret: "secret", // Устанавливаем секретный ключ для JWT.
			},
		},
		AuthService: NewUserService(userRepo), // Инициализируем сервис аутентификации.
	}

	return &handler, mock, nil // Возвращаем обработчик, мок и ошибку (если есть).
}

// TestRegisterHandlerSuccess тестирует успешную регистрацию пользователя.
func TestRegisterHandlerSuccess(t *testing.T) {
	// Инициализируем тестовую среду.
	handler, mock, err := bootstrap()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}

	// Определяем ожидаемые запросы к базе данных.
	rows := sqlmock.NewRows([]string{"email", "password", "name"}) // Пустой результат для SELECT.
	mock.ExpectQuery("SELECT").WillReturnRows(rows)                // Ожидаем SELECT-запрос.
	mock.ExpectBegin()                                             // Ожидаем начало транзакции.
	mock.ExpectQuery("INSERT").                                    // Ожидаем INSERT-запрос.
									WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1)) // Возвращаем ID нового пользователя.
	mock.ExpectCommit() // Ожидаем завершение транзакции.

	// Формируем JSON-данные для запроса регистрации.
	data, _ := json.Marshal(&RegisterRequest{
		Email:    "test@gmail.com", // Email пользователя.
		Password: "1",              // Пароль пользователя.
		Name:     "test",           // Имя пользователя.
	})
	reader := bytes.NewReader(data) // Преобразуем данные в формат для отправки.

	// Создаем тестовый HTTP-запрос и записываем ответ.
	w := httptest.NewRecorder() // Записываем ответ сервера.
	req := httptest.NewRequest(http.MethodPost, "/auth/register", reader)
	handler.Register()(w, req)

	// Проверяем статус ответа.
	if w.Code != http.StatusOK {
		t.Errorf("Status code should be 200 OK, but got %d", w.Code)
		return
	}
}

// TestLoginHandlerSuccess тестирует успешный вход пользователя.
func TestLoginHandlerSuccess(t *testing.T) {
	// Инициализируем тестовую среду.
	handler, mock, err := bootstrap()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}

	// Определяем ожидаемые запросы к базе данных.
	rows := sqlmock.NewRows([]string{"email", "password"}).
		AddRow("test@gmail.com", "$2a$10$kZVz97tThenmUzHVYxPG2e8KzhnZVG/8yHsbT4PKxwktbrO8JLsE2") // Хэшированный пароль.
	mock.ExpectQuery("SELECT").
		WillReturnRows(rows) // Ожидаем SELECT-запрос для поиска пользователя.

	// Формируем JSON-данные для запроса входа.
	data, _ := json.Marshal(&LoginRequest{
		Email:    "test@gmail.com", // Email пользователя.
		Password: "1",              // Пароль пользователя.
	})
	reader := bytes.NewReader(data) // Преобразуем данные в формат для отправки.

	// Создаем тестовый HTTP-запрос и записываем ответ.
	w := httptest.NewRecorder() // Записываем ответ сервера.
	req := httptest.NewRequest(http.MethodPost, "/auth/login", reader)
	handler.Login()(w, req)

	// Проверяем статус ответа.
	if w.Code != http.StatusOK {
		t.Errorf("Status code should be 200 OK, but got %d", w.Code)
		return
	}
}
