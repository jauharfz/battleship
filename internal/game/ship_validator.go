package game

import "fmt"

func ValidateAllShipPosition(ships []ShipPlacement) ([]*Ship, *[10][10]Cell, error) {
	if len(ships) != exactShipsNumber {
		return nil, nil, fmt.Errorf("invalid number of ships")
	}

	// max for each ship type is 2 -> total 10
	tempShipType := make(map[ShipType]uint8)
	var tempGrid [GridSize][GridSize]Cell
	var tempShips []*Ship
	for i, shipPlacement := range ships {
		if num := tempShipType[shipPlacement.Type]; num >= 2 {
			return nil, nil, fmt.Errorf("invalid ship type, max capapity for %s type", shipPlacement.Type)
		}
		tempShipType[shipPlacement.Type]++
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
