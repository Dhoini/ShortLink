package stats

import (
	"Lessons/pkg/db"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

// StatRepository представляет репозиторий для работы со статистикой.
type StatRepository struct {
	Database *db.Db // Указывает на базу данных, используемую для работы с сущностями Stats.
}

// NewStatRepository создает новый экземпляр StatRepository.
func NewStatRepository(db *db.Db) *StatRepository {
	return &StatRepository{db} // Инициализация базы данных.
}

// AddClick увеличивает счетчик кликов для указанной ссылки за текущую дату.
func (repo *StatRepository) AddClick(linkId uint) {
	var stat Stats
	currentDate := datatypes.Date(time.Now()) // Получаем текущую дату.

	// Ищем запись статистики для указанной ссылки и текущей даты.
	repo.Database.Find(&stat, "link_id = ? and date = ?", linkId, currentDate)

	if stat.ID == 0 {
		// Если запись не найдена, создаем новую запись с начальным количеством кликов.
		repo.Database.Create(&Stats{
			LinkID: linkId,      // ID ссылки.
			Clicks: 1,           // Начальное количество кликов.
			Date:   currentDate, // Текущая дата.
		})
	} else {
		// Если запись найдена, увеличиваем счетчик кликов и сохраняем изменения.
		stat.Clicks++
		repo.Database.Save(&stat)
	}
}

// GroupStats группирует статистику по дням или месяцам в указанном временном диапазоне.
func (repo *StatRepository) GroupStats(by string, from, to time.Time) []GetStatResponse {
	var stat []GetStatResponse
	var SelectQuery string

	// Формируем SQL-запрос для группировки данных в зависимости от параметра "by".
	switch by {
	case GroupByDay:
		SelectQuery = "to_char(date, 'YYYY-MM-DD') as period, sum(clicks)" // Группировка по дням.
	case GroupByMonth:
		SelectQuery = "to_char(date, 'YYYY-MM') as period, sum(clicks)" // Группировка по месяцам.
	}

	// Формируем основной запрос к базе данных.
	query := repo.Database.Table("stats").
		Select(SelectQuery).     // Выбираем период и сумму кликов.
		Session(&gorm.Session{}) // Создаем новую сессию для запроса.

	// Добавляем условие фильтрации (пример: count > 10).
	if true { // Это условие можно заменить на более осмысленную логику.
		query.Where("count > 10")
	}

	// Добавляем фильтрацию по временному диапазону, группировку и сортировку.
	query.
		Where("date BETWEEN ? AND ?", from, to). // Фильтруем данные по диапазону дат.
		Group("period").                         // Группируем данные по периоду.
		Order("period").                         // Сортируем данные по периоду.
		Scan(&stat)                              // Сканируем результаты в массив структур GetStatResponse.

	return stat // Возвращаем сгруппированную статистику.
}
