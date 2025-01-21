package user

import (
	"Lessons/pkg/db"
	"fmt"
)

//create findByemail

type UserRepository struct {
	Database *db.Db
}

// NewUserRepository создает новый экземпляр UserRepository.
// dataBase: объект базы данных.
// Возвращает: новый экземпляр LinkRepository.
func NewUserRepository(dataBase *db.Db) *UserRepository {
	return &UserRepository{
		Database: dataBase,
	}
}

// Create добавляет новую запись User в базу данных.
// user: указатель на объект User, который нужно сохранить.
// Возвращает: сохраненный объект User и возможную ошибку.
func (repo *UserRepository) Create(user *User) (*User, error) {
	result := repo.Database.DB.Create(user) // Выполняет операцию создания записи в базе данных.
	if result.Error != nil {
		return nil, result.Error // Возвращает ошибку, если операция завершилась неудачно.
	}
	fmt.Printf("User before saving: %+v\n", user)

	return user, nil // Возвращает сохраненный объект User при успешной операции.
}

func (repo *UserRepository) FindByEmail(email string) (*User, error) {
	var user User

	result := repo.Database.DB.First(&user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error // Возвращает ошибку, если операция завершилась неудачно.
	}

	return &user, nil
}
