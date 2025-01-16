package reg

import (
	"encoding/json"
	"io"
)

// Decode - универсальная функция для декодирования JSON-данных в структуру.
// Принимает в качестве аргумента тело запроса и возвращает структуру типа T.
// Если декодирование не удалось, возвращается ошибка.
func Decode[T any](body io.ReadCloser) (T, error) {
	var payload T
	// Используем json.NewDecoder для декодирования данных из потока.
	if err := json.NewDecoder(body).Decode(&payload); err != nil {
		// Если ошибка при декодировании, возвращаем пустую структуру и ошибку.
		return payload, err
	}
	// Возвращаем декодированные данные.
	return payload, nil
}
