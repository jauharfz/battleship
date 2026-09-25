package protocol

type Event string

const (
	EvtRoomCreated      Event = "EVENT_ROOM_CREATED"
	EvtPlaceShipsResult Event = "EVENT_PLACE_SHIPMENT_RESULT"
	EvtPlayerReady      Event = "EVENT_PLAYER_READY"
	EvtBattleStart      Event = "EVENT_BATTLE_START"
	EvtFireResult       Event = "EVENT_FIRE_RESULT"
	EvtError            Event = "EVENT_ERROR"
	EvtGameOver         Event = "EVENT_GAME_OVER"
)

type EndStatus string

const (
	Finished EndStatus = "FINISHED"
	Aborted  EndStatus = "ABORTED"
)

type RoomCreatedData struct {
	PlayerID string `json:"player_id"`
}

type PlaceShipsResultData struct {
	IsValid bool   `json:"is_valid"`
	Reason  string `json:"reason"`
}

type PlayerReadyData struct {
	PlayerId      string `json:"player_id"`
	IsPlayerReady bool   `json:"is_player_ready"`
}

type BattleStartData struct {
	Turn string `json:"turn"`
}

type FireResultData struct {
	FiredBy       string `json:"fired_by"`
	Target        string `json:"target"`
	Y             uint8  `json:"y"`
	X             uint8  `json:"x"`
	Hit           bool   `json:"hit"`
	ShipDestroyed string `json:"ship_destroyed,omitempty"`
	NextTurn      string `json:"next_turn"`
}

type GameOverData struct {
	Status   string `json:"status"`
	WinnerId string `json:"winner_id,omitempty"`
	Reason   string `json:"reason"`
}

// not yet/used, can be deleted
// type EventErrorData struct {
// 	Code    string `json:"code"`
// 	Message string `json:"message"`
// }
