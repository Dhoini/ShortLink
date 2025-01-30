package stats

import (
	"Lessons/pkg/event"
	"log"
)

// StatServiceDeps определяет зависимости для сервиса статистики.
type StatServiceDeps struct {
	EventBus       *event.EventBus // Шина событий для обработки событий, связанных со статистикой.
	StatRepository *StatRepository // Репозиторий для работы с данными статистики.
}

// StatService представляет сервис для работы со статистикой.
type StatService struct {
	EventBus       *event.EventBus // Шина событий для подписки на события.
	StatRepository *StatRepository // Репозиторий для доступа к данным статистики.
}

// NewStatService создает новый экземпляр StatService.
func NewStatService(deps *StatServiceDeps) *StatService {
	return &StatService{
		EventBus:       deps.EventBus,       // Инициализация шины событий.
		StatRepository: deps.StatRepository, // Инициализация репозитория статистики.
	}
}

// AddClick обрабатывает события посещения ссылок и увеличивает счетчик кликов.
func (s *StatService) AddClick() {
	// Подписываемся на события через EventBus.
	for msg := range s.EventBus.Subscribe() {
		// Проверяем, является ли событие типом EventLinkVisited (посещение ссылки).
		if msg.Type == event.EventLinkVisited {
			// Пробуем преобразовать данные события в uint (ID ссылки).
			id, ok := msg.Data.(uint)
			if !ok {
				// Если преобразование не удалось, логируем ошибку и продолжаем обработку.
				log.Fatalln("Bad EventLinkVisited Data: ", msg.Data)
				continue
			}

			// Увеличиваем счетчик кликов для указанной ссылки.
			s.StatRepository.AddClick(id)
		}
	}
}
