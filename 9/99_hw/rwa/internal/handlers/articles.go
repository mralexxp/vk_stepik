package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"rwa/internal/dto"
	"rwa/internal/errs"
	"rwa/internal/utils"
)

func (h *Handlers) GetArticlesByFilter(w http.ResponseWriter, r *http.Request) {
	const op = "Handlers.GetAllArticles"

	query := r.URL.Query()

	response, err := h.Svc.ArticlesByFilter(&query)
	if err != nil {
		errs.SendError(w, http.StatusUnprocessableEntity, "unknown error: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("%s: %s", op, err.Error())
		return
	}

}

func (h *Handlers) CreateArticle(w http.ResponseWriter, r *http.Request) {
	const op = "Handlers.CreateArticle"

	token, err := utils.GetHeaderToken(r)
	if err != nil {
		errs.SendError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	aReq := &dto.ArticleRequest{}
	err = json.NewDecoder(r.Body).Decode(aReq)
	if err != nil {
		errs.SendError(w, http.StatusUnprocessableEntity, "unknown error: "+err.Error())
		return
	}
	defer r.Body.Close()

	aRes, err := h.Svc.CreateArticle(aReq, token)
	if err != nil {
		errs.SendError(w, http.StatusUnprocessableEntity, "unknown error: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(aRes)
	if err != nil {
		log.Printf("%s: %s", op, err.Error())
		return
	}
}

func (h *Handlers) GetArticle(w http.ResponseWriter, r *http.Request) {
	const op = "Handlers.GetArticle"

	panic(op + ": not yet implemented")
}

func (h *Handlers) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	const op = "Handlers.UpdateArticle"

	panic(op + ": not yet implemented")
}

func (h *Handlers) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	const op = "Handlers.DeleteArticle"

	panic(op + ": not yet implemented")
}
