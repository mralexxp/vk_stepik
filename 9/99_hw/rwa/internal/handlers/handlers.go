package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"rwa/internal/dto"
	"rwa/internal/service"
)

const APIURL = "/api"

var NoAuthURL map[string]string = map[string]string{
	// users
	APIURL + "/users/login": http.MethodPost,
	APIURL + "/users":       http.MethodPost,

	// articles
	APIURL + "/articles": http.MethodGet,
}

// При наличии большго количества подобных URL стоит реализовать префиксное дерево
// example: /api/profiles/ = /api/profiles/*
var NoAuthPrefURL map[string]string = map[string]string{
	APIURL + "/profiles/": http.MethodGet, // /api/profiles/{username}
	APIURL + "/articles/": http.MethodGet, // /api/articles/{slug}
}

type UserServicer interface {
	RegisterUser(*dto.UserRequest) (*dto.UserResponse, error)
	LoginUser(*dto.UserRequest) (*dto.UserResponse, error)
	GetCurrentUser(string) (*dto.UserResponse, error)
	UpdateUser(*dto.UserRequest) (*dto.UserResponse, error)

	GetSessionManager() service.SessManager
}

type Handlers struct {
	router     *mux.Router
	Svc        UserServicer
	NoAuth     map[string]string
	NoAuthPref map[string]string
}

func NewHandlers(svc UserServicer) *Handlers {
	h := &Handlers{
		router:     mux.NewRouter(),
		Svc:        svc,
		NoAuth:     NoAuthURL,
		NoAuthPref: NoAuthPrefURL,
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

	// Article handlers
	h.router.Handle(APIURL+"/articles", http.HandlerFunc(h.GetAllArticles)).Methods(http.MethodGet)
	h.router.Handle(APIURL+"/articles/feed", http.HandlerFunc(h.GetFeedArticles)).Methods(http.MethodGet)
	h.router.Handle(APIURL+"/articles", http.HandlerFunc(h.CreateArticle)).Methods(http.MethodPost)
	h.router.Handle(APIURL+"/articles/{slug}", http.HandlerFunc(h.GetArticle)).Methods(http.MethodGet)
	h.router.Handle(APIURL+"/articles/{slug}", http.HandlerFunc(h.UpdateArticle)).Methods(http.MethodPut)
	h.router.Handle(APIURL+"/articles/{slug}", http.HandlerFunc(h.DeleteArticle)).Methods(http.MethodDelete)
}

func (h *Handlers) GetRouter() *mux.Router {
	return h.router
}
