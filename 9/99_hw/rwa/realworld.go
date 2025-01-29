package main

import (
	"net/http"
	"rwa/internal/handlers"
	"rwa/internal/service"
	"rwa/internal/sessions"
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
	sessionManager := sessions.NewSessionManager()

	svc := service.NewService(usersStore, sessionManager)

	return &App{
		H: handlers.NewHandlers(svc),
	}
}
