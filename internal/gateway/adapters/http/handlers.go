package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"

	"github.com/zennify/backend/internal/gateway/core/services"
	"github.com/zennify/backend/internal/shared/grpcerr"
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

func (h *Handler) writeErrorJSON(w http.ResponseWriter, err error) {
	st := grpcerr.Convert(err)
	switch st.Code() {
	case codes.InvalidArgument:
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": grpcerr.MsgInvalidRequest})
	case codes.AlreadyExists:
		writeJSON(w, http.StatusConflict, map[string]string{"error": st.Message()})
	case codes.Unauthenticated:
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": st.Message()})
	case codes.NotFound:
		writeJSON(w, http.StatusNotFound, map[string]string{"error": st.Message()})
	case codes.DeadlineExceeded, codes.Unavailable:
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": grpcerr.MsgUpstreamUnavailable})
	case codes.Internal:
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": grpcerr.MsgInternal})
	default:
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": grpcerr.MsgInternal})
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
