package middleware

import (
	"log"      // Импортируем пакет для логирования.
	"net/http" // Импортируем пакет для работы с HTTP-запросами и ответами.
	"time"     // Импортируем пакет для работы со временем.
)

// Logging — middleware для логирования HTTP-запросов.
// next: следующий HTTP-обработчик в цепочке.
// Возвращает: HTTP-обработчик, который логирует информацию о запросе и времени его выполнения.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now() // Запоминаем время начала обработки запроса.

		// Создаем обертку WrapperWriter для перехвата статуса ответа.
		wrapper := &WrapperWriter{
			ResponseWriter: w,             // Передаем оригинальный ResponseWriter.
			StatusCode:     http.StatusOK, // Устанавливаем статус по умолчанию (200 OK).
		}

		// Передаем управление следующему обработчику с использованием обертки.
		next.ServeHTTP(wrapper, r)

		// Логируем информацию о запросе:
		// - Код статуса HTTP-ответа.
		// - Метод HTTP-запроса.
		// - Путь URL запроса.
		// - Время выполнения запроса.
		log.Println(wrapper.StatusCode, r.Method, r.URL.Path, time.Since(start))
	})
}
