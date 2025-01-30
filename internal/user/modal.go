package user

import "gorm.io/gorm"

// User представляет структуру для хранения данных пользователя.
type User struct {
	gorm.Model // Встроенная структура GORM, содержащая поля ID, CreatedAt, UpdatedAt и DeletedAt.

	Email    string `json:"email" gorm:"type:varchar(255);uniqueIndex"` // Email пользователя. Уникальный индекс гарантирует, что email не повторяется.
	Password string `json:"password"`                                   // Хэшированный пароль пользователя.
	Name     string `json:"name" gorm:"column:name;type:varchar(255)"`  // Имя пользователя. Ограничение на длину строки до 255 символов.
}
