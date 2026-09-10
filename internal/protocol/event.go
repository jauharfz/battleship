package protocol

import "github.com/jauharfz/battleship/internal/game"

type Event string

const (
	EvtPlaceShipsResult Event = "EVENT_PLACE_SHIPMENT_RESULT"
	EvtPlayerReady      Event = "EVENT_PLAYER_READY"
	EvtBattleStart      Event = "EVENT_BATTLE_START"
	EvtFireResult       Event = "EVENT_FIRE_RESULT"
	EvtError            Event = "EVENT_ERROR"
	EvtGameOver         Event = "EVENT_GAME_OVER"
)

type PlaceShipsResultData struct {
	IsValid bool   `json:"is_valid"`
	Reason  string `json:"reason"`
}

type PlayerReadyData struct {
	PlayerId      string `json:"player_id"`
	IsPlayerReady bool   `json:"is_player_ready"`
}

type BattleStartData struct {
	FirstTurn string `json:"first_turn"`
}

type FireResultData struct {
	FiredBy       string
	Target        string
	Y             int
	X             int
	Hit           bool
	ShipDestroyed game.ShipName
	NextTurn      string
}

type GameOverData struct {
	Status   string `json:"status"`
	WinnerId string `json:"winner_id"`
	Reason   string `json:"reason"`
}

// not yet/used, can be deleted
// type EventErrorData struct {
// 	Code    string `json:"code"`
// 	Message string `json:"message"`
// }
