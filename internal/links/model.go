package links

import (
	"Lessons/internal/stats"
	"gorm.io/gorm"
	"math/rand"
)

// Link представляет сущность для хранения ссылок.
type Link struct {
	gorm.Model               // Встроенная структура GORM, содержащая поля ID, CreatedAt, UpdatedAt и DeletedAt.
	Url        string        `json:"url" gorm:"type:varchar(255);not null"`          // Оригинальная ссылка, обязательное поле.
	Hash       string        `json:"hash" gorm:"uniqueIndex"`                        // Уникальный хэш для сокращенной ссылки.
	Stats      []stats.Stats `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // Связанные статистические данные.
}

// NewLink создает новую сущность Link с указанным URL и случайным хэшем.
// url: оригинальный URL.
// Возвращает: указатель на созданную сущность Link.
func NewLink(url string) *Link {
	link := &Link{
		Url: url, // Устанавливаем оригинальный URL.
	}
	link.GenereateHash() // Генерируем уникальный хэш для ссылки.
	return link
}

// GenereateHash генерирует случайный хэш для ссылки.
func (link *Link) GenereateHash() {
	link.Hash = RandStringRunes(10) // Генерируем строку длиной 10 символов.
}

// letterRunes содержит набор символов, используемых для генерации хэша.
var letterRunes = []rune("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM123456789")

// RandStringRunes генерирует случайную строку указанной длины.
// n: длина строки.
// Возвращает: случайную строку.
func RandStringRunes(n int) string {
	newString := make([]rune, n) // Создаем срез рун заданной длины.

	// Заполняем срез случайными символами из letterRunes.
	for i := range newString {
		newString[i] = letterRunes[rand.Intn(len(letterRunes))] // Случайно выбираем символ из letterRunes.
	}

	return string(newString) // Преобразуем срез рун в строку и возвращаем.
}
