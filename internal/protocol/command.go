package protocol

import "github.com/jauharfz/battleship/internal/game"

type Command string

const (
	CmdPlaceShips  Command = "CMD_PLACE_SHIPS"
	CmdCancelReady Command = "CMD_CANCEL_READY"
	CmdFire        Command = "CMD_FIRE"
	CmdSurrender   Command = "CMD_SURRENDER"
)

type PlaceShipsData struct {
	Ships []game.ShipPlacement `json:"ships"`
}

type CancelReadyData struct{}

type FireData struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type SurrenderData struct{}
