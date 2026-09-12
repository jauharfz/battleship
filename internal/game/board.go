package game

type CellState uint8

const (
	CellEmpty CellState = iota
	CellShip
	CellHit
	CellMiss
)

type Cell struct {
	State     CellState
	ShipIndex uint8
}

const GridSize = 10

type Direction uint8

const (
	Horizontal Direction = iota
	Vertical
)
