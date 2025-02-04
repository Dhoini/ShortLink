package stats

import (
	"Lessons/configs"
	"Lessons/pkg/middleware"
	"Lessons/pkg/res"
	"net/http"
	"time"
)

// Константы для группировки статистики.
const (
	GroupByDay   = "day"   // Группировка по дням.
	GroupByMonth = "month" // Группировка по месяцам.
)

// StatHandlerDeps определяет зависимости для обработчика статистики.
type StatHandlerDeps struct {
	StatRepository *StatRepository // Репозиторий для работы со статистикой.
	Config         *configs.Config // Конфигурация приложения.
}

// StatHandler отвечает за обработку HTTP-запросов, связанных со статистикой.
type StatHandler struct {
	StatRepository *StatRepository // Репозиторий для доступа к данным статистики.
}

// NewStatHandler создает новый обработчик статистики и регистрирует маршруты.
func NewStatHandler(router *http.ServeMux, deps StatHandlerDeps) {
	handler := &StatHandler{
		StatRepository: deps.StatRepository, // Инициализация репозитория статистики.
	}

	// Регистрация маршрутов.
	router.Handle("GET /stat", middleware.IsAuthenticated(handler.GetStat(), deps.Config)) // Маршрут для получения статистики (требует аутентификации).
}

// GetStat обрабатывает запросы на получение статистики.
func (h *StatHandler) GetStat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем параметр "from" из URL и преобразуем его в объект time.Time.
		from, err := time.Parse("2006-01-02", r.URL.Query().Get("from"))
		if err != nil {
			http.Error(w, "Invalid param :from:", http.StatusBadRequest) // Возвращаем 400, если параметр некорректен.
			return
		}

		// Получаем параметр "to" из URL и преобразуем его в объект time.Time.
		to, err := time.Parse("2006-01-02", r.URL.Query().Get("to"))
		if err != nil {
			http.Error(w, "Invalid param :to:", http.StatusBadRequest) // Возвращаем 400, если параметр некорректен.
			return
		}

		// Получаем параметр "by" из URL для определения типа группировки.
		by := r.URL.Query().Get("by")
		if by != GroupByDay && by != GroupByMonth {
			http.Error(w, "Invalid param :by:", http.StatusBadRequest) // Возвращаем 400, если параметр некорректен.
			return
		}

		// Получаем статистику из репозитория с учетом параметров группировки.
		stats := h.StatRepository.GroupStats(by, from, to)

		// Возвращаем статистику в формате JSON.
		res.Json(w, stats, http.StatusOK)
	}
}
