package room

import (
	"encoding/json"
	"log/slog"
	"uuid"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Send chan []byte
	Room *Room
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		ID:   uuid.New().String(),
		Conn: conn,
		Send: make(chan []byte, 256),
	}
}

func (c *Client) ReadPump() {
	defer func() {
		if c.Room != nil {
			select {
			case c.Room.Unregister <- c:
			default:
			}
		} else {
			close(c.Send)
		}
		c.Conn.Close()
	}()
	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			slog.Error("failed to read message from message reader", "error", err)
			break
		}

		if c.Room != nil {
			var payload InboundPayload
			err = json.Unmarshal(p, &payload)
			if err != nil {
				slog.Error("failed to unmarshal byte from message reader", "error", err)
			}
			c.Room.Ingress <- Message{c, payload}
		}
	}
}

func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()
	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			slog.Error("failed to write message to websocket connection", "error", err)
			break
		}
	}
}
