package main

import (
	"Lessons/configs"
	"Lessons/internal/links"
	auth "Lessons/internal/outh"
	"Lessons/internal/stats"
	"Lessons/internal/user"
	"Lessons/pkg/db"
	"Lessons/pkg/event"
	"Lessons/pkg/middleware"
	"net/http"
)

// App создает и настраивает HTTP-приложение.
func App() http.Handler {
	// Загружаем конфигурацию из файла или переменных окружения.
	conf := configs.LoadConfig()

	// Инициализируем подключение к базе данных, используя загруженную конфигурацию.
	dB := db.NewDb(conf)

	// Создаем маршрутизатор для обработки HTTP-запросов.
	router := http.NewServeMux()

	// Инициализируем шину событий (Event Bus) для обработки событий в приложении.
	eventBus := event.NewEventBus()

	// Инициализация репозиториев для работы с базой данных.
	LinkRepository := links.NewLinkRepository(dB) // Репозиторий для работы с ссылками.
	UserRepository := user.NewUserRepository(dB)  // Репозиторий для работы с пользователями.
	StatRepository := stats.NewStatRepository(dB) // Репозиторий для работы со статистикой.

	// Инициализация сервисов, которые используют репозитории и другие зависимости.
	AuthService := auth.NewUserService(UserRepository) // Сервис для работы с аутентификацией пользователей.
	statService := stats.NewStatService(&stats.StatServiceDeps{
		EventBus:       eventBus,       // Шина событий для обработки событий статистики.
		StatRepository: StatRepository, // Репозиторий для хранения данных статистики.
	})

	// Инициализация обработчиков для маршрутов аутентификации.
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,        // Передаем конфигурацию в зависимости обработчика.
		AuthService: AuthService, // Сервис для работы с аутентификацией.
	})

	// Инициализация обработчиков для работы с ссылками.
	links.NewLinkHandler(router, links.LinkHandlerDeps{
		LinkRepository: LinkRepository, // Репозиторий для работы с ссылками.
		EventBus:       eventBus,       // Шина событий для обработки событий ссылок.
		Config:         conf,           // Конфигурация для настройки обработчика.
	})

	// Инициализация обработчиков для работы со статистикой.
	stats.NewStatHandler(router, stats.StatHandlerDeps{
		StatRepository: StatRepository, // Репозиторий для работы со статистикой.
		Config:         conf,           // Конфигурация для настройки обработчика.
	})

	// Запускаем фоновую задачу для обработки событий статистики.
	go statService.AddClick()

	// Настраиваем цепочку middleware для обработки запросов.
	stackMiddlwares := middleware.Chain(
		middleware.CORS,    // Middleware для обработки CORS.
		middleware.Logging, // Middleware для логирования запросов.
	)

	// Возвращаем маршрутизатор, обернутый в middleware.
	return stackMiddlwares(router)
}
