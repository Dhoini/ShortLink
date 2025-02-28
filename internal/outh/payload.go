package auth

// LoginRequest представляет структуру для запроса на вход.
// Содержит email и пароль пользователя.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"` // Поле email должно быть обязательным и валидным.
	Password string `json:"password" validate:"required"`    // Поле password обязательно для заполнения.
}

// LoginResponse представляет структуру ответа на запрос входа.
// Содержит токен для пользователя после успешной авторизации.
type LoginResponse struct {
	Token string `json:"token"` // Токен, который выдается после успешного входа.
}

// RegisterRequest представляет структуру для запроса на регистрацию.
// Содержит имя пользователя и данные для входа.
type RegisterRequest struct {
	Name     string `json:"name"`                      // Поле username обязательно для заполнения.
	Email    string `json:"email" validate:"required"` // Вложенная структура LoginRequest для логина.
	Password string `json:"password" validate:"required"`
}

// RegisterResponse представляет структуру ответа на запрос регистрации.
// Содержит токен для пользователя после успешной регистрации.
type RegisterResponse struct {
	Token string `json:"token"` // Токен, который выдается после успешной регистрации.
}
