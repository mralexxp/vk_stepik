package handlers

import (
	"context"
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
		if method, ok := h.NoAuth[r.URL.Path]; ok && method == r.Method {
			next.ServeHTTP(w, r)
			return
		}

		token, err := utils.GetHeaderToken(r)
		if err != nil {
			errs.SendError(w, http.StatusUnauthorized, err.Error())
			return
		}

		id, ok := h.Svc.GetSessionManager().Check(token)
		if !ok {
			errs.SendError(w, http.StatusUnauthorized, err.Error())
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "id", id)))
	})
}
