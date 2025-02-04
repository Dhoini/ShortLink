package auth

import (
	"Lessons/configs"
	"Lessons/pkg/Jwt"
	"Lessons/pkg/reg"
	"Lessons/pkg/res"
	"net/http"
)

// AuthHandlerDeps определяет зависимости для обработчика авторизации.
type AuthHandlerDeps struct {
	*configs.Config // Конфигурация приложения (например, секретный ключ для JWT).
	*AuthService    // Сервис для работы с аутентификацией и регистрацией пользователей.
}

// AuthHandler реализует методы для обработки запросов авторизации.
type AuthHandler struct {
	*configs.Config // Конфигурация приложения.
	*AuthService    // Сервис для работы с аутентификацией и регистрацией.
}

// NewAuthHandler регистрирует маршруты для авторизации.
func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,      // Инициализация конфигурации.
		AuthService: deps.AuthService, // Инициализация сервиса аутентификации.
	}

	// Регистрируем маршруты для входа и регистрации.
	router.HandleFunc("POST /auth/login", handler.Login())       // Маршрут для входа пользователя.
	router.HandleFunc("POST /auth/register", handler.Register()) // Маршрут для регистрации пользователя.
}

// Login обрабатывает запросы на вход пользователя.
func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Обрабатываем тело запроса и преобразуем в структуру LoginRequest.
		body, err := reg.HandleBody[LoginRequest](&w, r)
		if err != nil {
			return // Если ошибка при обработке тела запроса, прекращаем выполнение.
		}

		// Выполняем аутентификацию пользователя через AuthService.
		email, err := handler.AuthService.Login(body.Email, body.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized) // Возвращаем 401, если аутентификация не удалась.
			return
		}

		// Создаем JWT-токен для пользователя.
		token, err := Jwt.NewJWT(handler.Config.Auth.Secret).Create(Jwt.JwtDate{
			Email: email, // Email пользователя для включения в токен.
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // Возвращаем 500, если ошибка при создании токена.
			return
		}

		// Формируем ответ с токеном.
		data := LoginResponse{
			Token: token, // Токен для пользователя.
		}
		res.Json(w, data, http.StatusOK) // Возвращаем успешный ответ с токеном.
	}
}

// Register обрабатывает запросы на регистрацию пользователя.
func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Обрабатываем тело запроса и преобразуем в структуру RegisterRequest.
		body, err := reg.HandleBody[RegisterRequest](&w, r)
		if err != nil {
			return // Если ошибка при обработке тела запроса, прекращаем выполнение.
		}

		// Регистрируем пользователя через AuthService.
		email, err := handler.AuthService.Register(body.Email, body.Password, body.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized) // Возвращаем 401, если регистрация не удалась.
			return
		}

		// Создаем JWT-токен для нового пользователя.
		token, err := Jwt.NewJWT(handler.Config.Auth.Secret).Create(Jwt.JwtDate{
			Email: email, // Email пользователя для включения в токен.
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // Возвращаем 500, если ошибка при создании токена.
			return
		}

		// Формируем ответ с токеном.
		data := RegisterResponse{
			Token: token, // Токен для пользователя.
		}
		res.Json(w, data, http.StatusOK) // Возвращаем успешный ответ с токеном.
	}
}
