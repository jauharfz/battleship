package game

type CellState int

const (
	CellEmpty CellState = iota
	CellShip
	CellHit
	CellMiss
)

type Cell struct {
	State     CellState
	ShipIndex int
}

type Direction int

const (
	Horizontal Direction = iota
	Vertical
)

const GridSize = 10
