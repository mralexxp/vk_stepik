package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"rwa/internal/dto"
	"rwa/internal/errs"
)

func (h *Handlers) UserLogin(w http.ResponseWriter, r *http.Request) {
	const op = "Handlers.UsersLogin"

	panic(op + ": not implemented")
}

func (h *Handlers) UserRegister(w http.ResponseWriter, r *http.Request) {
	const op = "Handlers.UsersRegister"

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(op, err)
		errs.SendError(w, http.StatusBadRequest, "body read error: "+err.Error())
		return
	}

	RequestDTO := &dto.UserRegisterRequest{}

	err = json.Unmarshal(body, RequestDTO)
	if err != nil {
		log.Println(op, err)
		return
	}

	ResponseDTO, err := h.Svc.Add(RequestDTO)
	if err != nil {
		log.Println(op + ": " + err.Error())
		errs.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(ResponseDTO)
	if err != nil {
		log.Println(op + ": " + err.Error())
		return
	}
}

func (h *Handlers) UserGet(w http.ResponseWriter, r *http.Request) {
	const op = "Handlers.UsersGet"
	panic(op + ": not implemented")
}

func (h *Handlers) UserUpdate(w http.ResponseWriter, r *http.Request) {
	const op = "Handlers.UsersUpdate"
	panic(op + ": not implemented")
}
