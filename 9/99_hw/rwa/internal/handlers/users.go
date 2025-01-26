package handlers

import (
	"encoding/json"
	"net/http"
	"rwa/internal/dto"
)

func (h *Handlers) UserLogin(w http.ResponseWriter, r *http.Request) {
	const op = "Handlers.UsersLogin"

	panic(op + ": not implemented")
}

//func (h *Handlers) UsersLogout(w http.ResponseWriter, r *http.Request) {
//	const op = "Handlers.UsersLogout"
//	panic(op + ": not implemented")
//}

func (h *Handlers) UserRegister(w http.ResponseWriter, r *http.Request) {
	const op = "Handlers.UsersRegister"

	RequestDTO := &dto.UserRegisterRequest{}

	err := json.NewDecoder(r.Body).Decode(RequestDTO)
	if err != nil {
		panic(op + ": " + err.Error())
	}

	// TODO: Реализовать
	ResponseDTO, err := h.Svc.Add(RequestDTO)
	if err != nil {
		//errs.SendError(fmt.Sprintf("%s: bad request: %v", op, RequestDTO))
		// TODO: Возвращаем ошибку, если вернулась ошибка из бизнеса
	}

	err = json.NewEncoder(w).Encode(ResponseDTO)
	if err != nil {
		panic(op + ": " + err.Error())
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
