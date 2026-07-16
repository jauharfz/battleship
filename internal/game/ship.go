package game

var shipSizes = map[string]int{
	"aircraft_carrier": 5,
	"battleship":       4,
	"heavy_cruiser":    3,
	"light_cruiser":    3,
	"destroyer":        2,
}

type Ship struct {
	Name      string
	Type      string
	Size      int
	X         int
	Y         int
	Direction Direction
	Health    int
}
