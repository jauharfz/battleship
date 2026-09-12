package room

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"uuid"

	"github.com/jauharfz/battleship/internal/game"
	"github.com/jauharfz/battleship/internal/protocol"
)

type Room struct {
	ID         string
	Clients    map[*Client]struct{}
	Match      game.Match
	Register   chan *Client
	Unregister chan *Client
	Ingress    chan Message
	Done       chan struct{}
}

// ingress -> handle message from client, carrying client & byte
// then validate the client, parse
func NewRoom() *Room {
	return &Room{
		ID:         uuid.New().String(),
		Clients:    make(map[*Client]struct{}, 2),
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

func (r *Room) Run() error {
	defer func() {
		for client := range r.Clients {
			delete(r.Clients, client)
			close(client.Send)
		}
	}()
	for {
		select {
		// future TODO: matchmaking using REST first, and send the client_id foreach client if room has been created(?)
		case client := <-r.Register:
			client.Room = r
			r.Clients[client] = struct{}{}
			if len(r.Clients) == 2 {
				var clientIDList [2]string
				idx := 0
				for client := range r.Clients {
					clientIDList[idx] = client.ID
					idx++
				}
				r.Match = *game.NewMatch(clientIDList[0], clientIDList[1])
				// temp for debugging purpose
				for client := range r.Clients {
					// data, _ := json.Marshal("ROOM CREATED SUCCESFULLY")
					if err := r.roomCreatedHandler(client); err != nil {
						return err
					}
				}
			}
		// this usually happen when the client is disconnected/timeout
		case client := <-r.Unregister:
			if _, exists := r.Clients[client]; exists {
				if len(r.Clients) == 0 {
					return nil
				}
				if r.Match.State == game.StateInBattle {
					if err := r.disconnectedTimeoutHandler(client, "finished"); err != nil {
						return err
					}
				}
				if r.Match.State == game.StatePreparation {
					if err := r.disconnectedTimeoutHandler(client, "aborted"); err != nil {
						return err
					}
				}
				return nil
			}
		// determine the type -> unmarshal by type -> process the request -> sendEvent/sendError
		case message := <-r.Ingress:
			if err := r.handleMessage(message.Sender, message.Payload); err != nil {
				r.sendError(message.Sender, err)
			}

		case <-r.Done:
			return nil
		}
	}
}

func (r *Room) handleMessage(sender *Client, payload InboundPayload) error {
	switch payload.Type {
	case string(protocol.CmdPlaceShips):
		var data []protocol.ShipPlacementDTO
		if err := json.Unmarshal(payload.Data, &data); err != nil {
			slog.Error("failed to marshaling payload place_ships data", "error", err, "client_id", sender.ID)
			return err
		}
		placeShips := make([]game.ShipPlacement, len(data))
		for idx, place := range data {
			var err error
			placeShips[idx], err = place.ToDomain()
			if err != nil {
				slog.Warn("something isn't right with client during place ship", "error", err, "client_id", sender.ID, "data", data)
				r.sendError(sender, err)
				return err
			}
		}

		if err := r.Match.PlaceShips(sender.ID, placeShips); err != nil {
			slog.Warn("something isn't right with client during place ship", "error", err, "client_id", sender.ID, "data", data)
			r.sendError(sender, err)
			break
		}
		if err := r.placeShipsHandler(sender); err != nil {
			return err
		}
		if err := r.playerReadyHandler(sender); err != nil {
			return err
		}
		if r.Match.State == game.StateInBattle {
			// whoever is ready first -> opponent
			opponent, err := r.getOpponent(sender)
			if err != nil {
				return err
			}

			if err := r.battleStartHandler(opponent); err != nil {
				return err
			}
		}

	case string(protocol.CmdCancelReady):
		if err := r.Match.CancelReady(sender.ID); err != nil {
			slog.Warn("something isn't right with client during cancel_ready", "error", err, "client_id", sender.ID)
			r.sendError(sender, err)
			return err
		}
		if err := r.cancelReadyHandler(sender); err != nil {
			return err
		}
	case string(protocol.CmdFire):
		var data protocol.FireData
		if err := json.Unmarshal(payload.Data, &data); err != nil {
			slog.Error("failed to marshaling payload fire data", "error", err)
			return err
		}
		result, err := r.Match.Fire(sender.ID, data.X, data.Y)
		if err != nil {
			slog.Warn("something isn't right with client during fire", "error", err, "client_id", sender.ID, "data", data)
			r.sendError(sender, err)
			break
		}
		if err := r.fireHandler(result); err != nil {
			return err
		}
		if r.Match.State == game.StateFinished {
			if err := r.AllShipsDestroyedHandler(sender); err != nil {
				return err
			}
		}

	case string(protocol.CmdSurrender):
		if err := r.Match.Surrender(sender.ID); err != nil {
			slog.Warn("something isn't right with client during surrender", "error", err, "client_id", sender.ID)
			r.sendError(sender, err)
			break
		}
		if err := r.surrenderHandler(sender); err != nil {
			return err
		}

	default:
		slog.Warn("something isn't right with client, payload type command unrecognized", "client_id", sender.ID, "payload", payload)
		r.sendError(sender, errors.New("Invalid Type Command"))
	}
	return nil
}

func (r *Room) sendError(client *Client, err error) error {
	payload := OutboundPayload{
		Type: string(protocol.EvtError),
		Data: fmt.Sprintf("%s", err),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		slog.Error("failed to marshaling   at send_error", "error", err, "client_id", client.ID)
		return err
	}
	select {
	case client.Send <- data:
	default:
		r.Unregister <- client
	}
	return nil
}

func (r *Room) sendEvent(client *Client, eventType protocol.Event, data any) error {
	payload := OutboundPayload{
		Type: string(eventType),
		Data: data,
	}
	data2, err := json.Marshal(payload)
	if err != nil {
		slog.Error("failed to marshaling outbound payload at send_event", "error", err, "client_id", client.ID, "payload", payload)
		return err
	}
	select {
	case client.Send <- data2:
	default:
		r.Unregister <- client
	}
	return nil
}

// TODO: handle either for the sender or the receiver

func (r *Room) roomCreatedHandler(sender *Client) error {
	message := protocol.RoomCreatedData{
		PlayerID: sender.ID,
	}
	if err := r.sendEvent(sender, protocol.EvtRoomCreated, message); err != nil {
		return err
	}
	return nil
}

func (r *Room) placeShipsHandler(sender *Client) error {
	message := protocol.PlaceShipsResultData{
		IsValid: true,
		Reason:  "its true ofc",
	}
	if err := r.sendEvent(sender, protocol.EvtPlaceShipsResult, message); err != nil {
		return err
	}
	return nil
}

func (r *Room) playerReadyHandler(sender *Client) error {
	return r.readyHandlerHelper(sender.ID, true, "player_ready")
}

func (r *Room) cancelReadyHandler(sender *Client) error {
	return r.readyHandlerHelper(sender.ID, false, "cancel_ready")
}

// TODO: fix the error msg based on isPlayerReady
func (r *Room) readyHandlerHelper(senderID string, isPlayerReady bool, callLoc string) error {
	message := protocol.PlayerReadyData{
		PlayerId:      senderID,
		IsPlayerReady: isPlayerReady,
	}
	for client := range r.Clients {
		if err := r.sendEvent(client, protocol.EvtPlayerReady, message); err != nil {
			return err
		}
	}
	return nil
}

func (r *Room) battleStartHandler(sender *Client) error {
	message := protocol.BattleStartData{
		Turn: r.Match.Players[r.Match.Turn].PlayerID,
	}
	for client := range r.Clients {
		if err := r.sendEvent(client, protocol.EvtBattleStart, message); err != nil {
			return err
		}
	}
	return nil
}

// cell empty -> != destroyed
// cell ship -> reply as is & check if state finished/all ship destroyed
// cell miss/hit/etc -> sendErr, (it's already handled by caller)
func (r *Room) fireHandler(result game.FireResult) error {
	message := protocol.FireResultData{
		FiredBy:       result.FiredBy,
		Target:        result.Target,
		Y:             result.Y,
		X:             result.X,
		Hit:           result.Hit,
		ShipDestroyed: string(result.ShipDestroyed),
		NextTurn:      result.NextTurn,
	}
	for client := range r.Clients {
		if err := r.sendEvent(client, protocol.EvtFireResult, message); err != nil {
			return err
		}
	}
	return nil
}

func (r *Room) gameOverHelper(client *Client, reason game.EndReason, errLoc, status string) error {
	var winnerID string
	if status == "finished" {
		winnerID = client.ID
	}
	message := protocol.GameOverData{
		Status:   status,
		WinnerId: winnerID,
		Reason:   string(reason),
	}
	if errLoc == "disconnected_timeout" {
		if err := r.sendEvent(client, protocol.EvtGameOver, message); err != nil {
			return err
		}
	} else {
		for client := range r.Clients {
			if err := r.sendEvent(client, protocol.EvtGameOver, message); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Room) AllShipsDestroyedHandler(client *Client) error {
	return r.gameOverHelper(client, game.AllEnemyShipsDestroyed, "all_ship_destroyed", "finished")
}

func (r *Room) surrenderHandler(sender *Client) error {
	opponent, err := r.getOpponent(sender)
	if err != nil {
		return err
	}
	return r.gameOverHelper(opponent, game.OpponentSurrendered, "surrender", "finished")
}

// TODO:create all end reason handler
func (r *Room) disconnectedTimeoutHandler(client *Client, status string) error {
	if status == "finished" {
		if err := r.Match.Forfeit(client.ID); err != nil {
			if err := r.sendError(client, err); err != nil {
				return err
			}
		}
	}
	opponent, err := r.getOpponent(client)
	if err != nil {
		return err
	}
	return r.gameOverHelper(opponent, game.OpponentDisconnectedTimeout, "disconnected_timeout", status)
}

func (r *Room) getOpponent(sender *Client) (*Client, error) {
	if _, ok := r.Clients[sender]; !ok {
		err := errors.New("invalid sender client")
		slog.Error("failed to get opponent", "error", err, "client_id", sender.ID)
		return nil, err
	}
	for client := range r.Clients {
		if client != sender {
			return client, nil
		}
	}
	err := errors.New("idk")
	slog.Error("an error occured during get opponent", "error", err)
	return nil, err
}
