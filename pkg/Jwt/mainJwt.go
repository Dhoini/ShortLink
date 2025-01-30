package Jwt

import (
	"github.com/golang-jwt/jwt/v5" // Импортируем пакет для работы с JWT.
)

// JwtDate представляет данные, которые будут храниться в токене.
type JwtDate struct {
	Email string // Email пользователя, который будет закодирован в JWT.
}

// JwtSecret представляет секретный ключ, используемый для подписи и проверки JWT.
type JwtSecret struct {
	Secret string // Секретный ключ для подписи и проверки токена.
}

// NewJWT создает новый экземпляр JwtSecret.
// secret: секретный ключ, который будет использоваться для подписи и проверки токенов.
// Возвращает: новый экземпляр JwtSecret.
func NewJWT(secret string) *JwtSecret {
	return &JwtSecret{
		Secret: secret, // Инициализируем секретный ключ.
	}
}

// Create создает новый JWT токен на основе переданных данных.
// date: данные, которые нужно закодировать в токен (например, email).
// Возвращает: строку с JWT токеном и возможную ошибку.
func (j *JwtSecret) Create(date JwtDate) (string, error) {
	// Создаем новый токен с методом подписи HMAC SHA-256 и указываем данные (claims).
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": date.Email, // Добавляем email в claims токена.
	})
	// Подписываем токен с использованием секретного ключа.
	s, err := t.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err // Возвращаем ошибку, если подписание токена не удалось.
	}
	return s, nil // Возвращаем подписанный токен.
}

// Parse проверяет и извлекает данные из JWT токена.
// token: строка с JWT токеном, который нужно проверить.
// Возвращает: флаг валидности токена и извлеченные данные (если токен валиден).
func (j *JwtSecret) Parse(token string) (bool, *JwtDate) {
	// Парсим токен, используя функцию обратного вызова для получения секретного ключа.
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil // Возвращаем секретный ключ для проверки подписи.
	})
	if err != nil {
		return false, nil // Если произошла ошибка при парсинге, возвращаем false и nil.
	}
	// Извлекаем email из claims токена.
	email := t.Claims.(jwt.MapClaims)["email"].(string)
	// Возвращаем флаг валидности токена и извлеченные данные.
	return t.Valid, &JwtDate{
		Email: email, // Возвращаем email из токена.
	}
}
