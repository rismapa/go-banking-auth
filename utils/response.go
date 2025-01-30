package utils

import (
	"encoding/json"
	"net/http"

	"github.com/okyws/go-banking-auth/dto"
	logger "github.com/okyws/go-banking-lib/config"
)

func ResponseJSON(w http.ResponseWriter, data interface{}, code int, status string, message string) {
	resp := dto.SuccessResponseDTO[interface{}]{
		Status:  status,
		Code:    code,
		Message: message,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error encoding response" + err.Error()))

		logger.GetLog().Error().
			Err(err).
			Str("error", err.Error()).
			Msg("Failed to encode response")
		return
	}

	logger.GetLog().Info().
		Str("status", status).
		Int("code", code).
		Str("message", message).
		Msg("Response Success sent successfully")
}

func ErrorResponse(w http.ResponseWriter, code int, status string, message string) {
	resp := dto.ErrorResponseDTO{
		Status:  status,
		Code:    code,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(code)
		w.Write([]byte("Error encoding response" + err.Error()))

		logger.GetLog().Error().
			Err(err).
			Str("error", err.Error()).
			Msg("Failed to encode response")
		return
	}

	logger.GetLog().Info().
		Str("status", status).
		Int("code", code).
		Str("message", message).
		Msg("Response Error sent successfully")
}
