package auth

import (
	"Lessons/internal/user"
	"testing"
)

// MockUserRepository представляет мок-репозиторий для тестирования сервиса аутентификации.
type MockUserRepository struct{}

// Create имитирует создание нового пользователя в базе данных.
func (repo *MockUserRepository) Create(u *user.User) (*user.User, error) {
	// Возвращаем пользователя с предопределенным email для тестирования.
	return &user.User{
		Email: initialEmail, // Используем константу initialEmail.
	}, nil
}

// FindByEmail имитирует поиск пользователя по email.
func (repo *MockUserRepository) FindByEmail(email string) (*user.User, error) {
	// Возвращаем nil, чтобы имитировать отсутствие пользователя в базе данных.
	return nil, nil
}

// Константа initialEmail используется для тестирования.
const initialEmail = "test@test.com"

// TestRegisterSuccess тестирует успешную регистрацию пользователя.
func TestRegisterSuccess(t *testing.T) {
	// Создаем экземпляр AuthService с использованием мок-репозитория.
	authService := NewUserService(&MockUserRepository{})

	// Вызываем метод Register для регистрации нового пользователя.
	email, err := authService.Register(initialEmail, "0", "test")
	if err != nil {
		t.Fatal(err) // Если произошла ошибка, завершаем тест с сообщением об ошибке.
	}

	// Проверяем, что возвращенный email соответствует ожидаемому значению.
	if email != initialEmail {
		t.Fatalf("Expected email %s, but got %s", initialEmail, email)
	}
}
