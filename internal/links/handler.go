package links

import (
	"Lessons/configs"
	"Lessons/pkg/event"
	"Lessons/pkg/middleware"
	"Lessons/pkg/reg"
	"Lessons/pkg/res"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// LinkHendlerDeps определяет зависимости для LinkHendler.
type LinkHendlerDeps struct {
	LinkRepository *LinkRepository
	EventBus       *event.EventBus
	Config         *configs.Config // Репозиторий для работы с сущностями Link.
}

// LinkHendler отвечает за обработку HTTP-запросов, связанных с Link.
type LinkHendler struct {
	LinkRepository *LinkRepository // Репозиторий для доступа к данным Link.
	EventBus       *event.EventBus
}

// NewLinkHendler создает новый LinkHendler и регистрирует маршруты.
// router: маршрутизатор для обработки HTTP-запросов.
// deps: зависимости для LinkHendler.
func NewLinkHendler(router *http.ServeMux, deps LinkHendlerDeps) {
	handler := &LinkHendler{
		LinkRepository: deps.LinkRepository,
		EventBus:       deps.EventBus,
	}
	// Регистрация маршрутов.
	router.HandleFunc("POST /link", handler.Create())                                            // Создание новой ссылки.
	router.Handle("PATCH /link/{id}", middleware.IsAuthenticated(handler.Update(), deps.Config)) // Обновление существующей ссылки.
	router.HandleFunc("DELETE /link/{id}", handler.Delete())                                     // Удаление ссылки.
	router.HandleFunc("GET /link/{hash}", handler.GoTo())                                        // Переход по сокращенной ссылке.
	router.Handle("GET /link/", middleware.IsAuthenticated(handler.GetAll(), deps.Config))
}

// Create обрабатывает создание новой ссылки.
func (handler *LinkHendler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем тело запроса и преобразуем в LinkCreateRequest.
		body, err := reg.HandleBody[LinkCreateRequest](&w, r)
		if err != nil {
			return // Если ошибка, обработка прекращается.
		}

		link := NewLink(body.Url) //создаем ссылку

		for {
			existedLink, _ := handler.LinkRepository.GetByHash(link.Hash) //проверка на то что ссылка уже существует в БД
			//если не существует, выходим из цикла или создаем новую
			if existedLink == nil {
				break
			}
			link.GenereateHash()
		}

		// Сохраняем новую ссылку в репозитории.
		createdLink, err := handler.LinkRepository.Create(link)
		if err != nil {

			http.Error(w, err.Error(), http.StatusInternalServerError) // Возвращаем 500 в случае ошибки.
			return
		}

		// Возвращаем созданную ссылку в формате JSON.
		res.Json(w, createdLink, http.StatusOK)
	}
}

// Update обрабатывает обновление существующей ссылки (пока пустая реализация).
func (handler *LinkHendler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := r.Context().Value(middleware.ContextEmailKey).(string)
		if ok {
			fmt.Println(email)
		}
		// Здесь реализация обновления ссылки.
		body, err := reg.HandleBody[LinkUpdateRequest](&w, r)
		if err != nil {
			return // Если ошибка, обработка прекращается.
		}
		// Получаем значение параметра {id} из URL.
		idString := r.PathValue("id")

		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		link, err := handler.LinkRepository.Update(&Link{
			Model: gorm.Model{ID: uint(id)},
			Url:   body.Url,
			Hash:  body.Hash,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, link, http.StatusOK)

	}
}

// Delete обрабатывает удаление ссылки по ID.
func (handler *LinkHendler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем значение параметра {id} из URL.
		idString := r.PathValue("id")

		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = handler.LinkRepository.GetById(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = handler.LinkRepository.Database.Model(&Link{}).Delete(&Link{}, id).Error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Json(w, "status: deleted", http.StatusOK)
	}
}

// GoTo обрабатывает переход по сокращенной ссылке
func (handler *LinkHendler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")

		link, err := handler.LinkRepository.GetByHash(hash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // Возвращаем 500 в случае ошибки.
			return
		}

		go handler.EventBus.Publish(event.Event{
			Type: event.EventLinkVisited,
			Data: link.ID,
		})

		http.Redirect(w, r, link.Url, http.StatusTemporaryRedirect)

	}
}

func (handler *LinkHendler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}

		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			http.Error(w, "Invalid offset", http.StatusBadRequest)
			return
		}

		links := handler.LinkRepository.GetAll(limit, offset)
		count := handler.LinkRepository.Count()
		res.Json(w, GetAllResponse{
			Links: links,
			Count: count},
			http.StatusOK)
	}
}
