package auth

import (
	"Lessons/internal/user"
	"Lessons/pkg/di"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// AuthService представляет сервис для работы с аутентификацией и регистрацией пользователей.
type AuthService struct {
	UserRepository di.IUserRepository // Репозиторий для работы с пользователями.
}

// NewUserService создает новый экземпляр AuthService.
func NewUserService(userRepository di.IUserRepository) *AuthService {
	return &AuthService{
		UserRepository: userRepository, // Инициализация репозитория пользователей.
	}
}

// Register регистрирует нового пользователя в системе.
func (service *AuthService) Register(email, password, name string) (string, error) {
	// Проверяем, существует ли пользователь с таким email.
	existedUser, _ := service.UserRepository.FindByEmail(email)
	if existedUser != nil {
		return "", errors.New(ErrUserExisted) // Возвращаем ошибку, если пользователь уже существует.
	}

	// Хэшируем пароль перед сохранением в базу данных.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err // Возвращаем ошибку, если не удалось хэшировать пароль.
	}

	// Создаем нового пользователя.
	user := &user.User{
		Email:    email,                  // Устанавливаем email.
		Password: string(hashedPassword), // Устанавливаем хэшированный пароль.
		Name:     name,                   // Устанавливаем имя пользователя.
	}

	// Сохраняем пользователя в базе данных.
	_, err = service.UserRepository.Create(user)
	if err != nil {
		return "", err // Возвращаем ошибку, если не удалось создать пользователя.
	}

	return user.Email, nil // Возвращаем email зарегистрированного пользователя.
}

// Login выполняет вход пользователя в систему.
func (service *AuthService) Login(email, password string) (string, error) {
	// Ищем пользователя по email.
	existedUser, _ := service.UserRepository.FindByEmail(email)
	if existedUser == nil {
		return "", errors.New(ErrWrongCredetials) // Возвращаем ошибку, если пользователь не найден.
	}

	// Проверяем, соответствует ли пароль хэшу из базы данных.
	err := bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(password))
	if err != nil {
		return "", errors.New(ErrWrongCredetials) // Возвращаем ошибку, если пароль неверный.
	}

	return existedUser.Email, nil // Возвращаем email пользователя при успешном входе.
}
