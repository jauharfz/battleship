package server

import (
	"github.com/jauharfz/battleship/internal/lobby"
)

type Server struct {
	Matchmaker *lobby.Matchmaker
}

func NewServer(matchmaker *lobby.Matchmaker) *Server {
	return &Server{
		Matchmaker: matchmaker,
	}
}
