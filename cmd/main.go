package main

import (
	"Lessons/configs"            // Пакет для загрузки конфигурации.
	"Lessons/internal/links"     // Модуль для работы с ссылками.
	auth "Lessons/internal/outh" // Модуль аутентификации и авторизации.
	"Lessons/internal/stats"     // Модуль для работы со статистикой.
	"Lessons/internal/user"      // Модуль для работы с пользователями.
	"Lessons/pkg/db"             // Пакет для работы с базой данных.
	"Lessons/pkg/event"          // Пакет для работы с шиной событий (Event Bus).
	"Lessons/pkg/middleware"     // Пакет middleware (например, CORS, логирование).
	"fmt"                        // Стандартный пакет для форматированного вывода.
	"net/http"                   // Стандартный пакет для работы с HTTP.
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
	auth.NewAouthHendler(router, auth.AouthHendlerDeps{
		Config:      conf,        // Передаем конфигурацию в зависимости обработчика.
		AuthService: AuthService, // Сервис для работы с аутентификацией.
	})

	// Инициализация обработчиков для работы с ссылками.
	links.NewLinkHendler(router, links.LinkHendlerDeps{
		LinkRepository: LinkRepository, // Репозиторий для работы с ссылками.
		EventBus:       eventBus,       // Шина событий для обработки событий ссылок.
		Config:         conf,           // Конфигурация для настройки обработчика.
	})

	// Инициализация обработчиков для работы со статистикой.
	stats.NewStatHendler(router, stats.StatHendlerDeps{
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

// main — точка входа в приложение.
func main() {
	// Создаем HTTP-приложение.
	App := App()

	// Настраиваем HTTP-сервер.
	server := http.Server{
		Addr:    ":8080", // Порт, на котором будет запущен сервер.
		Handler: App,     // Указываем маршрутизатор как обработчик запросов.
	}

	// Выводим сообщение о запуске сервера.
	fmt.Println("Listening on port 8080")

	// Запускаем сервер и обрабатываем входящие запросы.
	err := server.ListenAndServe()
	if err != nil {
		// Если произошла ошибка при запуске сервера, завершаем программу.
		return
	}
}
