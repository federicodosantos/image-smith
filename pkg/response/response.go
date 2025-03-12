package pkg

import (
	"encoding/json"
	"net/http"
)

type HttpResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"obj,omitempty"`
}

func SuccessResponse(w http.ResponseWriter, status int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := HttpResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

func FailedResponse(w http.ResponseWriter, status int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errorResponse := HttpResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(errorResponse)
}
