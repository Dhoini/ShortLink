package stats

import (
	"gorm.io/datatypes" // Пакет для работы с типами данных, такими как Date.
	"gorm.io/gorm"      // ORM для работы с базой данных.
)

// Stats представляет структуру для хранения статистики по ссылкам.
type Stats struct {
	gorm.Model // Встроенная структура GORM, содержащая поля ID, CreatedAt, UpdatedAt и DeletedAt.

	LinkID uint           `json:"link_id"` // ID ссылки, к которой относится статистика.
	Clicks int64          `json:"clicks"`  // Количество кликов по ссылке за определенный период.
	Date   datatypes.Date `json:"date"`    // Дата, за которую собрана статистика.
}
