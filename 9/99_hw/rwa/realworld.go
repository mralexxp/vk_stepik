package main

import (
	"net/http"
	"rwa/internal/handlers"
)

type App struct {
	H *handlers.Handlers
}

func GetApp() http.Handler {
	a := NewApp()

	return a.H.GetRouter()
}

func NewApp() *App {

	return &App{
		H: handlers.NewHandlers(),
	}
}
