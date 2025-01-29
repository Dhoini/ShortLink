package stats

import (
	"Lessons/pkg/db"
	"gorm.io/datatypes"
	"time"
)

type StatRepository struct {
	Database *db.Db
}

func NewStatRepository(db *db.Db) *StatRepository {
	return &StatRepository{db}
}

func (repo *StatRepository) AddClick(linkId uint) {
	var stat Stats
	currentDate := datatypes.Date(time.Now())
	repo.Database.Find(&stat, "link_id = ? and date = ?", linkId, currentDate)
	if stat.ID == 0 {
		repo.Database.Create(&Stats{
			LinkID: linkId,
			Clicks: 1,
			Date:   currentDate,
		})
	} else {
		stat.Clicks++
		repo.Database.Save(&stat)
	}
}

func (repo *StatRepository) GroupStats(by string, from, to time.Time) []GetStatResponse {
	var stat []GetStatResponse
	var SelectQuery string
	switch by {
	case GroupByDay:
		SelectQuery = "to_char(date, 'YYYY-MM-DD') as period, sum(clicks) as clicks"
	case GroupByMonth:
		SelectQuery = "to_char(date, 'YYYY-MM') as period, sum(clicks) as clicks"
	}

	repo.Database.Table("stats").
		Select(SelectQuery).
		Where("date BETWEEN ? AND ?", from, to).
		Group("period").
		Order("period").
		Scan(&stat)
	return stat
}
