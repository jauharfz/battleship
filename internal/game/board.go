package game

type CellState string

const (
	CellEmpty CellState = ""
	CellShip  CellState = "Ship"
	CellHit   CellState = "Hit"
	CellMiss  CellState = "Miss"
)

type Cell struct {
	State     CellState
	ShipIndex int //using 1-based for ShipIndex
}

type Direction int

const (
	Horizontal Direction = iota
	Vertical
)

const maxPlacement = 5
const GridSize = 10
