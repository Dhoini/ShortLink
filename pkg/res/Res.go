package res

import (
	"encoding/json"
	"net/http"
)

// Json - функция для отправки данных в формате JSON в ответе.
// Устанавливает заголовок Content-Type как "application/json", код состояния ответа и кодирует переданные данные в формат JSON.
func Json(w http.ResponseWriter, data any, status int) {
	// Устанавливаем заголовок Content-Type в "application/json", чтобы клиент знал, что ответ в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// Устанавливаем код состояния ответа (например, 200, 400 и т.д.)
	w.WriteHeader(status)
	// Кодируем данные в JSON и отправляем их в ответ
	json.NewEncoder(w).Encode(data)
}
