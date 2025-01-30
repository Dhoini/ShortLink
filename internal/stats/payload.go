package stats

// GetStatResponse представляет структуру для ответа на запрос получения статистики.
type GetStatResponse struct {
	Period string `json:"period"` // Период, за который собрана статистика (например, "2023-10-01").
	Sum    int    `json:"sum"`    // Общее количество кликов за указанный период.
}
