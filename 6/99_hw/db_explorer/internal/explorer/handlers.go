package explorer

import (
	"net/http"
)

// Роутер
func (e *Explorer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.URL.Path + "\n"))
	w.Write([]byte(r.URL.Query().Get("admin") + "\n"))

}
