package response

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, statusCode int, data any) error {

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, statusCode int, message string) {

	data := map[string]string{
		"message": message,
	}

	_ = JSON(w, statusCode, data)
}
