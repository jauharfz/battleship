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

// TODO: fix the older serv ver, don't forget to consume the ticket
func (s *Server) serveWs(w http.ResponseWriter, r *http.Request) {
	// consume && validate before upgrade
	ticketID := r.URL.Query().Get("ticket")
	ticket, err := s.TicketStore.Consume(ticketID)
	if err != nil {
		slog.Warn("an error occured during consume ticket", "error", err)
		http.Error(w, "an error occured during consume ticket", http.StatusUnauthorized)
		return
	}

	slog.Info("ticket consumed", "ticket", ticket)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("failed to upgrade to websocket", "error", err)
		return
	}
	client := room.NewClient(conn)

	slog.Info("client connected", "remote_addr", conn.RemoteAddr())
	defer slog.Info("client disconnected", "remote_addr", conn.RemoteAddr())

	s.RoomManager.RegisterClient(ticket.RoomID, client)
	go client.WritePump()
	client.ReadPump()
}
