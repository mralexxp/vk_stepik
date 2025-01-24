package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
)

const (
	APIURL = "/api"
)

type Handlers struct {
	router *mux.Router
}

func NewHandlers() *Handlers {
	h := &Handlers{
		router: mux.NewRouter(),
	}

	// Регистрация ручек в роутере
	h.endpoints()

	return h
}

func (h *Handlers) endpoints() {
	// MethodNotAllowedPage
	h.router.MethodNotAllowedHandler = http.HandlerFunc(h.MethodNotAllowedHandler)

	// Auth handlers
	h.router.Handle(APIURL+"/users/login", http.HandlerFunc(h.UsersLogin)).Methods(http.MethodPost) // auth user
	h.router.Handle(APIURL+"/users", http.HandlerFunc(h.UsersRegister)).Methods(http.MethodPost)    // register new user
	h.router.Handle(APIURL+"/users", http.HandlerFunc(h.UsersGet)).Methods(http.MethodGet)          // get current user
	h.router.Handle(APIURL+"/users", http.HandlerFunc(h.UsersUpdate)).Methods(http.MethodPut)       // update current user
}

func (h *Handlers) GetRouter() *mux.Router {
	return h.router
}
