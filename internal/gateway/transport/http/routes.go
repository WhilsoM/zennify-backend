package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/zennify/backend/internal/gateway/app"
)

func NewRouter(svc *app.Service, logger *zap.Logger) http.Handler {
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

// TODO: fix this route
func mountProtectedRoutes(r chi.Router, h *Handler) {
	r.Get("/me", h.health)
}
