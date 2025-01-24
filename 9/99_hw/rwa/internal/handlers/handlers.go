package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Handlers struct {
	router *mux.Router
}

func NewHandlers() *Handlers {
	h := &Handlers{
		router: mux.NewRouter(),
	}

	// MethodNotAllowedPage
	h.router.MethodNotAllowedHandler = http.HandlerFunc(h.MethodNotAllowedHandler)

	// Auth handlers
	h.router.Handle("/users/login", http.HandlerFunc(h.UsersLogin)).Methods(http.MethodPost) // auth user
	h.router.Handle("/users/", http.HandlerFunc(h.UsersRegister)).Methods(http.MethodPost)   // register new user
	h.router.Handle("/users/", http.HandlerFunc(h.UsersGet)).Methods(http.MethodGet)         // get current user
	h.router.Handle("/users/", http.HandlerFunc(h.UsersUpdate)).Methods(http.MethodPut)      // update current user

	return h
}

func (h *Handlers) GetRouter() *mux.Router {
	return h.router
}
