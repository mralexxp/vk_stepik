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
	RegisterUser(*dto.UserRequest) (*dto.UserResponse, error)
	LoginUser(*dto.UserRequest) (*dto.UserResponse, error)
	GetCurrentUser(string) (*dto.UserResponse, error)
	UpdateUser(*dto.UserRequest) (*dto.UserResponse, error)
}

type Handlers struct {
	router *mux.Router
	Svc    UserServicer
	NoAuth map[string]struct{}
}

func NewHandlers(svc UserServicer) *Handlers {
	h := &Handlers{
		router: mux.NewRouter(),
		Svc:    svc,
		NoAuth: make(map[string]struct{}),
	}

	// Регистрация ручек в роутере
	h.endpoints()

	// Middleware
	h.router.Use(h.ContentTypeMiddleWare)
	h.router.Use(h.AuthMiddleWare)

	return h
}

func (h *Handlers) endpoints() {
	// Несовместимый метод
	h.router.MethodNotAllowedHandler = http.HandlerFunc(h.MethodNotAllowedHandler)

	// Auth handlers
	h.router.Handle(APIURL+"/users/login", http.HandlerFunc(h.UserLogin)).Methods(http.MethodPost) // auth user
	h.router.Handle(APIURL+"/users", http.HandlerFunc(h.UserRegister)).Methods(http.MethodPost)    // register new user
	h.router.Handle(APIURL+"/user", http.HandlerFunc(h.UserGet)).Methods(http.MethodGet)           // get current user
	h.router.Handle(APIURL+"/user", http.HandlerFunc(h.UserUpdate)).Methods(http.MethodPut)        // update current user
}

func (h *Handlers) GetRouter() *mux.Router {
	return h.router
}
