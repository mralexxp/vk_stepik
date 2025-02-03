package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"rwa/internal/dto"
	"rwa/internal/errs"
	"rwa/internal/utils"
)

func (h *Handlers) UserLogin(w http.ResponseWriter, r *http.Request) {
	const op = "Handlers.UsersLogin"

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errs.SendError(w, http.StatusUnprocessableEntity, "Error reading request body")
		return
	}

	requestDTO := &dto.UserRequest{}

	err = json.Unmarshal(body, requestDTO)
	if err != nil {
		errs.SendError(w, http.StatusUnprocessableEntity, "Error parsing request body")
		return
	}

	responseDTO, err := h.Svc.LoginUser(requestDTO)
	if err != nil {
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

	requestDTO := &dto.UserRequest{}

	err = json.Unmarshal(body, requestDTO)
	if err != nil {
		log.Println(op, err)
		return
	}

	responseDTO, err := h.Svc.RegisterUser(requestDTO)
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
	const op = "Handlers.UserGet"

	token, err := utils.GetHeaderToken(r)
	if err != nil {
		errs.SendError(w, http.StatusUnauthorized, err.Error())
	}

	responseDTO, err := h.Svc.GetCurrentUser(token)
	if err != nil {
		errs.SendError(w, http.StatusUnprocessableEntity, err.Error())
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(responseDTO)
	if err != nil {
		log.Println(op + ": " + err.Error())
	}
}

func (h *Handlers) UserUpdate(w http.ResponseWriter, r *http.Request) {
	const op = "Handlers.UserUpdate"

	// Для получения ID пользователя в бизнес-логике, в контексте передавать такую информацию опасно?
	token, err := utils.GetHeaderToken(r)
	if err != nil {
		errs.SendError(w, http.StatusUnauthorized, err.Error())
		return
	}

	requestDTO := &dto.UserRequest{}
	err = json.NewDecoder(r.Body).Decode(requestDTO)
	if err != nil {
		errs.SendError(w, http.StatusUnprocessableEntity, "Error parsing request body")
		return
	}

	requestDTO.User.Token = token
	responseDTO, err := h.Svc.UpdateUser(requestDTO)
	if err != nil {
		errs.SendError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(responseDTO)
	if err != nil {
		log.Println(op + ": " + err.Error())
	}
}
