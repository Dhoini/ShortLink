package reg

import (
	"Lessons/pkg/res"
	"net/http"
)

// HandleBody - функция для обработки тела запроса.
// Декодирует тело в структуру типа T, выполняет валидацию и возвращает структуру.
// В случае ошибки, отправляет ошибку в ответ и возвращает nil.
func HandleBody[T any](w *http.ResponseWriter, r *http.Request) (*T, error) {
	// Декодируем тело запроса в структуру типа T
	body, err := Decode[T](r.Body)
	if err != nil {
		// Если ошибка декодирования, отправляем ошибку в формате JSON и возвращаем nil.
		res.Json(*w, err.Error(), 402)
		return nil, err
	}

	// Проверка на валидацию
	err = IsValid(body)
	if err != nil {
		// Если валидация не пройдена, отправляем ошибку в формате JSON и возвращаем nil.
		res.Json(*w, err.Error(), 402)
		return nil, err
	}

	// Если все прошло успешно, возвращаем декодированное тело
	return &body, nil
}
