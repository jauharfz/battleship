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

var shipSizes = map[ShipType]int{
	AircraftCarrier: 5,
	Battleship:      4,
	HeavyCruiser:    3,
	LightCruiser:    3,
	Destroyer:       2,
}

type Ship struct {
	Name      string
	Type      ShipType
	Size      int
	X         int
	Y         int
	Direction Direction
	Health    int
}

type ShipPlacement struct {
	Name      string
	Type      ShipType
	X         int
	Y         int
	Direction Direction
}

// TODO: Bound Checking X & Y
func NewShip(name string, typ ShipType, direction Direction, x, y int) (*Ship, error) {
	if _, ok := shipSizes[typ]; !ok {
		return nil, fmt.Errorf("unknown ship type, please input valid ship type")
	}
	size := shipSizes[typ]

	if x < 0 || y < 0 || x >= GridSize || y >= GridSize {
		return nil, fmt.Errorf("coordinates out of the bounds")
	}

	switch direction {
	case Horizontal:
		if x+size > GridSize {
			return nil, fmt.Errorf("horizontal coordinates out of the bounds")
		}
	case Vertical:
		if y+size > GridSize {
			return nil, fmt.Errorf("vertical coordinates out of the bounds")
		}
	default:
		return nil, fmt.Errorf("invalid direction")
	}

	return &Ship{
		Name: name, Type: typ, Direction: direction, X: x, Y: y, Size: size, Health: size,
	}, nil
}
