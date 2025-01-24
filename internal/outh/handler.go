package auth

import (
	"Lessons/configs"
	"Lessons/pkg/JWT"
	"Lessons/pkg/reg"
	"Lessons/pkg/res"
	"net/http"
)

// AouthHendlerDeps определяет зависимости для обработчика авторизации.
// Содержит конфигурацию приложения.
type AouthHendlerDeps struct {
	*configs.Config
	*AuthService
}

// AouthHendler реализует методы для обработки запросов авторизации.
type AouthHendler struct {
	*configs.Config
	*AuthService
}

// NewAouthHendler регистрирует маршруты для авторизации.
// router: маршрутизатор HTTP.
// deps: зависимости для инициализации обработчика.
func NewAouthHendler(router *http.ServeMux, deps AouthHendlerDeps) {
	handler := &AouthHendler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}
	// Регистрируем маршруты для входа и регистрации.
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}

// Login обрабатывает запросы на вход пользователя.
func (handler *AouthHendler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Обрабатываем тело запроса и преобразуем в структуру LoginRequest.
		body, err := reg.HandleBody[LoginRequest](&w, r)
		if err != nil {
			return
		}

		email, err := handler.AuthService.Login(body.Email, body.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		token, err := JWT.NewJWT(handler.Config.Auth.Secret).Create(JWT.JWTDate{
			Email: email,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Ответ на успешный вход с фиктивным токеном.
		data := LoginResponse{
			Token: token,
		}
		res.Json(w, data, http.StatusOK)
	}
}

// Register обрабатывает запросы на регистрацию пользователя.
func (handler *AouthHendler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Обрабатываем тело запроса и преобразуем в структуру RegisterRequest.
		body, err := reg.HandleBody[RegisterRequest](&w, r)
		if err != nil {
			return
		}
		email, err := handler.AuthService.Register(body.Email, body.Password, body.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		token, err := JWT.NewJWT(handler.Config.Auth.Secret).Create(JWT.JWTDate{
			Email: email,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Ответ на успешный вход с фиктивным токеном.
		data := RegisterResponse{
			Token: token,
		}
		res.Json(w, data, http.StatusOK)
	}
}
