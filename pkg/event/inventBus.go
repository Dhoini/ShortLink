package event

// Константы для типов событий.
const (
	EventLinkVisited = "link-visited" // Событие, которое происходит при посещении ссылки.
)

// Event представляет собой структуру события.
type Event struct {
	Type string // Тип события (например, "link-visited").
	Data any    // Данные, связанные с событием. Может быть любого типа.
}

// EventBus представляет шину событий, которая позволяет публиковать и подписываться на события.
type EventBus struct {
	bus chan Event // Канал, через который передаются события.
}

// NewEventBus создает новый экземпляр EventBus.
func NewEventBus() *EventBus {
	return &EventBus{
		bus: make(chan Event), // Инициализируем канал для передачи событий.
	}
}

// Publish публикует событие в шину событий.
// event: событие, которое нужно опубликовать.
func (e *EventBus) Publish(event Event) {
	e.bus <- event // Отправляем событие в канал.
}

// Subscribe возвращает канал, через который можно получать события.
// Возвращаемый канал является только для чтения (<-chan).
func (e *EventBus) Subscribe() <-chan Event {
	return e.bus // Возвращаем канал для подписки на события.
}
