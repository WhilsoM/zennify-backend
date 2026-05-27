package http

import (
	"net/http"
)

func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {
	userID, _ := userClaimsFromContext(r.Context())
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	profile, err := h.svc.GetUserProfile(r.Context(), userID)
	if err != nil {
		writeErrorJSON(w, err)
		return
	}

	writeJSON(w, http.StatusOK, profile)
}
