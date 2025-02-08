package handlers

import (
	"log"
	"net/http"
)

func (h *Handlers) MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	const op = "Handlers.MethodNotAllowedHandler"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	_, err := w.Write([]byte(`{"message":"Method not allowed"}`))
	if err != nil {
		log.Println(op + ": " + err.Error())
	}
}
