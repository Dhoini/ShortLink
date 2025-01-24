package di

type IStatsRepository interface {
	AddClick(linkId uint)
}
