package links

// LinkCreateRequest представляет структуру для обработки входящих запросов на создание ссылки.
type LinkCreateRequest struct {
	Url string `json:"url" validate:"required,url"`
	// Поле Url содержит оригинальный URL, который необходимо сократить.
	// Валидация: поле является обязательным (required) и должно быть валидным URL-адресом (url).
}

// LinkUpdateRequest представляет структуру для обработки входящих запросов на обновление ссылки.
type LinkUpdateRequest struct {
	Url string `json:"url" validate:"required,url"`
	// Поле Url содержит новый оригинальный URL, который необходимо обновить.
	// Валидация: поле является обязательным (required) и должно быть валидным URL-адресом (url).

	Hash string `json:"hash"`
	// Поле Hash содержит новый хэш для сокращенной ссылки (опционально).
}

// LinkDeleteRequest представляет структуру для обработки входящих запросов на удаление ссылки.
type LinkDeleteRequest struct {
	Id uint `json:"id"`
	// Поле Id содержит уникальный идентификатор ссылки, которую необходимо удалить.
}

// GetAllResponse представляет структуру для ответа на запрос получения всех ссылок с пагинацией.
type GetAllResponse struct {
	Links []Link `json:"links"`
	// Поле Links содержит список ссылок, полученных из базы данных.

	Count int64 `json:"count"`
	// Поле Count содержит общее количество ссылок в базе данных.
}
