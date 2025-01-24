package links

// LinkCreateRequest представляет структуру для обработки входящих запросов на создание ссылки.
// Поле Url содержит оригинальный URL, который необходимо сократить.
type LinkCreateRequest struct {
	Url string `json:"url" validate:"required,url"` // Поле URL должно быть обязательным и содержать валидный URL.
}

type LinkUpdateRequest struct {
	Url  string `json:"url" validate:"required,url"` // Поле URL должно быть обязательным и содержать валидный URL.
	Hash string `json:"hash"`
}

type LinkDeleteRequest struct {
	Id uint `json:"id"`
}

type GetAllResponse struct {
	Links []Link `json:"links"`
	Count int64  `json:"count"`
}
