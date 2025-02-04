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

// LinkHandlerDeps определяет зависимости для LinkHandler.
type LinkHandlerDeps struct {
	LinkRepository *LinkRepository // Репозиторий для работы с сущностями Link.
	EventBus       *event.EventBus // Шина событий для обработки событий, связанных с ссылками.
	Config         *configs.Config // Конфигурация приложения.
}

// LinkHandler отвечает за обработку HTTP-запросов, связанных с Link.
type LinkHandler struct {
	LinkRepository *LinkRepository // Репозиторий для доступа к данным Link.
	EventBus       *event.EventBus // Шина событий для публикации событий.
}

// NewLinkHandler создает новый LinkHandler и регистрирует маршруты.
func NewLinkHandler(router *http.ServeMux, deps LinkHandlerDeps) {
	handler := &LinkHandler{
		LinkRepository: deps.LinkRepository,
		EventBus:       deps.EventBus,
	}

	// Регистрация маршрутов:
	router.Handle("POST /link", middleware.IsAuthenticated(handler.Create(), deps.Config))       // Создание новой ссылки (требует аутентификации).
	router.Handle("PATCH /link/{id}", middleware.IsAuthenticated(handler.Update(), deps.Config)) // Обновление существующей ссылки (требует аутентификации).
	router.HandleFunc("DELETE /link/{id}", handler.Delete())                                     // Удаление ссылки.
	router.HandleFunc("GET /link/{hash}", handler.GoTo())                                        // Переход по сокращенной ссылке.
	router.Handle("GET /link/", middleware.IsAuthenticated(handler.GetAll(), deps.Config))       // Получение списка всех ссылок (требует аутентификации).
}

// Create обрабатывает создание новой ссылки.
func (handler *LinkHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем тело запроса и преобразуем в структуру LinkCreateRequest.
		body, err := reg.HandleBody[LinkCreateRequest](&w, r)
		if err != nil {
			return // Если ошибка, обработка прекращается.
		}

		link := NewLink(body.Url) // Создаем новую ссылку на основе переданного URL.

		// Генерируем уникальный хеш для ссылки.
		for {
			existedLink, _ := handler.LinkRepository.GetByHash(link.Hash) // Проверяем, существует ли ссылка с таким хешем в БД.
			if existedLink == nil {
				break // Если ссылка не существует, выходим из цикла.
			}
			link.GenereateHash() // Если хеш уже занят, генерируем новый.
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

// Update обрабатывает обновление существующей ссылки.
func (handler *LinkHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := r.Context().Value(middleware.ContextEmailKey).(string)
		if ok {
			fmt.Println(email) // Логируем email пользователя (если он есть в контексте).
		}

		// Извлекаем тело запроса и преобразуем в структуру LinkUpdateRequest.
		body, err := reg.HandleBody[LinkUpdateRequest](&w, r)
		if err != nil {
			return // Если ошибка, обработка прекращается.
		}

		// Получаем значение параметра {id} из URL.
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) // Возвращаем 400, если ID некорректен.
			return
		}

		// Обновляем ссылку в репозитории.
		link, err := handler.LinkRepository.Update(&Link{
			Model: gorm.Model{ID: uint(id)}, // ID ссылки.
			Url:   body.Url,                 // Новый URL.
			Hash:  body.Hash,                // Новый хеш.
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) // Возвращаем 400 в случае ошибки.
			return
		}

		// Возвращаем обновленную ссылку в формате JSON.
		res.Json(w, link, http.StatusOK)
	}
}

// Delete обрабатывает удаление ссылки по ID.
func (handler *LinkHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем значение параметра {id} из URL.
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) // Возвращаем 400, если ID некорректен.
			return
		}

		// Проверяем, существует ли ссылка с указанным ID.
		_, err = handler.LinkRepository.GetById(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) // Возвращаем 400, если ссылка не найдена.
			return
		}

		// Удаляем ссылку из базы данных.
		err = handler.LinkRepository.Database.Model(&Link{}).Delete(&Link{}, id).Error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // Возвращаем 500 в случае ошибки.
			return
		}

		// Возвращаем статус успешного удаления.
		res.Json(w, "status: deleted", http.StatusOK)
	}
}

// GoTo обрабатывает переход по сокращенной ссылке.
func (handler *LinkHandler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash") // Получаем хеш из URL.

		// Ищем ссылку по хешу в репозитории.
		link, err := handler.LinkRepository.GetByHash(hash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // Возвращаем 500 в случае ошибки.
			return
		}

		// Публикуем событие о посещении ссылки.
		go handler.EventBus.Publish(event.Event{
			Type: event.EventLinkVisited, // Тип события: посещение ссылки.
			Data: link.ID,                // ID ссылки.
		})

		// Перенаправляем пользователя на оригинальный URL.
		http.Redirect(w, r, link.Url, http.StatusTemporaryRedirect)
	}
}

// GetAll возвращает список всех ссылок с пагинацией.
func (handler *LinkHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем параметр limit из запроса.
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			http.Error(w, "Invalid limit", http.StatusBadRequest) // Возвращаем 400, если limit некорректен.
			return
		}

		// Получаем параметр offset из запроса.
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			http.Error(w, "Invalid offset", http.StatusBadRequest) // Возвращаем 400, если offset некорректен.
			return
		}

		// Получаем список ссылок с учетом пагинации.
		links := handler.LinkRepository.GetAll(limit, offset)
		count := handler.LinkRepository.Count() // Получаем общее количество ссылок.

		// Возвращаем результат в формате JSON.
		res.Json(w, GetAllResponse{
			Links: links, // Список ссылок.
			Count: count, // Общее количество ссылок.
		}, http.StatusOK)
	}
}
