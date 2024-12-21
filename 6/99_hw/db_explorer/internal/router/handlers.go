package router

import (
	"encoding/json"
	"net/http"
)

type Handlers struct {
	Router
}

// TODO: ТАБЛИЦЫ ДОЛЖНЫ ПОЛУЧАТЬ НЕПОСРЕДСТВЕННО при инициализации
func (h *Handlers) Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := map[string]string{
			"error": "Method Not Allowed",
		}
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			// TODO: !!! Не паниковать при закрытии соед.!!!
			panic(err)
		}
		return
	}

	tables, err := h.Explorer.ShowTables()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]string{
			"error": "Internal Server Error: " + err.Error(),
		}
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			// TODO: !!! Не паниковать при закрытии соед.!!!
			panic(err)
		}
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(tables)
	if err != nil {
		// TODO: !!! Не паниковать при закрытии соед.!!!
		panic(err)
	}
}

// Получаем содержимое таблиц. Параметры limit и offset по умолчанию 5 и 0 при отсутствии.
func (h *Handlers) GETTableValues(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GETTableValuesOK"))
}
