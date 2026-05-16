package httpapi

import (
	"encoding/json"
	"net/http"

	"google.golang.org/grpc/codes"

	"github.com/go-playground/validator/v10"

	"github.com/zennify/backend/internal/gateway/app"
	"github.com/zennify/backend/internal/gateway/ports"
	"github.com/zennify/backend/internal/shared/grpcerr"
)

type Handler struct {
	svc *app.Service
	vld *validator.Validate
}

func newHandler(svc *app.Service) *Handler {
	return &Handler{
		svc: svc,
		vld: validator.New(),
	}
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var req ports.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	if err := h.vld.Struct(req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": grpcerr.MsgInvalidRequest})
		return
	}

	resp, err := h.svc.Register(r.Context(), &req)
	if err != nil {
		h.writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req ports.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if err := h.vld.Struct(req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": grpcerr.MsgInvalidRequest})
		return
	}
	resp, err := h.svc.Login(r.Context(), &req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	var req ports.RefreshTokensRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	if err := h.vld.Struct(req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": grpcerr.MsgInvalidRequest})
		return
	}

	resp, err := h.svc.RefreshTokens(r.Context(), &req)
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	st := grpcerr.Convert(err)
	switch st.Code() {
	case codes.InvalidArgument:
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": st.Message()})
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
