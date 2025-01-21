package main

import (
	"Lessons/configs"
	"Lessons/internal/links"
	"Lessons/internal/outh"
	"Lessons/internal/user"
	"Lessons/pkg/db"
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

	// Инициализация репозиториев.
	LinkRepository := links.NewLinkRepository(dB)
	UserRepository := user.NewUserRepository(dB)

	//Services
	AuthService := auth.NewUserService(UserRepository)

	// Инициализация обработчиков для маршрутов аутентификации.
	auth.NewAouthHendler(router, auth.AouthHendlerDeps{
		Config:      conf, // Передаем конфигурацию в зависимости обработчика
		AuthService: AuthService,
	})

	// Инициализация обработчиков для работы с сущностями Link.
	links.NewLinkHendler(router, links.LinkHendlerDeps{
		LinkRepository: LinkRepository, // Передаем репозиторий ссылок в зависимости обработчика.
	})
	//middlwares
	stackMiddlwares := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
		middleware.IsAuthenticated,
	)
	// Настраиваем сервер.
	server := http.Server{
		Addr:    ":8080",                 // Порт, на котором будет запущен сервер.
		Handler: stackMiddlwares(router), // Указываем маршрутизатор как обработчик.
	}

	// Выводим сообщение о запуске сервера.
	fmt.Println("Listening on port 8080")

	// Запускаем сервер и обрабатываем входящие запросы.
	server.ListenAndServe()
}
