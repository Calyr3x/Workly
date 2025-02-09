package app

import (
	"log"
	"net/http"
	"workly/db"
	"workly/routes"
)

const Version = "1.0.0"

func Run() error {
	// Инициализация базы данных
	if err := db.InitDB(); err != nil {
		return err
	}

	// Инициализация зависимостей
	deps, err := InitDependencies()
	if err != nil {
		return err
	}

	// Регистрация маршрутов
	routes.RegisterRoutes(deps.UserHandler, deps.TaskHandler, deps.TeamHandler)

	// Запуск сервера
	log.Printf(">>>>>> Version: %v <<<<<\n", Version)
	log.Println("Server is running on http://localhost:8080")
	return http.ListenAndServe(":8080", nil)
}
