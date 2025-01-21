package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"type:varchar(255);uniqueIndex"`
	Password string `json:"password"`
	Name     string `json:"name" gorm:"column:name;type:varchar(255)"`
}
