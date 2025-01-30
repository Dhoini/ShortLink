package links

import (
	"Lessons/pkg/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// LinkRepository представляет репозиторий для работы с сущностью Link.
type LinkRepository struct {
	Database *db.Db // Указывает на базу данных, используемую для работы с сущностями Link.
}

// NewLinkRepository создает новый экземпляр LinkRepository.
func NewLinkRepository(dataBase *db.Db) *LinkRepository {
	return &LinkRepository{
		Database: dataBase, // Инициализация базы данных.
	}
}

// Create добавляет новую запись Link в базу данных.
func (repo *LinkRepository) Create(link *Link) (*Link, error) {
	result := repo.Database.DB.Create(link) // Выполняет операцию создания записи в базе данных.
	if result.Error != nil {
		return nil, result.Error // Возвращает ошибку, если операция завершилась неудачно.
	}
	return link, nil // Возвращает сохраненный объект Link при успешной операции.
}

// GetByHash ищет ссылку по хешу в базе данных.
func (repo *LinkRepository) GetByHash(hash string) (*Link, error) {
	var link Link
	result := repo.Database.DB.First(&link, "hash = ?", hash) // Ищет первую запись с указанным хешем.
	if result.Error != nil {
		return nil, result.Error // Возвращает ошибку, если операция завершилась неудачно.
	}
	return &link, nil // Возвращает найденную ссылку.
}

// Update обновляет существующую запись Link в базе данных.
func (repo *LinkRepository) Update(link *Link) (*Link, error) {
	result := repo.Database.DB.Clauses(clause.Returning{}).Updates(link) // Обновляет запись и возвращает обновленные данные.
	if result.Error != nil {
		return nil, result.Error // Возвращает ошибку, если операция завершилась неудачно.
	}
	return link, nil // Возвращает обновленную ссылку.
}

// Delete удаляет запись Link из базы данных по ID.
func (repo *LinkRepository) Delete(id uint) error {
	result := repo.Database.DB.Delete(&Link{}, id) // Удаляет запись с указанным ID.
	if result.Error != nil {
		return result.Error // Возвращает ошибку, если операция завершилась неудачно.
	}
	return nil // Возвращает nil при успешном удалении.
}

// GetById ищет ссылку по ID в базе данных.
func (repo *LinkRepository) GetById(id uint) (*Link, error) {
	var link Link
	result := repo.Database.DB.First(&link, id) // Ищет первую запись с указанным ID.
	if result.Error != nil {
		return nil, result.Error // Возвращает ошибку, если операция завершилась неудачно.
	}
	return &link, nil // Возвращает найденную ссылку.
}

// Count возвращает общее количество активных ссылок в базе данных.
func (repo *LinkRepository) Count() int64 {
	var count int64
	repo.Database.
		Table("links").              // Указываем таблицу для запроса.
		Where("deleted_at IS NULL"). // Фильтруем только активные записи (не удаленные).
		Count(&count)                // Подсчитываем количество записей.
	return count // Возвращаем общее количество активных ссылок.
}

// GetAll возвращает список ссылок с пагинацией.
func (repo *LinkRepository) GetAll(limit, offset int) []Link {
	var links []Link
	query := repo.Database.
		Table("links").              // Указываем таблицу для запроса.
		Select("*").                 // Выбираем все поля.
		Where("deleted_at IS NULL"). // Фильтруем только активные записи (не удаленные).
		Session(&gorm.Session{})     // Создаем новую сессию для запроса.

	query.
		Order("id").         // Сортируем записи по ID.
		Limit(int(limit)).   // Ограничиваем количество записей согласно параметру limit.
		Offset(int(offset)). // Пропускаем записи согласно параметру offset.
		Scan(&links)         // Сканируем результаты в массив ссылок.

	return links // Возвращаем список ссылок.
}
