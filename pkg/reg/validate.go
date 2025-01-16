package reg

import (
	"github.com/go-playground/validator/v10"
)

// IsValid - функция для валидации структуры типа T.
// Использует пакет go-playground/validator для проверки обязательных полей и других правил валидации.
func IsValid[T any](payload T) error {
	// Создаем новый валидатор
	validate := validator.New()
	// Выполняем валидацию структуры
	err := validate.Struct(payload)
	// Возвращаем ошибку, если валидация не прошла
	return err
}
