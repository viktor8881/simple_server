package middleware

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func AuthMiddleware(logger *zap.Logger, apiKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Warn("Missing Authorization header")
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				logger.Warn("Invalid Authorization header format")
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			if token != apiKey {
				logger.Warn("Invalid token", zap.String("token", token))
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Создание нового контекста с информацией о пользователе
			ctx := context.WithValue(r.Context(), "user", "extracted_user_info") // Замените на реальную информацию

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
