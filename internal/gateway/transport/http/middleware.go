package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/zennify/backend/internal/gateway/app"
)

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.NewString()
		}
		w.Header().Set("X-Request-ID", reqID)
		next.ServeHTTP(w, r.WithContext(withRequestID(r.Context(), reqID)))
	})
}

func recoverMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("panic recovered",
						zap.Any("panic", rec),
						zap.String("path", r.URL.Path),
						zap.String("request_id", requestIDFromContext(r.Context())),
					)
					writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func loggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			sw := &statusResponseWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(sw, r)

			userID, _ := userClaimsFromContext(r.Context())
			logger.Info("gateway request",
				zap.String("request_id", requestIDFromContext(r.Context())),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", sw.status),
				zap.Duration("duration", time.Since(start)),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_id", userID),
			)
		})
	}
}

func authMiddleware(svc *app.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
			if authHeader == "" || !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing bearer token"})
				return
			}
			token := strings.TrimSpace(authHeader[len("bearer "):])
			claims, err := svc.ValidateAccessToken(token)
			if err != nil {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token"})
				return
			}

			r.Header.Set("X-User-ID", claims.UserID)
			if claims.Username != "" {
				r.Header.Set("X-Username", claims.Username)
			}
			next.ServeHTTP(w, r.WithContext(withUserClaims(r.Context(), claims.UserID, claims.Username)))
		})
	}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
