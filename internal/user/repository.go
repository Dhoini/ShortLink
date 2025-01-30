package user

import (
	"Lessons/pkg/db" // Импортируем пакет db, который содержит логику работы с базой данных.
	"fmt"            // Импортируем пакет fmt для форматированного вывода.
)

// UserRepository представляет репозиторий для работы с пользователями в базе данных.
type UserRepository struct {
	Database *db.Db // Ссылка на объект базы данных, через который выполняются запросы.
}

// NewUserRepository создает новый экземпляр UserRepository.
// dataBase: объект базы данных.
// Возвращает: новый экземпляр UserRepository.
func NewUserRepository(dataBase *db.Db) *UserRepository {
	return &UserRepository{
		Database: dataBase, // Инициализируем поле Database переданным объектом базы данных.
	}
}

// Create добавляет новую запись User в базу данных.
// user: указатель на объект User, который нужно сохранить.
// Возвращает: сохраненный объект User и возможную ошибку.
func (repo *UserRepository) Create(user *User) (*User, error) {
	result := repo.Database.DB.Create(user) // Выполняем операцию создания записи в базе данных.
	if result.Error != nil {
		return nil, result.Error // Если произошла ошибка при создании, возвращаем её.
	}
	fmt.Printf("User before saving: %+v\n", user) // Выводим информацию о пользователе перед сохранением (для отладки).
	return user, nil                              // Возвращаем сохраненный объект User при успешной операции.
}

// FindByEmail ищет пользователя в базе данных по его email.
// email: строка, содержащая email пользователя.
// Возвращает: найденный объект User и возможную ошибку.
func (repo *UserRepository) FindByEmail(email string) (*User, error) {
	var user User                                               // Объявляем переменную для хранения найденного пользователя.
	result := repo.Database.DB.First(&user, "email = ?", email) // Ищем первую запись, где email совпадает с заданным.
	if result.Error != nil {
		return nil, result.Error // Если произошла ошибка или пользователь не найден, возвращаем ошибку.
	}
	return &user, nil // Возвращаем указатель на найденного пользователя.
}
