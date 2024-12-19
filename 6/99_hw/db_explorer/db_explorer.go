package main

import (
	"database/sql"
	"db_explorer/internal/explorer"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func NewDbExplorer(db *sql.DB) (http.Handler, error) {
	e := &explorer.Explorer{DB: db}
	e.InitDBStruct()

	return e, nil
}
