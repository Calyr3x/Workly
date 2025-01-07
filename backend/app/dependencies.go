package app

import (
	"workly/db"
	"workly/handlers"
	"workly/repository"
	"workly/usecase"
)

type Dependencies struct {
	UserHandler *handlers.UserHandler
	TaskHandler *handlers.TaskHandler
	TeamHandler *handlers.TeamHandler
}

func InitDependencies() (*Dependencies, error) {

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(db.DB)
	taskRepo := repository.NewTaskRepository(db.DB)
	teamRepo := repository.NewTeamRepository(db.DB)

	// Инициализация юзкейсов
	userUC := usecase.NewUserUseCase(userRepo)
	taskUC := usecase.NewTaskUseCase(taskRepo)
	teamUC := usecase.NewTeamUseCase(teamRepo)

	// Инициализация хендлеров
	userHandler := handlers.NewUserHandler(userUC)
	taskHandler := handlers.NewTaskHandler(taskUC)
	teamHandler := handlers.NewTeamHandler(teamUC)

	return &Dependencies{
		UserHandler: userHandler,
		TaskHandler: taskHandler,
		TeamHandler: teamHandler,
	}, nil
}
