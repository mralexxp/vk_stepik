package handlers

import (
	"github.com/gorilla/mux"
	"rwa/internal/service"
)

type Handlers struct {
	router     *mux.Router
	Svc        *service.Service
	NoAuth     map[string]string
	NoAuthPref map[string]string
}

func NewHandlers(svc *service.Service) *Handlers {
	h := &Handlers{
		router:     mux.NewRouter(),
		Svc:        svc,
		NoAuth:     NoAuthURL,
		NoAuthPref: NoAuthPrefURL,
	}

	h.endpoints()

	h.router.Use(h.ContentTypeMiddleWare)
	h.router.Use(h.AuthMiddleWare)

	return h
}
