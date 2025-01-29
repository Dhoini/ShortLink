package stats

import (
	"Lessons/configs"
	"Lessons/pkg/middleware"
	"Lessons/pkg/res"
	"net/http"
	"time"
)

const (
	GroupByDay   = "day"
	GroupByMonth = "month"
)

type StatHendlerDeps struct {
	StatRepository *StatRepository
	Config         *configs.Config // Репозиторий для работы с сущностями Link.
}

// LinkHendler отвечает за обработку HTTP-запросов, связанных с Link.
type StatHandler struct {
	StatRepository *StatRepository // Репозиторий для доступа к данным Link.
}

// NewLinkHendler создает новый LinkHendler и регистрирует маршруты.
// router: маршрутизатор для обработки HTTP-запросов.
// deps: зависимости для LinkHendler.
func NewStatHendler(router *http.ServeMux, deps StatHendlerDeps) {
	handler := &StatHandler{
		StatRepository: deps.StatRepository,
	}
	// Регистрация маршрутов. 	// Переход по сокращенной ссылке.
	router.Handle("GET /stat", middleware.IsAuthenticated(handler.GetStat(), deps.Config))
}
func (h *StatHandler) GetStat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		from, err := time.Parse("2006-01-02", r.URL.Query().Get("from"))
		if err != nil {
			http.Error(w, "Invalid param :from:", http.StatusBadRequest)
			return
		}

		to, err := time.Parse("2006-01-02", r.URL.Query().Get("to"))
		if err != nil {
			http.Error(w, "invalid param :to:", http.StatusBadRequest)
			return
		}
		by := r.URL.Query().Get("by")
		if by != GroupByDay && by != GroupByMonth {
			http.Error(w, "Invalid param :by:", http.StatusBadRequest)
			return
		}
		stats := h.StatRepository.GroupStats(by, from, to)
		res.Json(w, stats, http.StatusOK)
	}
}
