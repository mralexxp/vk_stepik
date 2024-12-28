package router

import (
	"db_explorer/internal/errors"
	"encoding/json"
	"github.com/go-sql-driver/mysql"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Handlers struct {
	Router
}

// TODO: ТАБЛИЦЫ ДОЛЖНЫ ПОЛУЧАТЬ НЕПОСРЕДСТВЕННО при инициализации
func (h *Handlers) Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errors.SendJSONError(w, http.StatusMethodNotAllowed, "method not allowed: "+r.Method)
		return
	}

	tables, err := h.Explorer.ShowTables()
	if err != nil {
		errors.SendJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = SendResponse(w, http.StatusOK, map[string]interface{}{
		"tables": tables,
	})
	if err != nil {
		panic(err)
	}
}

// Получаем содержимое таблиц. Параметры limit и offset по умолчанию 5 и 0 при отсутствии.
func (h *Handlers) GetTableValues(w http.ResponseWriter, r *http.Request) {
	table := strings.Split(r.URL.Path, "/")[1]

	response, err := h.Explorer.ShowTable(table, GetParams(r))
	if err != nil {
		errors.SendJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = SendResponse(w, http.StatusOK, map[string]interface{}{
		"records": response,
	})
	if err != nil {
		panic(err)
	}
}

// GET: Получение записи по ID
func (h *Handlers) ShowTuple(w http.ResponseWriter, r *http.Request) {
	URL := strings.Split(r.URL.Path, "/")
	if len(URL) != 3 {
		errors.SendJSONError(w, http.StatusBadRequest, "bad request URL: "+r.URL.Path)
		return
	}

	response := map[string]interface{}{}

	var id int
	var err error
	if id, err = strconv.Atoi(URL[2]); err != nil {
		errors.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err = h.Explorer.GetTuple(URL[1], id)
	if err != nil {
		errors.SendJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Если не нашли
	if len(response) == 0 {
		errors.SendJSONError(w, http.StatusNotFound, "record not found")
		return
	}

	err = SendResponse(w, http.StatusOK, map[string]interface{}{
		"record": response,
	})
	if err != nil {
		panic(err)
	}
}

func (h *Handlers) PutTuple(w http.ResponseWriter, r *http.Request) {
	table := strings.Split(r.URL.Path, "/")[1]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errors.SendJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer r.Body.Close()

	data := map[string]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		errors.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	added, err := h.Explorer.PutTuple(table, data)
	if err != nil {
		// TODO: Статус-коды mySQL можно имплементировать в определенные ошибки, чтобы не отдавать юзеру встроенную ошибку
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			errors.SendJSONError(w, http.StatusInternalServerError, mysqlErr.Error())
			return

		} else {
			errors.SendJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	err = SendResponse(w, http.StatusOK, added)
	if err != nil {
		panic(err)
	}

}

func (h *Handlers) UpdateTuple(w http.ResponseWriter, r *http.Request) {
	URL := strings.Split(r.URL.Path, "/")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		errors.SendJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(body, &data)
	if err != nil {
		errors.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	table := URL[1]
	id, err := strconv.Atoi(URL[2])
	if err != nil {
		errors.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	updated, err := h.Explorer.UpdateTuple(table, id, data)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			errors.SendJSONError(w, http.StatusInternalServerError, mysqlErr.Error())
			return
		} else {
			errors.SendJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	err = SendResponse(w, http.StatusOK, map[string]interface{}{
		"updated": updated,
	})
	if err != nil {
		panic(err)
	}
}

func (h *Handlers) DeleteTuple(w http.ResponseWriter, r *http.Request) {
	URL := strings.Split(r.URL.Path, "/")
	table := URL[1]
	id, err := strconv.Atoi(URL[2])
	if err != nil {
		errors.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	deleted, err := h.Explorer.DeleteTuple(table, id)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			errors.SendJSONError(w, http.StatusInternalServerError, mysqlErr.Error())
			return
		}
		errors.SendJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = SendResponse(w, http.StatusOK, map[string]interface{}{
		"deleted": deleted,
	})
	if err != nil {
		panic(err)
	}
}

func GetParams(r *http.Request) map[string]int {
	params := make(map[string]int)

	var limit int
	var offset int
	var err error

	if limit, err = strconv.Atoi(r.URL.Query().Get("limit")); err == nil {
		params["limit"] = limit
	} else {
		params["limit"] = 5
	}

	if offset, err = strconv.Atoi(r.URL.Query().Get("offset")); err == nil {
		params["offset"] = offset
	} else {
		params["offset"] = 0
	}

	return params
}
