package router

import "net/http"

// TODO: полноценный MW с m.Use(Mw func)
func TempMiddleWare(w *http.ResponseWriter) {
	(*w).Header().Add("content-type", "application/JSON; charset=utf-8")
}
