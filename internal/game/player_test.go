package game

import "testing"

func TestPlayer_ValidateAllShipPosition(t *testing.T) {
	tests := []struct {
		name    string
		ships   []ShipPlacement
		wantErr bool
	}{
		{
			"valid for all ship placement",
			[]ShipPlacement{
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
			},
			false,
		}, {
			"invalid total number of ships",
			[]ShipPlacement{{}, {}, {}, {}, {}, {}},
			true,
		}, {
			"invalid ship type, duplicated ship type",
			[]ShipPlacement{
				{
					IJNTaihou,
					AircraftCarrier,
					5, 9,
					Horizontal,
				}, {
					IJNTaihou,
					AircraftCarrier,
					0, 0,
					Vertical,
				},
			},
			true,
		}, {
			"invalid ship placement",
			[]ShipPlacement{
				{
					IJNTaihou,
					AircraftCarrier,
					6, 9,
					Horizontal,
				},
			},
			true,
		}, {
			"invalid ship coordinate, already placed w/ another ship",
			[]ShipPlacement{
				{
					IJNTaihou,
					AircraftCarrier,
					5, 9,
					Horizontal,
				}, {
					HMSBelfast,
					LightCruiser,
					9, 7,
					Vertical,
				},
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewPlayer("666").validateAllShipPosition(tt.ships)
			if (err != nil) != tt.wantErr {
				t.Errorf("test failed, error: %v, wantErr: %v", err, tt.wantErr)
			}
		})
	}
}
