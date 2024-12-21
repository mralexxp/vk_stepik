package router

import "net/http"

// TODO: Обернуть все ответы в application/JSON
func TempMiddleWare(w *http.ResponseWriter) {
	(*w).Header().Add("content-type", "application/JSON; charset=utf-8")
}
