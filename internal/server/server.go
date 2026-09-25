package server

import (
	"github.com/jauharfz/battleship/internal/auth"
	"github.com/jauharfz/battleship/internal/lobby"
	"github.com/jauharfz/battleship/internal/room"
)

const maxBytes = 1_048_576

type Server struct {
	JWTManager  *auth.JWTManager
	TicketStore *auth.TicketStore
	Matchmaker  *lobby.Matchmaker
	RoomManager *room.RoomManager
}

func NewServer(jwtManager *auth.JWTManager, ticketStore *auth.TicketStore, matchmaker *lobby.Matchmaker) *Server {
	return &Server{
		JWTManager:  jwtManager,
		TicketStore: ticketStore,
		Matchmaker:  matchmaker,
		RoomManager: room.NewRoomManager(),
	}
}
