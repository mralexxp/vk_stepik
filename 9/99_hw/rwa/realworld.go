package main

import (
	"net/http"
	"rwa/internal/handlers"
	"rwa/internal/service"
	"rwa/internal/sessions"
	"rwa/internal/storage/articles"
	"rwa/internal/storage/profile"
	"rwa/internal/storage/users"
)

type App struct {
	H *handlers.Handlers
	S *service.Service
}

func GetApp() http.Handler {
	a := NewApp()

	return a.H.GetRouter()
}

func NewApp() *App {
	usersStore := users.NewUsersStore()
	profileStore := profile.NewStore()
	articlesStore := articles.NewStore()
	sessionManager := sessions.NewSessionManager()

	svc := service.NewService(articlesStore, profileStore, usersStore, sessionManager)

	return &App{
		H: handlers.NewHandlers(svc),
		S: svc,
	}
}
