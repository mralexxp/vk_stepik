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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errs.SendError(w, http.StatusUnprocessableEntity, "Error reading request body")
		return
	}

	requestDTO := &dto.UserLoginRequest{}

	err = json.Unmarshal(body, requestDTO)
	if err != nil {
		errs.SendError(w, http.StatusUnprocessableEntity, "Error parsing request body")
		return
	}

	responseDTO, err := h.Svc.Login(requestDTO)
	if err != nil {
		// TODO: Возможно, ошибки придется отделить
		log.Println(op + ": " + err.Error())
		errs.SendError(w, http.StatusUnauthorized, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(responseDTO)
	if err != nil {
		log.Println(op + ": " + err.Error())
		return
	}
}

func (h *Handlers) UserRegister(w http.ResponseWriter, r *http.Request) {
	const op = "Handlers.UsersRegister"

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(op, err)
		errs.SendError(w, http.StatusUnprocessableEntity, "body read error: "+err.Error())
		return
	}

	requestDTO := &dto.UserRegisterRequest{}

	err = json.Unmarshal(body, requestDTO)
	if err != nil {
		log.Println(op, err)
		return
	}

	responseDTO, err := h.Svc.Register(requestDTO)
	if err != nil {
		log.Println(op + ": " + err.Error())
		errs.SendError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(responseDTO)
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
