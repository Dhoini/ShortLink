package Jwt

import (
	"testing" // Импортируем пакет testing для написания unit-тестов.
)

// TestJWT_Create тестирует методы Create и Parse из пакета Jwt.
func TestJWT_Create(t *testing.T) {
	const email = "test@test.com" // Задаем тестовый email, который будет использоваться в JWT.

	// Создаем новый экземпляр JwtSecret с секретным ключом.
	jwtService := NewJWT("RxbxgRcFCFes0enila83XSdWzejBmKuw4cHiPuMgiU8")

	// Создаем JWT токен с использованием метода Create.
	token, err := jwtService.Create(JwtDate{
		Email: "test@test.com", // Передаем тестовый email в данные токена.
	})
	if err != nil {
		t.Fatal(err) // Если произошла ошибка при создании токена, завершаем тест с ошибкой.
	}

	// Проверяем созданный токен с помощью метода Parse.
	isValid, data := jwtService.Parse(token)
	if !isValid {
		t.Fatal("token is not valid") // Если токен невалиден, завершаем тест с ошибкой.
	}

	// Проверяем, что email извлеченных данных соответствует исходному email.
	if data.Email != email {
		t.Fatalf("data:%s is not equal to %s", data.Email, email) // Если email не совпадает, завершаем тест с ошибкой.
	}
}
