package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JSON отвечет на запрос в json формате
func JSON(rw http.ResponseWriter, data interface{}, code int) {
	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)

	// nolint:errcheck
	json.NewEncoder(rw).Encode(data)
}

// JSONError отвечет json ошибкой
func JSONError(rw http.ResponseWriter, err string, code int) {
	rw.Header().Set("Content-Type", "application/json")
	// ХЗ зачем это, скопировать из http.Error
	rw.Header().Set("X-Content-Type-Options", "nosniff")
	rw.WriteHeader(code)
	//nolint
	fmt.Fprintln(rw, fmt.Sprintf(`{"error": "%s"}`, err))
}
