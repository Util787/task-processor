package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type errResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(w http.ResponseWriter, log *slog.Logger, statusCode int, message string, err error) {
	log.Error(message, slog.String("error", err.Error()))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	respErr := json.NewEncoder(w).Encode(errResponse{Message: message})
	if respErr != nil {
		log.Error("Failed to write error response", slog.String("error", respErr.Error()))
	}
}
