package router

import (
	"encoding/json"
	"errors"
	"github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Handlers struct {
	Router
}

func (h *Handlers) Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		SendJSONError(w, http.StatusMethodNotAllowed, "method not allowed: "+r.Method)
		return
	}

	tables, err := h.Explorer.GetTables()
	if err != nil {
		SendJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = SendResponse(w, http.StatusOK, map[string]interface{}{
		"tables": tables,
	})
	if err != nil {
		log.Printf("error in sending response: %v", err)
	}
}

func (h *Handlers) GetTableTuples(w http.ResponseWriter, r *http.Request) {
	table := strings.Split(r.URL.Path, "/")[1]

	response, err := h.Explorer.ShowTable(table, GetParams(r))
	if err != nil {
		SendJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = SendResponse(w, http.StatusOK, map[string]interface{}{
		"records": response,
	})
	if err != nil {
		log.Printf("error in sending response: %v", err)
	}
}

func (h *Handlers) GetTuple(w http.ResponseWriter, r *http.Request) {
	URL := strings.Split(r.URL.Path, "/")
	if len(URL) != 3 {
		SendJSONError(w, http.StatusBadRequest, "bad request URL: "+r.URL.Path)
		return
	}

	response := map[string]interface{}{}

	id, err := strconv.Atoi(URL[2])
	if err != nil {
		SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err = h.Explorer.GetTuple(URL[1], id)
	if err != nil {
		SendJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(response) == 0 {
		SendJSONError(w, http.StatusNotFound, "record not found")
		return
	}

	err = SendResponse(w, http.StatusOK, map[string]interface{}{
		"record": response,
	})
	if err != nil {
		log.Printf("error in sending response: %v", err)
	}
}

func (h *Handlers) PutTuple(w http.ResponseWriter, r *http.Request) {
	table := strings.Split(r.URL.Path, "/")[1]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		SendJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer r.Body.Close()

	data := map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	added, err := h.Explorer.PutTuple(table, data)
	if err != nil {
		// Статус-коды mySQL можно имплементировать в определенные ошибки, чтобы не отдавать юзеру встроенную ошибку
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			SendJSONError(w, http.StatusInternalServerError, mysqlErr.Error())
			return

		} else {
			SendJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	err = SendResponse(w, http.StatusOK, added)
	if err != nil {
		log.Printf("error in sending response: %v", err)
	}
}

func (h *Handlers) UpdateTuple(w http.ResponseWriter, r *http.Request) {
	URL := strings.Split(r.URL.Path, "/")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		SendJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(body, &data)
	if err != nil {
		SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	table := URL[1]
	id, err := strconv.Atoi(URL[2])
	if err != nil {
		SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	updated, err := h.Explorer.UpdateTuple(table, id, data)
	if err != nil {
		var mysqlErr *mysql.MySQLError

		if errors.As(err, &mysqlErr) {
			SendJSONError(w, http.StatusInternalServerError, mysqlErr.Error())
			return
		}

		SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = SendResponse(w, http.StatusOK, map[string]interface{}{
		"updated": updated,
	})
	if err != nil {
		log.Printf("error in sending response: %v", err)
	}
}

func (h *Handlers) DeleteTuple(w http.ResponseWriter, r *http.Request) {
	URL := strings.Split(r.URL.Path, "/")
	table := URL[1]
	id, err := strconv.Atoi(URL[2])
	if err != nil {
		SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	deleted, err := h.Explorer.DeleteTuple(table, id)
	if err != nil {
		var mysqlErr *mysql.MySQLError

		if errors.As(err, &mysqlErr) {
			SendJSONError(w, http.StatusInternalServerError, mysqlErr.Error())
			return
		}

		SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = SendResponse(w, http.StatusOK, map[string]interface{}{
		"deleted": deleted,
	})
	if err != nil {
		log.Printf("error in sending response: %v", err)
	}
}
