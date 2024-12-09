package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	useXML = fullXml

	fullXml = "dataset.xml"
)

func main() {
	router := mux.NewRouter()
	router.Handle("/", http.HandlerFunc(SearchServer))

	err := http.ListenAndServe("127.0.0.1:8080", router)
	if err != nil {
		panic(err)
	}
}

// SearchServer handler
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
		panic("write response error: " + err.Error())
	}
}
