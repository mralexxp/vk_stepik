package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"rwa/internal/dto"
	"rwa/internal/errs"
	"rwa/internal/utils"
)

func (h *Handlers) UserLogin(w http.ResponseWriter, r *http.Request) {
	const op = "Handlers.UsersLogin"

	requestDTO := &dto.UserRequest{}

	err := json.NewDecoder(r.Body).Decode(requestDTO)
	defer r.Body.Close()
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

	requestDTO := &dto.UserRequest{}

	err := json.NewDecoder(r.Body).Decode(requestDTO)
	defer r.Body.Close()
	if err != nil {
		errs.SendError(w, http.StatusUnprocessableEntity, "Error parsing request body")
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
		return
	}

	responseDTO, err := h.Svc.GetCurrentUser(token)
	if err != nil {
		errs.SendError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(responseDTO)
	if err != nil {
		log.Println(op + ": " + err.Error())
		return
	}
}

func (h *Handlers) UserUpdate(w http.ResponseWriter, r *http.Request) {
	const op = "Handlers.UserUpdate"

	// Через контекст от MW нельзя
	token, err := utils.GetHeaderToken(r)
	if err != nil {
		errs.SendError(w, http.StatusUnauthorized, err.Error())
		return
	}

	requestDTO := &dto.UserRequest{}
	err = json.NewDecoder(r.Body).Decode(requestDTO)
	defer r.Body.Close()
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
		return
	}
}

func (h *Handlers) UserLogout(w http.ResponseWriter, r *http.Request) {
	const op = "Handlers.UserLogout"

	token, err := utils.GetHeaderToken(r)
	if err != nil {
		errs.SendError(w, http.StatusUnauthorized, err.Error())
		return
	}

	_, err = h.Svc.LogoutUser(token)
	if err != nil {
		errs.SendError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
