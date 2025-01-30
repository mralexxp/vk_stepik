package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"rwa/internal/dto"
)

const (
	APIURL = "/api"
)

type UserServicer interface {
	Register(*dto.UserRegisterRequest) (*dto.UserRegisterResponse, error)
	Login(*dto.UserLoginRequest) (*dto.UserLoginResponse, error)
}

type Handlers struct {
	router *mux.Router
	Svc    UserServicer
}

func NewHandlers(svc UserServicer) *Handlers {
	h := &Handlers{
		router: mux.NewRouter(),
		Svc:    svc,
	}

	// Регистрация ручек в роутере
	h.endpoints()

	// Middleware
	h.router.Use(h.ContentTypeMiddleWare)

	return h
}

func (h *Handlers) endpoints() {
	// Несовместимый метод
	h.router.MethodNotAllowedHandler = http.HandlerFunc(h.MethodNotAllowedHandler)

	// Auth handlers
	h.router.Handle(APIURL+"/users/login", http.HandlerFunc(h.UserLogin)).Methods(http.MethodPost) // auth user
	h.router.Handle(APIURL+"/users", http.HandlerFunc(h.UserRegister)).Methods(http.MethodPost)    // register new user
	h.router.Handle(APIURL+"/users", http.HandlerFunc(h.UserGet)).Methods(http.MethodGet)          // get current user
	h.router.Handle(APIURL+"/users", http.HandlerFunc(h.UserUpdate)).Methods(http.MethodPut)       // update current user
}

func (h *Handlers) GetRouter() *mux.Router {
	return h.router
}
