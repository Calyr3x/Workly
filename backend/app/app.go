package app

import (
	"log"
	"net/http"
	"workly/config"
	"workly/db"
	"workly/routes"
)

func Run() error {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig("config/config.json")
	if err != nil {
		return err
	}

	// Инициализация базы данных
	if err := db.InitDB(cfg.DB); err != nil {
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
	log.Println("Server is running on http://localhost:8080")
	return http.ListenAndServe(":8080", nil)
}
