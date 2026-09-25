package server

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
)

func (s *Server) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			authValue := r.Header.Get("Authorization")
			if !strings.HasPrefix(authValue, "Bearer ") {
				slog.Warn("invalid auth header format")
				s.respondWithError(w, http.StatusUnauthorized, "invalid request")
				return
			}
			token := strings.TrimPrefix(authValue, "Bearer ")
			claims, err := s.JWTManager.Verify(token)
			if err != nil {
				slog.Warn("user failed to verify bearer token", "error", err)
				s.respondWithError(w, http.StatusUnauthorized, "invalid request token/session expired")
				return
			}
			ctx := context.WithValue(r.Context(), contextKeyPlayerID, claims.Subject)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		},
	)
}
