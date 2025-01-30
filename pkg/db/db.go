package db

import (
	"Lessons/configs"         // Импортируем пакет configs, который содержит конфигурацию приложения.
	"gorm.io/driver/postgres" // Импортируем драйвер PostgreSQL для GORM.
	"gorm.io/gorm"            // Импортируем GORM — ORM для работы с базой данных.
)

// Db оборачивает gorm.DB и позволяет работать с базой данных.
type Db struct {
	*gorm.DB // Встраиваем gorm.DB для доступа к его методам.
}

// NewDb создает и возвращает новый объект Db.
// Он устанавливает соединение с базой данных PostgreSQL на основе конфигурации из config.
func NewDb(conf *configs.Config) *Db {
	// Открываем соединение с базой данных, используя DSN (Data Source Name) из конфигурации.
	db, err := gorm.Open(postgres.Open(conf.Db.Dsn), &gorm.Config{})
	if err != nil {
		// Если соединение не удалось, программа завершится с ошибкой.
		// Это критическая ошибка, так как работа без базы данных невозможна.
		panic(err)
	}
	// Возвращаем обертку Db с экземпляром gorm.DB.
	// Теперь можно использовать этот объект для выполнения запросов к базе данных.
	return &Db{db}
}
