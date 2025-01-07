package routes

import (
	"net/http"
	"workly/handlers"
	"workly/middleware"
)

// RegisterRoutes регистрирует все маршруты приложения
func RegisterRoutes(userHandler *handlers.UserHandler, taskHandler *handlers.TaskHandler, teamHandler *handlers.TeamHandler) {
	http.HandleFunc("/login", middleware.WithCORS(userHandler.Login))
	http.HandleFunc("/register", middleware.WithCORS(userHandler.Register))
	http.HandleFunc("/tasks", middleware.WithCORS(taskHandler.GetTasks))
	http.HandleFunc("/tasks/create", middleware.WithCORS(taskHandler.CreateTask))
	http.HandleFunc("/tasks/update", middleware.WithCORS(taskHandler.UpdateTask))
	http.HandleFunc("/tasks/delete", middleware.WithCORS(taskHandler.DeleteTask))
	http.HandleFunc("/tasks/", middleware.WithCORS(taskHandler.GetTaskByID))
	http.HandleFunc("/task_access", middleware.WithCORS(taskHandler.CreateTaskAccess))
	http.HandleFunc("/updateAvatar", middleware.WithCORS(userHandler.UpdateAvatar))
	http.HandleFunc("/getUserData", middleware.WithCORS(userHandler.GetUserData))
	http.HandleFunc("/updateUsername", middleware.WithCORS(userHandler.UpdateUsername))
	http.HandleFunc("/createTeam", middleware.WithCORS(teamHandler.CreateTeam))
	http.HandleFunc("/getUserAvatar", middleware.WithCORS(teamHandler.GetUserAvatar))
	http.HandleFunc("/getTeams", middleware.WithCORS(teamHandler.GetTeams))
	http.HandleFunc("/removeMember", middleware.WithCORS(teamHandler.RemoveMember))
	http.HandleFunc("/addMember", middleware.WithCORS(teamHandler.AddMember))
	http.HandleFunc("/getUserIds", middleware.WithCORS(userHandler.GetUserIDs))
}
