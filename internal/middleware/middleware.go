package middleware

import (
	"avenger/pkg/utils"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("[%s] %s %s from %s", time.Now().Format("2006/01/02 15:04:05"), r.Method, r.URL.Path, r.RemoteAddr)

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		slog.Info("HTTP request completed", slog.String("method", r.Method), slog.String("path", r.URL.Path), slog.Duration("duration", duration))
	})
}

func AuthMiddleware(next http.HandlerFunc, allowedRoles ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			slog.Warn("Missing authorization header", slog.String("path", r.URL.Path))
			writeAuthError(w, http.StatusUnauthorized, "Missing authorization token")
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			slog.Warn("Invalid Authorization format", slog.String("path", r.URL.Path))
			writeAuthError(w, http.StatusUnauthorized, "Invalid authorization format. Use: Bearer <token>")
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			slog.Warn("Empty token provided", slog.String("path", r.URL.Path))
			writeAuthError(w, http.StatusUnauthorized, "Empty authorization token")
			return
		}

		claims, err := utils.ValidateToken(token)
		if err != nil {
			slog.Warn("Invalid token", slog.String("path", r.URL.Path), slog.Any("error", err))
			writeAuthError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		if len(allowedRoles) > 0 {
			roleAllowed := false
			for _, role := range allowedRoles {
				if claims.Role == role {
					roleAllowed = true
					break
				}
			}

			if !roleAllowed {
				slog.Warn("Access forbidden", slog.String("path", r.URL.Path), slog.String("user_role", claims.Role), slog.Any("allowed_roles", allowedRoles))
				writeAuthError(w, http.StatusForbidden, "You dont have permission to access this resource")
				return
			}

			slog.Debug("Authentication successful", slog.Int("user_id", claims.UserID), slog.String("role", claims.Role), slog.String("path", r.URL.Path))

			next(w, r)
		}
	}
}

func writeAuthError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{
		"message": message,
		"errors":  nil,
	})
}
