package server

import (
	"fmt"
	"log/slog"
	"net/http"
)

func (s *Server) EnterMatchmaking(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
	playerID := r.Context().Value(contextKeyPlayerID).(string)
	if playerID == s.Matchmaker.Waiting {
		slog.Warn(fmt.Sprintf("user trying to enter matchmaking twice lmao, userID: %s", playerID))
		s.respondWithError(w, http.StatusBadRequest, "Bad Request")
		return
	}
	s.Matchmaker.Pool <- playerID
	response := jsonResponse{
		Status: "queued",
		Data: map[string]string{
			"message": "Searching for an opponent",
		},
	}
	if err := s.respondWithJSON(w, http.StatusCreated, response); err != nil {
		slog.Error("failed to encode json response", "error", err)
		s.respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	// http.Redirect(w, r, "GET /api/v1/matchmaking/status", http.StatusTemporaryRedirect)
}

func (s *Server) CheckStatus(w http.ResponseWriter, r *http.Request) {
	var response jsonResponse
	//
	playerID := r.Context().Value(contextKeyPlayerID).(string)
	ticketID, err := s.TicketStore.Get(playerID)
	if err != nil {
		response = jsonResponse{
			Status: "searching",
		}
	} else {
		response = jsonResponse{
			Status: "matched",
			Data: map[string]string{
				"ws_token": ticketID,
			},
		}
	}
	// process a ticket issue w/ playerid && rooms
	if err := s.respondWithJSON(w, http.StatusOK, response); err != nil {
		slog.Error("failed to encode json response", "error", err)
		s.respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
}

// TODO: create a handler(?) which matchmaking by hardening room id -> matchmaking by room? -> issue?, : need to adjust register guest to provide default room id? create room lookup too?
