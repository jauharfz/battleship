package game

import "fmt"

type Player struct {
	PlayerID string
	IsReady  bool
	Ships    []*Ship
	Grid     [GridSize][GridSize]Cell
}

func NewPlayer(playerID string) *Player {
	return &Player{
		PlayerID: playerID,
	}
}

// ClientSide Validate each Ship Position
// ServerSide ReValidate All Ship Position

// check len, ship assign exactly one, overlap betweeb ships
func (p *Player) validateAllShipPosition(ships []ShipPlacement) error {
	if len(ships) != len(shipSizes) {
		return fmt.Errorf("invalid number of ships")
	}

	tempShipType := make(map[ShipType]struct{})
	var tempGrid [GridSize][GridSize]Cell
	var tempShips []*Ship
	for i, shipPlacement := range ships {
		if _, ok := tempShipType[shipPlacement.Type]; ok {
			return fmt.Errorf("invalid ship type, already have %s type", shipPlacement.Type)
		}
		tempShipType[shipPlacement.Type] = struct{}{}
		ship, err := NewShip(shipPlacement.Name, shipPlacement.Type, shipPlacement.Direction, shipPlacement.X, shipPlacement.Y)
		if err != nil {
			return fmt.Errorf("failed to to place ship %v", err)
		}
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
		tempShips = append(tempShips, ship)
	}
	p.Grid = tempGrid
	p.Ships = tempShips
	p.IsReady = true
	return nil
}
