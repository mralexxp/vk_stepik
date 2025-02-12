package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
)

const APIURL = "/api"

var NoAuthURL = map[string]string{
	APIURL + "/users/login": http.MethodPost,
	APIURL + "/users":       http.MethodPost,

	APIURL + "/articles": http.MethodGet,
}

// При наличии большго количества подобных URL стоит реализовать префиксное дерево
// example: /api/profiles/ = /api/profiles/*
var NoAuthPrefURL map[string]string = map[string]string{
	APIURL + "/profiles/": http.MethodGet, // /api/profiles/{username}
	APIURL + "/articles/": http.MethodGet, // /api/articles/{slug}
}

func (h *Handlers) endpoints() {
	// Несовместимый метод
	h.router.MethodNotAllowedHandler = http.HandlerFunc(h.MethodNotAllowedHandler)

	h.router.Handle(APIURL+"/users/login", http.HandlerFunc(h.UserLogin)).Methods(http.MethodPost)  // auth user
	h.router.Handle(APIURL+"/users", http.HandlerFunc(h.UserRegister)).Methods(http.MethodPost)     // register new user
	h.router.Handle(APIURL+"/user", http.HandlerFunc(h.UserGet)).Methods(http.MethodGet)            // get current user
	h.router.Handle(APIURL+"/user", http.HandlerFunc(h.UserUpdate)).Methods(http.MethodPut)         // update current user
	h.router.Handle(APIURL+"/user/logout", http.HandlerFunc(h.UserLogout)).Methods(http.MethodPost) // logout current user

	h.router.Handle(APIURL+"/articles", http.HandlerFunc(h.GetArticlesByFilter)).Methods(http.MethodGet)
	h.router.Handle(APIURL+"/articles", http.HandlerFunc(h.CreateArticle)).Methods(http.MethodPost)
	h.router.Handle(APIURL+"/articles/{slug}", http.HandlerFunc(h.GetArticle)).Methods(http.MethodGet)
	h.router.Handle(APIURL+"/articles/{slug}", http.HandlerFunc(h.UpdateArticle)).Methods(http.MethodPut)
	h.router.Handle(APIURL+"/articles/{slug}", http.HandlerFunc(h.DeleteArticle)).Methods(http.MethodDelete)
}

func (h *Handlers) GetRouter() *mux.Router {
	return h.router
}
