package router

import (
	"db_explorer/internal/explorer"
	"fmt"
	"net/http"
)

type Router struct {
	Route    map[string]func(http.ResponseWriter, *http.Request)
	Explorer *explorer.Explorer
}

func NewRouter(e *explorer.Explorer) http.Handler {
	r := Router{
		Explorer: e,
	}
	h := &Handlers{
		Router: r,
	}
	r.endpoints(h)
	return &r
}

// Роутер
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	TempMiddleWare(&w)
	fmt.Println(r.URL.Path)
	if fn, ok := router.Route[r.URL.Path]; ok {
		fn(w, r)
		return
	}

	http.NotFound(w, r)
}

func (router *Router) endpoints(h *Handlers) {
	router.Route = make(map[string]func(http.ResponseWriter, *http.Request))

	// Главная со списком таблиц
	router.Route["/"] = func(w http.ResponseWriter, r *http.Request) { h.Index(w, r) }

	// Создаем статические пути к нашим таблицам
	tables, err := router.Explorer.ShowTables()
	if err != nil {
		panic(err)
	}

}
