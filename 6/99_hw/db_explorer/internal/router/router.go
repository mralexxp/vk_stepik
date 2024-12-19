package router

import (
	"db_explorer/internal/explorer"
	"net/http"
)

type Router struct {
	Route    map[string]http.Handler
	Explorer *explorer.Explorer
}

func NewRouter(e *explorer.Explorer) http.Handler {
	r := Router{
		Explorer: e,
	}
	return &r
}

// Роутер
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.URL.Query().Get("admin") + "\n"))
}

func (router *Router) endpoints() {
	router.Route["/"] = Handlers.Tables
}
