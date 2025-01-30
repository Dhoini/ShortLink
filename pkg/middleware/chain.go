package middleware

import "net/http"

// Middleware представляет функцию, которая оборачивает HTTP-обработчик (http.Handler)
// и возвращает новый HTTP-обработчик. Это позволяет добавлять дополнительную логику
// (например, аутентификацию, логирование) перед или после выполнения основного обработчика.
type Middleware func(http.Handler) http.Handler

// Chain объединяет несколько middleware в одну цепочку.
// middlewares: список middleware, которые нужно объединить.
// Возвращает: единый middleware, который последовательно применяет все переданные middleware.
func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		// Применяем middleware в обратном порядке (с конца к началу),
		// чтобы они выполнялись в правильной последовательности.
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next) // Каждый middleware оборачивает предыдущий.
		}
		return next // Возвращаем финальный обработчик, который включает всю цепочку middleware.
	}
}
