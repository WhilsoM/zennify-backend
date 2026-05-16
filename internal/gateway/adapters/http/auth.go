package httpapi

import (
	"encoding/json"
	"net/http"

	"google.golang.org/grpc/codes"

	"github.com/zennify/backend/internal/gateway/core/ports"
	"github.com/zennify/backend/internal/shared/grpcerr"
)

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var req ports.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorJSON(w, grpcerr.ClientError(codes.InvalidArgument, grpcerr.MsgInvalidRequest))
		return
	}
	if err := h.vld.Struct(req); err != nil {
		h.writeErrorJSON(w, grpcerr.ClientError(codes.InvalidArgument, grpcerr.MsgInvalidRequest))
		return
	}

	resp, err := h.svc.Register(r.Context(), &req)
	if err != nil {
		h.writeErrorJSON(w, err)
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
		h.writeErrorJSON(w, grpcerr.ClientError(codes.InvalidArgument, grpcerr.MsgInvalidRequest))
		return
	}

	resp, err := h.svc.Login(r.Context(), &req)
	if err != nil {
		h.writeErrorJSON(w, err)
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
		h.writeErrorJSON(w, grpcerr.ClientError(codes.InvalidArgument, grpcerr.MsgInvalidRequest))
		return
	}

	resp, err := h.svc.RefreshTokens(r.Context(), &req)
	if err != nil {
		h.writeErrorJSON(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
