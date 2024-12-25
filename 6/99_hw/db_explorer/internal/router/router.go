package router

import (
	"db_explorer/internal/errors"
	"db_explorer/internal/explorer"
	"fmt"
	"net/http"
	"strings"
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
	target := r.Method + r.URL.Path
	fmt.Println(target)
	splitPath := strings.Split(target, "/")

	// Route с ID
	if len(splitPath) == 3 && splitPath[len(splitPath)-1] != "" {
		target = strings.Join(splitPath[:len(splitPath)-1], "/") + "/id"
	}

	if fn, ok := router.Route[target]; ok {
		fn(w, r)
		return
	}

	errors.SendJSONError(w, http.StatusNotFound, "unknown table")
}

func (router *Router) endpoints(h *Handlers) {
	router.Route = make(map[string]func(http.ResponseWriter, *http.Request))

	// Главная со списком таблиц
	router.Route["GET/"] = func(w http.ResponseWriter, r *http.Request) { h.Index(w, r) }

	// Создаем статические пути к нашим таблицам
	tables, err := router.Explorer.ShowTables()
	if err != nil {
		panic(err)
	}

	for _, table := range tables {
		router.Route[fmt.Sprintf("GET/%s", table)] = h.GetTableValues
		router.Route[fmt.Sprintf("GET/%s/", table)] = h.GetTableValues
		router.Route[fmt.Sprintf("GET/%s/id", table)] = h.ShowTuple
		router.Route[fmt.Sprintf("PUT/%s/", table)] = h.PutTuple
		router.Route[fmt.Sprintf("PUT/%s", table)] = h.PutTuple
		router.Route[fmt.Sprintf("POST/%s/id", table)] = h.UpdateTuple
		router.Route[fmt.Sprintf("DELETE/%s/id", table)] = h.DeleteTuple
	}

	// GET: Отдаем одну запись по ID
}
