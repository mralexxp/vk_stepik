package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

type Explorer struct {
	DB *sql.DB
}

func NewDbExplorer(db *sql.DB) (http.Handler, error) {
	e := &Explorer{DB: db}
	return e, nil
}

func (e *Explorer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.URL.Path + "\n"))
	w.Write([]byte(r.URL.Query().Get("admin") + "\n"))

}
