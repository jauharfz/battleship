package game

type MatchState int

const (
	StatePreparation MatchState = iota
	StateInBattle
	StateFinished
	StateAborted
)

type Match struct {
	State    MatchState
	Player1  Player
	Player2  Player
	TurnP1   bool
	WinnerID string
}
