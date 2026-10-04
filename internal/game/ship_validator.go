package game

import "fmt"

func ValidateAllShipPosition(ships []ShipPlacement) ([]*Ship, *[10][10]Cell, error) {
	if len(ships) != len(shipSizes) {
		return nil, nil, fmt.Errorf("invalid number of ships")
	}

	tempShipType := make(map[ShipType]struct{})
	var tempGrid [GridSize][GridSize]Cell
	var tempShips []*Ship
	for i, shipPlacement := range ships {
		if _, ok := tempShipType[shipPlacement.Type]; ok {
			return nil, nil, fmt.Errorf("invalid ship type, already have %s type", shipPlacement.Type)
		}
		tempShipType[shipPlacement.Type] = struct{}{}
		ship, err := NewShip(shipPlacement.Name, shipPlacement.Type, shipPlacement.Direction, shipPlacement.X, shipPlacement.Y)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to to place ship %v", err)
		}
		tempShips = append(tempShips, ship)
		if err := validateGrid(ship, &tempGrid, i); err != nil {
			return nil, nil, err
		}
	}
	return tempShips, &tempGrid, nil
}

func validateGrid(ship *Ship, tempGrid *[10][10]Cell, i int) error {
	for idx := uint8(0); idx < shipSizes[ship.Type]; idx++ {
		var currX, currY uint8
		if ship.Direction == Horizontal {
			currX = ship.X + idx
			currY = ship.Y
		} else {
			currX = ship.X
			currY = ship.Y + idx
		}

		if tempGrid[currY][currX].State != CellEmpty {
			return fmt.Errorf("invalid coordinate for %s, already placed w/ another ships", ship.Name)
		}

		tempGrid[currY][currX] = Cell{
			State:     CellShip,
			ShipIndex: uint8(i),
		}

	}
	return nil
}
