package server

import (
	"encoding/json"
	"net/http"
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
