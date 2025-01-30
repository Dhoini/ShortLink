package res

import (
	"encoding/json" // Импортируем пакет для работы с JSON.
	"net/http"      // Импортируем пакет для работы с HTTP-запросами и ответами.
)

// Json - функция для отправки данных в формате JSON в ответе.
// Устанавливает заголовок Content-Type как "application/json", код состояния ответа и кодирует переданные данные в формат JSON.
func Json(w http.ResponseWriter, data any, status int) {
	// Устанавливаем заголовок Content-Type в "application/json",
	// чтобы клиент знал, что ответ представлен в формате JSON.
	w.Header().Set("Content-Type", "application/json")

	// Устанавливаем код состояния HTTP-ответа (например, 200 OK, 400 Bad Request и т.д.).
	w.WriteHeader(status)

	// Кодируем данные в формат JSON и отправляем их в тело HTTP-ответа.
	// Если возникнет ошибка при кодировании, она будет обработана автоматически.
	json.NewEncoder(w).Encode(data)
}
