package errors

import (
	"encoding/json"
	"net/http"
)

// Отправка ошибок:
func SendJSONError(w http.ResponseWriter, code int, text string) {
	data := map[string]interface{}{
		"error": text,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		// TODO: Не паникуем
		panic(err)
	}
}
