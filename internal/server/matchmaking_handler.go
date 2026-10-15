package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func (s *Server) EnterMatchmaking(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
	var roomID string

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&roomID); err != nil {
		slog.Warn("failed to decode json request", "error", err)
		s.respondWithError(w, http.StatusBadRequest, "bad request")
		return
	}

	playerID := r.Context().Value(contextKeyPlayerID).(string)
	if roomID == "" {
		s.Matchmaker.Pool <- playerID
		slog.Info("player queued to random matchmaker", "playerID", playerID)
	} else {
		opponent, exist := s.Matchmaker.CustomRooms[roomID]
		if exist && opponent.PlayerID != playerID {
			s.TicketStore.Issue(playerID, roomID)
			s.TicketStore.Issue(opponent.PlayerID, roomID)
			s.Matchmaker.DeleteCustomMatch(roomID, opponent.PlayerID)
			slog.Info("another player registered to the custom room, ticked issued", "playerID", playerID, "opponentID", opponent, "roomID", roomID)
		} else if !exist {
			s.Matchmaker.AssignCustomMatch(roomID, playerID)
			slog.Info("player registered to the custom room", "playerID", playerID, "roomID", roomID)
		}
	}
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
}

// TODO: handle lastChecked and reassign
func (s *Server) CheckStatus(w http.ResponseWriter, r *http.Request) {
	var response jsonResponse
	playerID := r.Context().Value(contextKeyPlayerID).(string)
	ticketID, err := s.TicketStore.Get(playerID)

	if err != nil {
		response = jsonResponse{
			Status: "searching",
		}
		if _, exist := s.Matchmaker.CustomPlayers[playerID]; exist {
			s.Matchmaker.UpdateLastChecked(playerID, true)
		} else if s.Matchmaker.Waiting != nil && s.Matchmaker.Waiting.PlayerID == playerID {
			s.Matchmaker.UpdateLastChecked(playerID, false)
		} else {
			response = jsonResponse{
				Status: "timedout",
			}
		}
	} else {
		slog.Info("matched", "playerID", playerID, "ticketID", ticketID)
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
