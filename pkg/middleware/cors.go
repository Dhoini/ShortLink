package middleware

import "net/http"

// CORS — middleware для обработки Cross-Origin Resource Sharing (CORS).
// next: следующий HTTP-обработчик в цепочке.
// Возвращает: HTTP-обработчик, который добавляет заголовки CORS к ответу.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем значение заголовка Origin из HTTP-запроса.
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Если заголовок Origin отсутствует, передаем управление следующему обработчику без добавления CORS-заголовков.
			next.ServeHTTP(w, r)
			return
		}

		// Получаем доступ к заголовкам HTTP-ответа.
		header := w.Header()

		// Разрешаем запросы только с указанного источника (Origin).
		header.Set("Access-Control-Allow-Origin", origin)

		// Разрешаем отправку учетных данных (например, cookies) вместе с запросом.
		header.Set("Access-Control-Allow-Credentials", "true")

		// Если метод запроса — OPTIONS, это предварительный (preflight) запрос.
		if r.Method == http.MethodOptions {
			// Разрешаем указанные HTTP-методы.
			header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH, HEAD")

			// Устанавливаем максимальное время кэширования результатов предварительного запроса (в секундах).
			header.Set("Access-Control-Max-Age", "86400") // 86400 секунд = 1 день.

			// Разрешаем указанные заголовки в запросах.
			header.Set("Access-Control-Allow-Headers", "authorization, content-type, content-length")

			// Завершаем обработку предварительного запроса, так как основной запрос будет выполнен позже.
			return
		}

		// Передаем управление следующему обработчику для выполнения основного запроса.
		next.ServeHTTP(w, r)
	})
}
