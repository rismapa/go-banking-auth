package service

import (
	"fmt"
	"time"

	adapter "github.com/okyws/go-banking-auth/adapter/repository"
	config "github.com/okyws/go-banking-auth/config"
	logger "github.com/okyws/go-banking-lib/config"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	LoginAccount(username, password string) (string, string, error)
	ValidateToken(token string) (bool, error)
}

type AuthAdapterDB struct {
	repo adapter.AuthRepository
}

func NewAuthService(repo adapter.AuthRepository) *AuthAdapterDB {
	return &AuthAdapterDB{repo: repo}
}

func (u *AuthAdapterDB) LoginAccount(username, password string) (string, string, error) {
	logger.GetLog().Info().Str("username", username).Msg("LoginAccount started")
	user, err := u.repo.GetAccountByUsername(username)
	if err != nil {
		return "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logger.GetLog().Error().Err(err).Msg("Username or password is incorrect. Failed to login")
		return "", "", fmt.Errorf("invalid password: %v", err)
	}

	token, expiresAt, err := config.GenerateJWT(user.ID, user.Username)
	if err != nil {
		logger.GetLog().Error().Err(err).Msg("Failed to generate token")
		return "", "", fmt.Errorf("could not generate token: %v", err)
	}

	err = u.repo.SaveToken(user.ID, token, expiresAt)
	if err != nil {
		logger.GetLog().Error().Err(err).Msg("Failed to save token")
		return "", "", fmt.Errorf("could not save token: %v", err)
	}

	logger.GetLog().Info().
		Str("username", user.Username).
		Str("expiresAt", expiresAt).
		Msg("Login success")

	return token, expiresAt, nil
}

func (s *AuthAdapterDB) ValidateToken(token string) (bool, error) {
	// Get token expiration from the repository
	expiresAt, err := s.repo.GetTokenExpiration(token)
	if err != nil {
		logger.GetLog().Error().Err(err).Msg("Failed to validate token")
		return false, fmt.Errorf("failed to validate token: %v", err)
	}

	// Check if the token is expired
	if expiresAt.Before(time.Now()) {
		logger.GetLog().Info().Str("token", token).Msg("Token is expired")
		return false, nil
	}

	logger.GetLog().Info().Msg("Token is valid")
	return true, nil
}
