package stats

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Stats struct {
	gorm.Model
	LinkID uint           `json:"link_id"`
	Clicks int64          `json:"clicks"`
	Date   datatypes.Date `json:"date"`
}
