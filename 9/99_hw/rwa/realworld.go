package main

import (
	"net/http"
	"rwa/internal/handlers"
	"rwa/internal/service"
	"rwa/internal/sessions"
	"rwa/internal/storage/profile"
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
	profileStore := profile.NewStore()
	sessionManager := sessions.NewSessionManager()

	svc := service.NewService(usersStore, sessionManager, profileStore)

	return &App{
		H: handlers.NewHandlers(svc),
	}
}
