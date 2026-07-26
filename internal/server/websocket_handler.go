package server

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
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
	defer conn.Close()

	slog.Info("client connected", "remote_addr", conn.RemoteAddr())
	defer slog.Info("client disconnected", "remote_addr", conn.RemoteAddr())

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				slog.Error("failed to reading to a buffer", "error", err)
			}
			break
		}
		if err := conn.WriteMessage(messageType, p); err != nil {
			slog.Error("failed to writing the message", "error", err)
			break
		}
	}
}
