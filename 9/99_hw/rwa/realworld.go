package main

import (
	"net/http"
	"rwa/internal/handlers"
)

// сюда писать код

func GetApp() http.Handler {
	h := handlers.NewHandlers()

	return h.GetRouter()
}
