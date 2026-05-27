package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zennify/backend/internal/gateway/core/ports"
)

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var req ports.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJSON(w, status.Error(codes.InvalidArgument, ""))
		return
	}
	if err := h.vld.Struct(req); err != nil {
		writeErrorJSON(w, status.Error(codes.InvalidArgument, ""))
		return
	}

	resp, err := h.svc.Register(r.Context(), &req)
	if err != nil {
		writeErrorJSON(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req ports.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJSON(w, status.Error(codes.InvalidArgument, ""))
		return
	}
	if err := h.vld.Struct(req); err != nil {
		writeErrorJSON(w, status.Error(codes.InvalidArgument, ""))
		return
	}

	resp, err := h.svc.Login(r.Context(), &req)
	if err != nil {
		writeErrorJSON(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	var req ports.RefreshTokensRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJSON(w, status.Error(codes.InvalidArgument, ""))
		return
	}
	fmt.Println("req", req.RefreshToken)
	if err := h.vld.Struct(req); err != nil {
		writeErrorJSON(w, status.Error(codes.InvalidArgument, ""))
		return
	}

	resp, err := h.svc.RefreshTokens(r.Context(), &req)
	if err != nil {
		writeErrorJSON(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
