package reg

import (
	"Lessons/pkg/res" // Импортируем пакет для работы с JSON-ответами.
	"net/http"        // Импортируем пакет для работы с HTTP-запросами и ответами.
)

// HandleBody - функция для обработки тела запроса.
// Декодирует тело в структуру типа T, выполняет валидацию и возвращает структуру.
// В случае ошибки, отправляет ошибку в ответ и возвращает nil.
func HandleBody[T any](w *http.ResponseWriter, r *http.Request) (*T, error) {
	// Декодируем тело запроса в структуру типа T
	body, err := Decode[T](r.Body)
	if err != nil {
		// Если ошибка декодирования, отправляем ошибку в формате JSON с кодом 402 (Payment Required).
		res.Json(*w, err.Error(), 402)
		return nil, err // Возвращаем nil и ошибку декодирования.
	}

	// Проверка на валидацию данных
	err = IsValid(body)
	if err != nil {
		// Если валидация не пройдена, отправляем ошибку в формате JSON с кодом 402.
		res.Json(*w, err.Error(), 402)
		return nil, err // Возвращаем nil и ошибку валидации.
	}

	// Если все прошло успешно, возвращаем указатель на декодированное тело.
	return &body, nil
}
