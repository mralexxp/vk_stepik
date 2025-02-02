package handlers

import (
	"net/http"
	"rwa/internal/errs"
	"rwa/internal/utils"
)

func (h *Handlers) ContentTypeMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		next.ServeHTTP(w, r)
	})
}

func (h *Handlers) AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := h.NoAuth[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
		}

		token, err := utils.GetHeaderToken(r)
		if err != nil {
			errs.SendError(w, http.StatusUnauthorized, err.Error())
		}

	})
}
