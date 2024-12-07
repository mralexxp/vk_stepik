package main

import (
	"encoding/json"
	"net/http"
)

const (
	fileName = "dataset.xml"
)

func main() {
	http.HandleFunc("/", SearchServer)

	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		panic(err)
	}
}

// Handler
func SearchServer(w http.ResponseWriter, r *http.Request) {
	sr, err := parseRequest(r)
	if err != nil {
		responseError(w, err)
		return
	}

	users, err := Search(sr)
	if err != nil {
		responseError(w, NewErr(err, http.StatusInternalServerError))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = NewErr(json.NewEncoder(w).Encode(*users), http.StatusInternalServerError)
	if err != nil {
		// TODO: заменить панику: не падать во время обрыва
		panic("write response error: " + err.Error())
	}
}
