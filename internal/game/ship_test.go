package game

import "testing"

func TestNewShip(t *testing.T) {
	tests := []struct {
		name      string
		shipName  ShipName
		shipType  ShipType
		direction Direction
		x, y      int
		wantErr   bool
	}{
		{
			"valid ship placement",
			KMSPrinzEugen,
			HeavyCruiser,
			Horizontal,
			7, 9,
			false,
		}, {
			"invalid ship type",
			KMSPrinzEugen,
			"random_type_ship",
			Horizontal,
			7, 9,
			true,
		}, {
			"invalid coordinates, out of the bounds",
			KMSPrinzEugen,
			HeavyCruiser,
			Horizontal,
			7, 10,
			true,
		}, {
			"invalid horizontal coordinates, ship placement out of the bounds",
			KMSPrinzEugen,
			HeavyCruiser,
			Horizontal,
			8, 9,
			true,
		}, {
			"invalid vertical coordinates, ship placement out of the bounds",
			KMSPrinzEugen,
			HeavyCruiser,
			Vertical,
			9, 8,
			true,
		}, {
			"invalid ship direction",
			KMSPrinzEugen,
			HeavyCruiser,
			2, // neither horizontal nor vertical
			7, 7,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewShip(tt.shipName, tt.shipType, tt.direction, tt.x, tt.y)
			if (err != nil) != tt.wantErr {
				t.Errorf("test failed, error: %v, wantErr: %v", err, tt.wantErr)
			}
		})
	}
}
