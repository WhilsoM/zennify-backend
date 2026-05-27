package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/zennify/backend/internal/gateway/core/services"
)

func NewRouter(svc *services.Service, logger *zap.Logger) http.Handler {
	h := newHandler(svc)
	r := chi.NewRouter()

	r.Use(requestIDMiddleware)
	r.Use(recoverMiddleware(logger))
	r.Use(loggingMiddleware(logger))

	r.Get("/health", h.health)

	r.Route("/auth", mountAuthRoutes(h))

	r.Group(func(gr chi.Router) {
		gr.Use(authMiddleware(svc))
		mountProtectedRoutes(gr, h)

		gr.Get("/ws", h.connectWs)
	})

	return r
}

func mountAuthRoutes(h *Handler) func(chi.Router) {
	return func(r chi.Router) {
		r.Post("/register", h.register)
		r.Post("/login", h.login)
		r.Post("/refresh", h.refresh)
	}
}

func mountProtectedRoutes(r chi.Router, h *Handler) {
	r.Get("/users/me", h.getProfile)
}
