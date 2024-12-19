package router

import "net/http"

type Handlers struct {
	Router
}

// ТАБЛИЦЫ ДОЛЖНЫ ПОЛУЧАТЬ НЕПОСРЕДСТВЕННО при инициализации
func (h *Handlers) Tables(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {

	}
}
