package links

import (
	"Lessons/pkg/db"
	"gorm.io/gorm/clause"
)

// LinkRepository представляет репозиторий для работы с сущностью Link.
type LinkRepository struct {
	Database *db.Db // Указывает на базу данных, используемую для работы с сущностями Link.
}

// NewLinkRepository создает новый экземпляр LinkRepository.
// dataBase: объект базы данных.
// Возвращает: новый экземпляр LinkRepository.
func NewLinkRepository(dataBase *db.Db) *LinkRepository {
	return &LinkRepository{
		Database: dataBase,
	}
}

// Create добавляет новую запись Link в базу данных.
// link: указатель на объект Link, который нужно сохранить.
// Возвращает: сохраненный объект Link и возможную ошибку.
func (repo *LinkRepository) Create(link *Link) (*Link, error) {
	result := repo.Database.DB.Create(link) // Выполняет операцию создания записи в базе данных.
	if result.Error != nil {
		return nil, result.Error // Возвращает ошибку, если операция завершилась неудачно.
	}
	return link, nil // Возвращает сохраненный объект Link при успешной операции.
}

// GetByHash ищет ссылку по hash в базе данных
func (repo *LinkRepository) GetByHash(hash string) (*Link, error) {
	var link Link

	result := repo.Database.DB.First(&link, "hash = ?", hash)
	if result.Error != nil {
		return nil, result.Error // Возвращает ошибку, если операция завершилась неудачно.
	}

	return &link, nil
}

func (repo *LinkRepository) Update(link *Link) (*Link, error) {
	result := repo.Database.DB.Clauses(clause.Returning{}).Updates(link)
	if result.Error != nil {
		return nil, result.Error // Возвращает ошибку, если операция завершилась неудачно.
	}
	return link, nil

}

func (repo *LinkRepository) Delete(id uint) error {
	result := repo.Database.DB.Delete(id)
	if result.Error != nil {
		return result.Error // Возвращает ошибку, если операция завершилась неудачно.
	}
	return nil
}

func (repo *LinkRepository) GetById(id uint) (*Link, error) {
	var link Link

	result := repo.Database.DB.First(&link, id)
	if result.Error != nil {
		return nil, result.Error // Возвращает ошибку, если операция завершилась неудачно.
	}

	return &link, nil
}

func (repo *LinkRepository) Count() int64 {
	var count int64
	repo.Database.
		Table("links").
		Where("deleted_at IS NULL").
		Count(&count)
	return count
}

func (repo *LinkRepository) GetAll(limit, offset int) []Link {
	var links []Link
	repo.Database.
		Table("links").
		Select("*").
		Where("deleted_at IS NULL").
		Order("id ").
		Limit(int(limit)).
		Offset(int(offset)).
		Scan(&links)
	return links
}
