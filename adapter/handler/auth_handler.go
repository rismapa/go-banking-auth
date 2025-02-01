package adapter

import (
	"encoding/json"
	"net/http"

	"github.com/rismapa/go-banking-auth/dto"
	"github.com/rismapa/go-banking-auth/service"
	"github.com/rismapa/go-banking-auth/utils"
	logger "github.com/rismapa/go-banking-lib/config"

	"github.com/go-playground/validator/v10"
)

type AuthHandlerDB struct {
	Service   service.AuthService
	Validator validator.Validate
}

func NewAuthHandlerDB(service service.AuthService) *AuthHandlerDB {
	return &AuthHandlerDB{Service: service, Validator: *validator.New()}
}

func (h *AuthHandlerDB) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, "error", "Method not allowed")
		return
	}

	logger.GetLog().Info().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Msg("Initializing Login")

	var req dto.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "error", "Invalid request body")
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		errorMessage := utils.CustomValidationError(err)
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, "error", errorMessage)
		return
	}

	token, expiresAt, err := h.Service.LoginAccount(req.Username, req.Password)
	if err != nil {
		logger.GetLog().Error().Err(err).Msg("Username or password is incorrect. Failed to login")
		utils.ErrorResponse(w, http.StatusUnauthorized, "error", err.Error())
		return
	}

	resp := dto.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}

	utils.ResponseJSON(w, resp, http.StatusOK, "success", "Login successful")
	logger.GetLog().Info().Str("username", req.Username).Msg("Login successful")
}
