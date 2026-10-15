package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"uuid"
)

type jsonResponse struct {
	Status string `json:"status"`
	Data   any    `json:"data,omitempty"`
	Error  string `json:"error,omitempty"`
}

func (s *Server) respondWithJSON(w http.ResponseWriter, code int, payload jsonResponse) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(&payload)
}

func (srv *Server) respondWithError(w http.ResponseWriter, code int, message string) error {
	response := jsonResponse{Status: "error", Error: message}
	return srv.respondWithJSON(w, code, response)
}

func (s *Server) createGuestToken(w http.ResponseWriter) (string, error) {
	playerID := fmt.Sprintf("player_%s", uuid.New().String())
	token, err := s.JWTManager.Sign(playerID)
	if err != nil {
		slog.Error("failed to create guest token", "error", err)
		s.respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return "", err
	}
	return token, nil
}
