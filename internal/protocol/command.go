package protocol

type Command string

const (
	CmdPlaceShips  Command = "CMD_PLACE_SHIPS"
	CmdCancelReady Command = "CMD_CANCEL_READY"
	CmdFire        Command = "CMD_FIRE"
	CmdSurrender   Command = "CMD_SURRENDER"
)

type PlaceShipsData struct {
	Ships []ShipPlacementDTO `json:"ships"`
}

type FireData struct {
	X uint8 `json:"x"`
	Y uint8 `json:"y"`
}

type ShipPlacementDTO struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	X         uint8  `json:"x"`
	Y         uint8  `json:"y"`
	Direction string `json:"direction"`
}
