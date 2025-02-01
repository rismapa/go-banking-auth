package config

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	logger "github.com/rismapa/go-banking-lib/config"
)

type Claims struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Secret key for JWT
var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// Generate JWT Token
func GenerateJWT(id, username string) (string, string, error) {
	logger.GetLog().Info().Msg("Initializing Generate JWT Token")
	expirationTime := time.Now().Add(2 * time.Hour)

	claims := &Claims{
		ID:       id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	converedtTime := claims.RegisteredClaims.ExpiresAt.Format("2006-01-02 15:04:05")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		logger.GetLog().Error().Err(err).Msg("Failed to generate token")
		return "", "", err
	}

	logger.GetLog().Info().
		Str("username", username).
		Str("expiresAt", converedtTime).
		Msg("Token generated successfully")

	return tokenString, converedtTime, nil
}

// ParseToken for validating JWT
func ParseToken(tokenString string) (*Claims, error) {
	logger.GetLog().Info().Msg("Initializing Parse JWT Token")
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		logger.GetLog().Error().Err(err).Msg("Failed to parse token")
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		logger.GetLog().Error().Msg("Invalid token")
		return nil, fmt.Errorf("invalid token")
	}

	logger.GetLog().Info().
		Str("username", claims.Username).
		Str("expiresAt", claims.ExpiresAt.String()).
		Msg("Token parsed successfully")

	return claims, nil
}
