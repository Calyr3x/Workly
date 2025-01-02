package main

import (
	"log"
	"net/http"
	"workly/repository"
	"workly/usecase"

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

	repo := repository.NewUserRepository(db.DB)
	uc := usecase.NewUserUseCase(repo)
	handler := handlers.NewUserHandler(uc)

	repo1 := repository.NewTaskRepository(db.DB)
	uc1 := usecase.NewTaskUseCase(repo1)
	handler1 := handlers.NewTaskHandler(uc1)

	// Регистрация обработчиков с использованием middleware CORS
	http.HandleFunc("/login", middleware.WithCORS(handler.Login))
	http.HandleFunc("/register", middleware.WithCORS(handler.Register))
	http.HandleFunc("/tasks", middleware.WithCORS(handler1.GetTasks))
	http.HandleFunc("/tasks/create", middleware.WithCORS(handler1.CreateTask))
	http.HandleFunc("/tasks/update", middleware.WithCORS(handler1.UpdateTask))
	http.HandleFunc("/tasks/delete", middleware.WithCORS(handler1.DeleteTask))
	http.HandleFunc("/tasks/", middleware.WithCORS(handler1.GetTaskByID))
	http.HandleFunc("/task_access", middleware.WithCORS(handler1.CreateTaskAccess))
	http.HandleFunc("/updateAvatar", middleware.WithCORS(handler.UpdateAvatar))
	http.HandleFunc("/getUserData", middleware.WithCORS(handler.GetUserData))
	http.HandleFunc("/updateUsername", middleware.WithCORS(handler.UpdateUsername))
	http.HandleFunc("/createTeam", middleware.WithCORS(handlers.HandleCreateTeam))
	http.HandleFunc("/getUserAvatar", middleware.WithCORS(handlers.HandleGetUserAvatar))
	http.HandleFunc("/getTeams", middleware.WithCORS(handlers.HandleGetTeams))
	http.HandleFunc("/removeMember", middleware.WithCORS(handlers.HandleRemoveMember))
	http.HandleFunc("/addMember", middleware.WithCORS(handlers.HandleAddMember))
	http.HandleFunc("/getUserIds", middleware.WithCORS(handler.GetUserIDs))

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
