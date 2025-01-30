package middleware

import "net/http"

// WrapperWriter оборачивает http.ResponseWriter для перехвата и сохранения HTTP-статуса ответа.
// Это позволяет отслеживать статус, который был отправлен клиенту.
type WrapperWriter struct {
	http.ResponseWriter     // Встраиваем оригинальный ResponseWriter для доступа к его методам.
	StatusCode          int // Поле для хранения HTTP-статуса ответа.
}

// WriteHeader перехватывает вызов метода WriteHeader и сохраняет статус в поле StatusCode.
// statusCode: HTTP-статус, который нужно отправить клиенту.
func (w *WrapperWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode) // Вызываем оригинальный метод WriteHeader.
	w.StatusCode = statusCode                // Сохраняем статус в поле StatusCode.
}
