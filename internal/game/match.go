package game

import (
	"errors"
)

const maxPlayers = 2

type MatchState int

const (
	StatePreparation MatchState = iota
	StateInBattle
	StateFinished
	StateAborted
)

var (
	ErrAlreadyShot      = errors.New("ALREADY_SHOT")
	ErrNotYourTurn      = errors.New("NOT_YOUR_TURN")
	ErrInvalidPlacement = errors.New("INVALID_PLACEMENT")
	ErrInvalidPhase     = errors.New("INVALID_PHASE")
	ErrUnknownError     = errors.New("UNKNOWN_ERROR")
)

type Event string

const (
	EvtPlaceShipsResult Event = "EVENT_PLACE_SHIPMENT_RESULT"
	EvtOpponentReady    Event = "EVENT_OPPONENT_READY"
	EvtBattleStart      Event = "EVENT_BATTLE_START"
	EvtFireResult       Event = "EVENT_FIRE_RESULT"
	EvtError            Event = "EVENT_ERROR"
	EvtGameOver         Event = "EVENT_GAME_OVER"
)

type EndReason string

const (
	AllEnemyShipsDestroyed      EndReason = "ALL_ENEMY_SHIP_DESTROYED"
	OpponentDisconnectedTimeout EndReason = "OPPONENT_DISCONNECTED_TIMEOUT"
	OpponentSurrendered         EndReason = "OPPONENT_SURRENDERED"
)

type FireResult struct {
	FiredBy       string
	Target        string
	Y             int
	X             int
	Hit           bool
	ShipDestroyed ShipName
	NextTurn      string
}

type Match struct {
	State     MatchState
	players   [maxPlayers]*Player
	turn      int
	WinnerID  string
	EndReason EndReason
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
		return ErrInvalidPhase
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

// fix: validate the turn , if ship destroyecd,
// copy by value,validate x&y
func (m *Match) Fire(playerID string, x, y int) (FireResult, error) {
	if m.State != StateInBattle {
		return FireResult{}, ErrInvalidPhase
	}

	if m.players[m.turn].PlayerID != playerID {
		return FireResult{}, ErrNotYourTurn
	}

	if x >= GridSize || x < 0 || y >= GridSize || y < 0 {
		return FireResult{}, ErrInvalidPlacement
	}

	var nextTurn int
	if m.turn == 0 {
		nextTurn = 1
	} else {
		nextTurn = 0
	}
	opponent := m.players[nextTurn]
	opponentGrid := &opponent.Grid[y][x]
	var err error
	var hit bool
	var destroyedShipName ShipName
	switch opponentGrid.State {
	case CellEmpty:
		opponentGrid.State = CellMiss
	case CellShip:
		opponentShip := opponent.Ships[opponentGrid.ShipIndex]
		opponentGrid.State = CellHit
		hit = true
		opponentShip.Health--
		if opponentShip.Health == 0 {
			destroyedShipName = opponent.Ships[opponentGrid.ShipIndex].Name
		}
		if m.isAllShipDestroyed(nextTurn) {
			m.State = StateFinished
			m.WinnerID = m.players[m.turn].PlayerID
			m.EndReason = AllEnemyShipsDestroyed
		}
		nextTurn = m.turn
	case CellHit:
		err = ErrAlreadyShot
	case CellMiss:
		err = ErrAlreadyShot
	default:
		err = ErrUnknownError
	}
	if err != nil {
		return FireResult{}, err
	}
	m.turn = nextTurn

	return FireResult{playerID, opponent.PlayerID, y, x, hit, destroyedShipName, m.players[nextTurn].PlayerID}, nil
}

func (m *Match) Surrender(playerID string) error {
	return m.endMatchPrematurely(playerID, OpponentSurrendered)
}

func (m *Match) Forfeit(playerID string) error {
	return m.endMatchPrematurely(playerID, OpponentDisconnectedTimeout)
}

// getPlayerIndex
func (m *Match) getPlayerIndex(playerID string) (int, error) {
	if m.players[0].PlayerID == playerID {
		return 0, nil
	}
	if m.players[1].PlayerID == playerID {
		return 1, nil
	}
	return -1, ErrUnknownError
}

// isAllShipDestroyed
func (m *Match) isAllShipDestroyed(playerIdx int) bool {
	for _, ship := range m.players[playerIdx].Ships {
		if ship.Health != 0 {
			return false
		}
	}
	return true
}

func (m *Match) endMatchPrematurely(playerID string, endReason EndReason) error {
	{
		if m.State != StateInBattle {
			return ErrInvalidPhase
		}

		idx, err := m.getPlayerIndex(playerID)
		if err != nil {
			return err
		}

		var winnerIdx int
		if idx == 0 {
			winnerIdx = 1
		} else {
			winnerIdx = 0
		}

		m.State = StateFinished
		m.WinnerID = m.players[winnerIdx].PlayerID
		m.EndReason = endReason

		return nil
	}
}
