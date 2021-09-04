package middleware

import (
	"encoding/json"
	"net/http"
	"xrate/services/api/auth"
)

const (
	exceptPath     = "/project"
	accessKeyQuery = "access_key"
)

type MAuth struct {
	auth auth.Auth
}

func NewMAuth(auth auth.Auth) *MAuth {
	return &MAuth{
		auth: auth,
	}
}

func (m *MAuth) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == exceptPath {
			next.ServeHTTP(w, r)
			return
		}

		accessKey := r.URL.Query().Get(accessKeyQuery)
		if !m.auth.Validate(accessKey) {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
