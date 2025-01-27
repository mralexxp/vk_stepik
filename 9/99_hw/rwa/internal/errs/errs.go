package errs

import (
	"encoding/json"
	"net/http"
)

type Err struct {
	Errors struct {
		Body []string
	}
}

func SendError(w http.ResponseWriter, code int, msg string) {
	err := NewErr(msg)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}

func NewErr(msg string) *Err {
	err := Err{
		Errors: struct {
			Body []string
		}{
			Body: []string{msg},
		},
	}

	return &err
}
