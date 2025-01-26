package main

import (
	"net/http"
	"rwa/internal/handlers"
	"rwa/internal/service"
	"rwa/internal/storage/users"
)

type App struct {
	H *handlers.Handlers
}

func GetApp() http.Handler {
	a := NewApp()

	return a.H.GetRouter()
}

func NewApp() *App {
	usersStore := users.NewUsersStore()

	svc := service.NewService(usersStore)

	return &App{
		H: handlers.NewHandlers(svc),
	}
}
