package game

import (
	"testing"
)

var valid5ShipsPlacement = []ShipPlacement{
	{
		IJNTaihou,
		AircraftCarrier,
		5, 9,
		Horizontal,
	}, {
		HMSBelfast,
		LightCruiser,
		0, 0,
		Vertical,
	}, {
		KMSBismarck,
		Battleship,
		6, 0,
		Horizontal,
	}, {
		KMSPrinzEugen,
		HeavyCruiser,
		0, 7,
		Vertical,
	}, {
		USSLaffey,
		Destroyer,
		5, 5,
		Horizontal,
	},
}

func TestMatch_PlaceShips(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() *Match
		playerID      string
		ShipPlacement []ShipPlacement
		wantErr       bool
		wantState     MatchState
	}{
		{
			"valid ship placement",
			func() *Match {
				m := NewMatch("666", "777")
				m.PlaceShips("777", valid5ShipsPlacement)
				return m
			},
			"666",
			valid5ShipsPlacement,
			false,
			StateInBattle,
		}, {
			"invalid match phase",
			func() *Match {
				m := NewMatch("666", "777")
				m.State = StateInBattle
				return m
			},
			"666",
			valid5ShipsPlacement,
			true,
			StateInBattle,
		}, {
			"invalid ship placement",
			func() *Match {
				return NewMatch("666", "777")
			},
			"666",
			[]ShipPlacement{},
			true,
			StatePreparation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.setup()
			err := m.PlaceShips(tt.playerID, tt.ShipPlacement)
			if (err != nil) != tt.wantErr {
				t.Errorf("test failed, error: %v, wantErr: %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if m.State != tt.wantState {
					t.Errorf("test failed, expected: %v, result: %v", tt.wantState, m.State)
				}
			}
		})
	}
}

func TestMatch_Fire(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() *Match
		playerID   string
		x, y       uint8
		wantErr    bool
		fireResult FireResult
	}{
		{
			"valid firing opponent ship",
			func() *Match {
				m := NewMatch("666", "777")
				m.PlaceShips("666", valid5ShipsPlacement)
				m.PlaceShips("777", valid5ShipsPlacement)
				return m
			},
			"666",
			7, 9,
			false,
			FireResult{
				"666", "777", 7, 9, true, "", "666",
			},
		}, {
			"invalid match phase",
			func() *Match {
				return NewMatch("666", "777")
			},
			"666",
			7, 9,
			true,
			FireResult{},
		}, {
			"invalid player turn",
			func() *Match {
				m := NewMatch("666", "777")
				m.PlaceShips("666", valid5ShipsPlacement)
				m.PlaceShips("777", valid5ShipsPlacement)
				return m
			},
			"777",
			7, 9,
			true,
			FireResult{},
		}, {
			"invalid firing placement, out of the bounds",
			func() *Match {
				m := NewMatch("666", "777")
				m.PlaceShips("666", valid5ShipsPlacement)
				m.PlaceShips("777", valid5ShipsPlacement)
				return m
			},
			"666",
			7, 10,
			true,
			FireResult{},
		}, {
			"invalid firing coordinate, already fired",
			func() *Match {
				m := NewMatch("666", "777")
				m.PlaceShips("666", valid5ShipsPlacement)
				m.PlaceShips("777", valid5ShipsPlacement)
				m.Fire("666", 7, 9)
				return m
			},
			"666",
			7, 9,
			true,
			FireResult{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.setup()
			res, err := m.Fire(tt.playerID, tt.x, tt.y)
			if (err != nil) != tt.wantErr {
				t.Errorf("test failed, error: %v, wantErr: %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if tt.fireResult.FiredBy != res.FiredBy {
					t.Errorf("test failed, expected: %s, result: %s", tt.playerID, res.FiredBy)
				}
				if tt.fireResult.Target != res.Target {
					t.Errorf("test failed, expected: %s, result: %s", tt.fireResult.Target, res.Target)
				}
				if tt.fireResult.Hit != res.Hit {
					t.Errorf("test failed, expected: %t, result: %t", tt.fireResult.Hit, res.Hit)
				}
				if tt.fireResult.NextTurn != res.NextTurn {
					t.Errorf("test failed, expected: %s, result: %s", tt.fireResult.NextTurn, res.NextTurn)
				}
			}
		})
	}
}

func TestMatch_Surrender(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() *Match
		playerID      string
		wantErr       bool
		wantState     MatchState
		wantWinnerID  string
		wantEndReason EndReason
	}{
		{
			"valid end match w/ surrender",
			func() *Match {
				m := NewMatch("666", "777")
				m.State = StateInBattle
				return m
			},
			"666",
			false,
			StateFinished,
			"777",
			OpponentSurrendered,
		}, {
			"invalid match phase",
			func() *Match {
				return NewMatch("666", "777")
			},
			"666",
			true,
			StatePreparation,
			"",
			"",
		}, {
			"invalid player id",
			func() *Match {
				m := NewMatch("666", "777")
				m.State = StateInBattle
				return m
			},
			"555",
			true,
			StatePreparation,
			"",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.setup()
			err := m.Surrender(tt.playerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("test failed, error: %v, wantErr: %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if m.State != tt.wantState {
					t.Errorf("test failed, expected: %v, result: %v", tt.wantState, m.State)
				}

				if m.WinnerID != tt.wantWinnerID {
					t.Errorf("test failed, expected: %s, result: %s", tt.wantWinnerID, m.WinnerID)
				}

				if m.EndReason != tt.wantEndReason {
					t.Errorf("test failed, expected: %s, result: %s", tt.wantEndReason, m.EndReason)
				}
			}
		})
	}
}

func TestMatch_Forfeit(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() *Match
		playerID      string
		wantErr       bool
		wantState     MatchState
		wantWinnerID  string
		wantEndReason EndReason
	}{
		{
			"valid end match w/ forfeit",
			func() *Match {
				m := NewMatch("666", "777")
				m.State = StateInBattle
				return m
			},
			"666",
			false,
			StateFinished,
			"777",
			OpponentDisconnectedTimeout,
		}, {
			"invalid match phase",
			func() *Match {
				return NewMatch("666", "777")
			},
			"666",
			true,
			StatePreparation,
			"",
			"",
		}, {
			"invalid player id",
			func() *Match {
				m := NewMatch("666", "777")
				m.State = StateInBattle
				return m
			},
			"555",
			true,
			StatePreparation,
			"",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.setup()
			err := m.Forfeit(tt.playerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("test failed, error: %v, wantErr: %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if m.State != tt.wantState {
					t.Errorf("test failed, expected: %v, result: %v", tt.wantState, m.State)
				}

				if m.WinnerID != tt.wantWinnerID {
					t.Errorf("test failed, expected: %s, result: %s", tt.wantWinnerID, m.WinnerID)
				}

				if m.EndReason != tt.wantEndReason {
					t.Errorf("test failed, expected: %s, result: %s", tt.wantEndReason, m.EndReason)
				}
			}
		})
	}
}
