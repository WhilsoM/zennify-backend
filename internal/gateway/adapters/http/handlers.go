package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/zennify/backend/internal/gateway/core/services"
)

type Handler struct {
	svc *services.Service
	vld *validator.Validate
}

func newHandler(svc *services.Service) *Handler {
	return &Handler{
		svc: svc,
		vld: validator.New(),
	}
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
