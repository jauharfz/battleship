package room

import (
	"encoding/json"
	"log/slog"
	"uuid"

	"github.com/jauharfz/battleship/internal/game"
)

type Room struct {
	ID         string
	Clients    map[*Client]struct{}
	Register   chan *Client
	Unregister chan *Client
	Ingress    chan Message
	Done       chan struct{}
}

func NewRoom() *Room {
	return &Room{
		ID:         uuid.New().String(),
		Clients:    make(map[*Client]struct{}),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Ingress:    make(chan Message),
		Done:       make(chan struct{}),
	}
}

type payload struct {
	Type    string         `json:"type"`
	Payload map[string]any `json:"payload"`
}

func (r *Room) Run() {
	defer func() {
		for client := range r.Clients {
			close(client.Send)
			delete(r.Clients, client)
		}
	}()
	for {
		select {
		case client := <-r.Register:
			client.Room = r
			r.Clients[client] = struct{}{}
			if len(r.Clients) == 2 {
				jsonData, err := createMessage(string(game.EvtBattleStart), map[string]any{"room_id": r.ID})
				if err != nil {
					return
				}
				for client := range r.Clients {
					select {
					case client.Send <- []byte(jsonData):
					default:
						r.Unregister <- client
					}
				}
			}
		case client := <-r.Unregister:
			if _, exists := r.Clients[client]; exists {
				close(client.Send)
				delete(r.Clients, client)
				if len(r.Clients) == 0 {
					return
				}

				jsonData, err := createMessage(string(game.EvtGameOver), map[string]any{"room_id": r.ID})
				if err != nil {
					return
				}
				for client := range r.Clients {
					select {
					case client.Send <- []byte(jsonData):
					default:
						r.Unregister <- client
					}
				}
			}
		case message := <-r.Ingress:
			jsonData, err := createMessage("EVENT_CHAT", map[string]any{"room_id": r.ID, "client_id": message.Sender.ID, "message": message.Data})
			if err != nil {
				return
			}
			for opponent := range r.Clients {
				if opponent != message.Sender {
					select {
					case opponent.Send <- jsonData:
					default:
						r.Unregister <- opponent

					}
				}
			}
		case <-r.Done:
			return
		}
	}
}

func createMessage(msg string, kv map[string]any) ([]byte, error) {
	data := payload{
		Type:    msg,
		Payload: kv,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		slog.Error("failed to marshaling payload data", "error", err)
		return nil, err
	}
	return jsonData, nil
}
