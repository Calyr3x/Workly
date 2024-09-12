package main

import (
	"log"
	"net/http"

	"workly/config"
	"workly/db"
	"workly/handlers"
	"workly/middleware"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatal("Failed to load configuration: ", err)
	}

	// Инициализация базы данных
	db.InitDB(cfg.DB)

	// Регистрация обработчиков с использованием middleware CORS
	http.HandleFunc("/login", middleware.WithCORS(handlers.HandleLogin))
	http.HandleFunc("/register", middleware.WithCORS(handlers.HandleRegister))
	http.HandleFunc("/tasks", middleware.WithCORS(handlers.HandleTasks))
	http.HandleFunc("/tasks/create", middleware.WithCORS(handlers.HandleCreateTask))
	http.HandleFunc("/tasks/update", middleware.WithCORS(handlers.HandleUpdateTask))
	http.HandleFunc("/tasks/delete", middleware.WithCORS(handlers.HandleDeleteTask))
	http.HandleFunc("/tasks/", middleware.WithCORS(handlers.HandleGetTaskByID))
	http.HandleFunc("/updateAvatar", handlers.WithCORS(handlers.HandleUpdateAvatar))
	http.HandleFunc("/getUserData", handlers.WithCORS(handlers.HandleGetUserData))
	http.HandleFunc("/updateUsername", handlers.WithCORS(handlers.HandleUpdateUsername))
	http.HandleFunc("/createTeam", handlers.WithCORS(handlers.HandleCreateTeam))
	http.HandleFunc("/getUserAvatar", handlers.WithCORS(handlers.HandleGetUserAvatar))

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
