package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"uuid"

	"github.com/jauharfz/battleship/internal/game"
	"github.com/jauharfz/battleship/internal/protocol"
)

func (s *Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := jsonResponse{
		Status: "health checked",
	}
	if err := s.respondWithJSON(w, http.StatusOK, response); err != nil {
		slog.Error("failed to check health", "error", err)
		s.respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}
}

// TODO: if client has guest_jwt before, create function to handle the guest_jwt/verify/checking if this guest_jwt is inthe match before
// and if expired && not in the match just auto create new jwt guest
func (s *Server) RegisterGuest(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
	playerID := fmt.Sprintf("player_%s", uuid.New().String())

	token, err := s.JWTManager.Sign(playerID)
	if err != nil {
		slog.Error("failed to create guest token", "error", err)
		s.respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	response := jsonResponse{
		Status: "token created succesfully",
		Data: map[string]string{
			"token": token,
		},
	}
	if err := s.respondWithJSON(w, http.StatusCreated, response); err != nil {
		slog.Error("failed to encode json response", "error", err)
		s.respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
}

func (s *Server) CheckShipPlacements(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	var data []protocol.ShipPlacementDTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&data); err != nil {
		slog.Warn("failed to decode json request", "error", err)
		s.respondWithError(w, http.StatusBadRequest, "bad request")
		return
	}
	// chekc ship validation

	shipPlacement := make([]game.ShipPlacement, len(data))
	for idx, ship := range data {
		var err error
		shipPlacement[idx], err = ship.ToDomain()
		if err != nil {
			slog.Warn("failed to casting DTO payload", "error", err)
			s.respondWithError(w, http.StatusBadRequest, "bad request")
			return
		}
	}
	if _, _, err := game.ValidateAllShipPosition(shipPlacement); err != nil {
		slog.Warn("failed to validate all ship position", "error", err)
		s.respondWithError(w, http.StatusBadRequest, "invalid ship position")
		return
	}

	response := jsonResponse{
		Status: "ship validated",
	}
	if err := s.respondWithJSON(w, http.StatusOK, response); err != nil {
		slog.Error("failed to encode json response", "error", err)
		s.respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}
}
