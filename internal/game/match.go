package game

import "fmt"

const maxPlayers = 2

type MatchState int

const (
	StatePreparation MatchState = iota
	StateInBattle
	StateFinished
	StateAborted
)

type FireResult struct {
	FiredBy       string
	Target        string
	X             int
	Y             int
	Hit           bool
	ShipDestroyed *ShipType // ptr bcs default null
	NextTurn      string
}

type Match struct {
	State     MatchState
	players   [maxPlayers]*Player
	turn      int
	WinnerID  string
	EndReason string
}

func NewMatch(p1ID, p2ID string) *Match {
	player1 := NewPlayer(p1ID)
	player2 := NewPlayer(p2ID)
	return &Match{
		State:   StatePreparation,
		players: [maxPlayers]*Player{player1, player2},
		turn:    0,
	}
}

// TODO:
// PlaceShips,Fire,Surrender,Forfeit
// ShipPlacement,FireResult,
// validate ships di player/match
//

// ship -> isvalid ship
// player -> is allship correct direction
// match -> is 2 player ready?
//
// check state, p1id -> ready, p2id -> change state
func (m *Match) PlaceShips(playerID string, shipPlacements []ShipPlacement) error {
	if m.State != StatePreparation {
		return fmt.Errorf("invalid state")
	}
	idx, err := m.getPlayerIndex(playerID)
	if err != nil {
		return err
	}
	err = m.players[idx].validateAllShipPosition(shipPlacements)
	if err != nil {
		return err
	}
	if m.players[0].IsReady && m.players[1].IsReady {
		m.State = StateInBattle
	}

	return nil
}

func (m *Match) Fire(playerID string, x, y int) (FireResult, error) {
	return FireResult{}, nil
}

func (m *Match) Surrender(playerID string) error {
	return nil
}

func (m *Match) Forfeit(playerID string) error {
	return nil
}

// getPlayerIndex
func (m *Match) getPlayerIndex(playerID string) (int, error) {
	if m.players[0].PlayerID == playerID {
		return 0, nil
	}
	if m.players[1].PlayerID == playerID {
		return 1, nil
	}
	return -1, fmt.Errorf("invalid player id %s", playerID)
}
