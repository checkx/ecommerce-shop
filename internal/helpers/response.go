package helpers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type SuccessResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func WriteSuccess(w http.ResponseWriter, message string, data interface{}) {
	resp := SuccessResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func WriteError(w http.ResponseWriter, code int, message, errDetail string, logger *zap.Logger) {
	resp := ErrorResponse{
		Status:  "error",
		Message: message,
		Error:   errDetail,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
	if logger != nil {
		logger.Error(message, zap.String("error", errDetail))
	}
}
