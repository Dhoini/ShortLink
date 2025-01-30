package middleware

import (
	"Lessons/configs"
	"Lessons/pkg/Jwt"
	"context"  // Импортируем пакет для работы с контекстом.
	"net/http" // Импортируем пакет для работы с HTTP-запросами и ответами.
	"strings"  // Импортируем пакет для работы со строками.
)

// key определяет тип для ключей контекста.
type key string

// ContextEmailKey — ключ для хранения email в контексте запроса.
const (
	ContextEmailKey key = "ContextEmailKey"
)

// writeUnauthorized записывает ответ "401 Unauthorized" в HTTP-ответ.
func writeUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)                              // Устанавливаем статус 401 Unauthorized.
	_, err := w.Write([]byte(http.StatusText(http.StatusUnauthorized))) // Записываем текст ошибки.
	if err != nil {
		panic(err) // Если произошла ошибка при записи, завершаем программу.
	}
}

// IsAuthenticated — middleware для проверки аутентификации пользователя.
// next: следующий обработчик HTTP-запроса.
// config: конфигурация приложения, содержащая секретный ключ для JWT.
// Возвращает: HTTP-обработчик, который проверяет токен и передает управление следующему обработчику.
func IsAuthenticated(next http.Handler, config *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем заголовок Authorization из HTTP-запроса.
		authedHeader := r.Header.Get("Authorization")
		// Проверяем, что заголовок начинается с "Bearer ".
		if !strings.HasPrefix(authedHeader, "Bearer ") {
			writeUnauthorized(w) // Если заголовок неверный, возвращаем 401 Unauthorized.
			return
		}

		// Извлекаем токен из заголовка, удаляя префикс "Bearer ".
		token := strings.TrimPrefix(authedHeader, "Bearer ")

		// Парсим токен с использованием секретного ключа из конфигурации.
		isValid, data := Jwt.NewJWT(config.Auth.Secret).Parse(token)
		if !isValid {
			writeUnauthorized(w) // Если токен невалиден, возвращаем 401 Unauthorized.
			return
		}

		// Добавляем email из токена в контекст запроса.
		ctx := context.WithValue(r.Context(), ContextEmailKey, data.Email)
		req := r.WithContext(ctx) // Создаем новый запрос с обновленным контекстом.

		// Передаем управление следующему обработчику с обновленным запросом.
		next.ServeHTTP(w, req)
	})
}
