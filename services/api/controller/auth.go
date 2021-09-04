package controller

import (
	"encoding/json"
	"net/http"
	"strings"
	"xrate/services/api/auth"
)

type AuthHandler struct {
	auth auth.Auth
}

func NewAuthHandler(auth auth.Auth) *AuthHandler {
	return &AuthHandler{
		auth: auth,
	}
}

func (h *AuthHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectName := r.URL.Query().Get("name")
	if len(strings.TrimSpace(projectName)) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "project name can't be empty"})
		return
	}

	token, err := h.auth.Create(projectName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"access_key": token})
}
