package server

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jauharfz/battleship/internal/room"
)

// NOTE: change if its Cross-Origin in prod
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// TODO: defer conn, rmv http.err
func (s *Server) serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("failed to upgrade to websocker", "error", err)
		return
	}
	client := room.NewClient(conn)
	defer func() {
		if client.Room == nil {
			s.Matchmaker.Dequeue()
		}
		conn.Close()
	}()

	slog.Info("client connected", "remote_addr", conn.RemoteAddr())
	defer slog.Info("client disconnected", "remote_addr", conn.RemoteAddr())

	go client.WritePump()
	s.Matchmaker.Enqueue(client)
	client.ReadPump()
}
