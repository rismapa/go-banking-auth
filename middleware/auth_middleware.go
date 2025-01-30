package middleware

import (
	"context"
	"net/http"
	"strings"

	config "github.com/okyws/go-banking-auth/config"
	"github.com/okyws/go-banking-auth/service"
	"github.com/okyws/go-banking-auth/utils"
)

// AuthMiddleware untuk validasi JWT
func AuthMiddleware(authService service.AuthService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			utils.ErrorResponse(w, http.StatusUnauthorized, "error", "Unauthorized, missing Authorization header")
			return
		}

		// remove "Bearer " prefix
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims, err := config.ParseToken(tokenString)
		if err != nil {
			utils.ErrorResponse(w, http.StatusUnauthorized, "error", err.Error())
			return
		}

		isValid, err := authService.ValidateToken(tokenString)
		if err != nil || !isValid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// get context from request then set data
		ctx := r.Context()
		ctx = context.WithValue(ctx, "id", claims.ID)
		ctx = context.WithValue(ctx, "username", claims.Username)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
