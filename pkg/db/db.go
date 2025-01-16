package db

import (
	"Lessons/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Db оборачивает gorm.DB и позволяет работать с базой данных.
type Db struct {
	*gorm.DB // Встраиваем gorm.DB для доступа к его методам.
}

// NewDb создает и возвращает новый объект Db.
// Он устанавливает соединение с базой данных PostgreSQL на основе конфигурации из config.
func NewDb(conf *configs.Config) *Db {
	// Открываем соединение с базой данных, используя DSN из конфигурации.
	db, err := gorm.Open(postgres.Open(conf.Db.Dsn), &gorm.Config{})
	if err != nil {
		// Если соединение не удалось, программа завершится с ошибкой.
		panic(err)
	}
	// Возвращаем обертку Db с экземпляром gorm.DB.
	return &Db{db}
}
