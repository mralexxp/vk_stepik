package main

import (
	"encoding/json"
	"net/http"
)

type CustomError interface {
	error
	Code() int
}

type Error struct {
	msg  string
	code int
}

func (e *Error) Error() string {
	return e.msg
}

func (e *Error) Code() int {
	return e.code
}

func NewErr(err error, code int) CustomError {
	if err == nil {
		return nil
	}
	return &Error{err.Error(), code}
}

func responseError(w http.ResponseWriter, e CustomError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code())
	err := json.NewEncoder(w).Encode(SearchErrorResponse{Error: e.Error()})
	if err != nil {
		panic(err)
	}
}
