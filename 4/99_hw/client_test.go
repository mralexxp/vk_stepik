package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var Server *httptest.Server

func TestMain(m *testing.M) {
	Server = httptest.NewServer(ishandler)
}

func TestSearchServer(t *testing.T) {

}

func ishandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}
