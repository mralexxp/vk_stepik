package router

import (
	"encoding/json"
	"net/http"
)

func SendResponse(w http.ResponseWriter, code int, data map[string]interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	response := map[string]map[string]interface{}{
		"response": data,
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return err
	}

	return nil
}
