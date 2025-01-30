package di

import "Lessons/internal/user"

type IStatsRepository interface {
	AddClick(linkId uint)
}
type IUserRepository interface {
	Create(user *user.User) (*user.User, error)
	FindByEmail(email string) (*user.User, error)
}
