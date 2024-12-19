package main

import (
	"database/sql"
	"db_explorer/internal/explorer"
	"db_explorer/internal/router"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func NewDbExplorer(db *sql.DB) (http.Handler, error) {
	e := &explorer.Explorer{
		DB: db,
	}

	r := router.NewRouter(e)

	e.InitDBStruct()

	return r, nil
}
