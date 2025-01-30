package di

import "Lessons/internal/user"

// IStatsRepository определяет интерфейс для репозитория статистики.
// Этот интерфейс используется для абстракции операций, связанных со статистикой.
type IStatsRepository interface {
	// AddClick увеличивает счетчик кликов для указанной ссылки.
	// linkId: идентификатор ссылки, для которой нужно увеличить счетчик.
	AddClick(linkId uint)
}

// IUserRepository определяет интерфейс для репозитория пользователей.
// Этот интерфейс используется для абстракции операций, связанных с пользователями.
type IUserRepository interface {
	// Create создает нового пользователя в базе данных.
	// user: указатель на объект User, который нужно сохранить.
	// Возвращает: сохраненный объект User и возможную ошибку.
	Create(user *user.User) (*user.User, error)

	// FindByEmail ищет пользователя в базе данных по его email.
	// email: строка, содержащая email пользователя.
	// Возвращает: найденный объект User и возможную ошибку.
	FindByEmail(email string) (*user.User, error)
}
