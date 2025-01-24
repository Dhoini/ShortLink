package main

import (
	"Lessons/configs"
	"Lessons/internal/links"
	"Lessons/internal/outh"
	"Lessons/internal/stats"
	"Lessons/internal/user"
	"Lessons/pkg/db"
	"Lessons/pkg/event"
	"Lessons/pkg/middleware"
	"fmt"
	"net/http"
)

func main() {
	// Загружаем конфигурацию из файла или среды.
	conf := configs.LoadConfig()

	// Инициализируем подключение к базе данных.
	dB := db.NewDb(conf)

	// Создаем маршрутизатор для обработки HTTP-запросов.
	router := http.NewServeMux()
	eventBus := event.NewEventBus()

	// Инициализация репозиториев.
	LinkRepository := links.NewLinkRepository(dB)
	UserRepository := user.NewUserRepository(dB)
	StatRepository := stats.NewStatRepository(dB)
	//Services
	AuthService := auth.NewUserService(UserRepository)
	statService := stats.NewStatService(&stats.StatServiceDeps{
		EventBus:       eventBus,
		StatRepository: StatRepository,
	})
	// Инициализация обработчиков для маршрутов аутентификации.
	auth.NewAouthHendler(router, auth.AouthHendlerDeps{
		Config:      conf, // Передаем конфигурацию в зависимости обработчика
		AuthService: AuthService,
	})

	// Инициализация обработчиков для работы с сущностям.
	links.NewLinkHendler(router, links.LinkHendlerDeps{
		LinkRepository: LinkRepository, // Передаем репозиторий ссылок в зависимости обработчика.
		EventBus:       eventBus,
		Config:         conf,
	})

	//middlwares
	stackMiddlwares := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)
	// Настраиваем сервер.
	server := http.Server{
		Addr:    ":8080",                 // Порт, на котором будет запущен сервер.
		Handler: stackMiddlwares(router), // Указываем маршрутизатор как обработчик.
	}

	go statService.AddClick()
	// Выводим сообщение о запуске сервера.
	fmt.Println("Listening on port 8080")

	// Запускаем сервер и обрабатываем входящие запросы.
	err := server.ListenAndServe()
	if err != nil {
		return
	}
}
