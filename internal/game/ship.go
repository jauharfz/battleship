package game

import (
	"fmt"
)

type ShipType string

const (
	AircraftCarrier ShipType = "aircraft_carrier"
	Battleship      ShipType = "battleship"
	HeavyCruiser    ShipType = "heavy_cruiser"
	LightCruiser    ShipType = "light_cruiser"
	Destroyer       ShipType = "destroyer"
)

var shipSizes = map[ShipType]uint8{
	AircraftCarrier: 5,
	Battleship:      4,
	HeavyCruiser:    3,
	LightCruiser:    3,
	Destroyer:       2,
}

type ShipName string

const (
	KMSPrinzEugen ShipName = "KMS_Prinz_Eugen"
	IJNTaihou     ShipName = "IJN_Taihou"
	HMSBelfast    ShipName = "HMS_Belfast"
	USSLaffey     ShipName = "USS_Laffey"
	KMSBismarck   ShipName = "KMS_Bismarck"
)

type Ship struct {
	Name      ShipName
	Type      ShipType
	Size      uint8
	X         uint8
	Y         uint8
	Direction Direction
	Health    uint8
}

type ShipPlacement struct {
	Name      ShipName
	Type      ShipType
	X         uint8
	Y         uint8
	Direction Direction
}

// TODO: Bound Checking X & Y
func NewShip(name ShipName, typ ShipType, direction Direction, x, y uint8) (*Ship, error) {
	size, ok := shipSizes[typ]
	if !ok {
		return nil, fmt.Errorf("unknown ship type, please input valid ship type")
	}

	ship := &Ship{
		Name: name, Type: typ, Direction: direction, X: x, Y: y, Size: size, Health: size,
	}
	err := ship.validate(size)
	if err != nil {
		return nil, err
	}

	return ship, nil
}

func (s *Ship) validate(size uint8) error {
	if s.X < 0 || s.Y < 0 || s.X >= GridSize || s.Y >= GridSize {
		return fmt.Errorf("coordinates out of the bounds")
	}

	switch s.Direction {
	case Horizontal:
		if s.X+size > GridSize {
			return fmt.Errorf("horizontal coordinates out of the bounds")
		}
	case Vertical:
		if s.Y+size > GridSize {
			return fmt.Errorf("vertical coordinates out of the bounds")
		}
	default:
		return fmt.Errorf("invalid direction")
	}
	return nil
}
