package handlers

import "net/http"

func (h *Handlers) MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte(`{"message":"Method not allowed"}`)) // TODO: подогнать под стандарт ошибок
}
