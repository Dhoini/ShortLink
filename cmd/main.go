package main

import (
	"fmt"
	"net/http"
)

func main() {
	app := App()
	// Настраиваем сервер.
	server := http.Server{
		Addr:    ":8080", // Порт, на котором будет запущен сервер.
		Handler: app,     // Указываем маршрутизатор как обработчик.
	}
	// Выводим сообщение о запуске сервера.
	fmt.Println("Listening on port 8080")

	// Запускаем сервер и обрабатываем входящие запросы.
	err := server.ListenAndServe()
	if err != nil {
		return
	}
}
