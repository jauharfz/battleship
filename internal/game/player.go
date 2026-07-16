package game

type Player struct {
	PlayerID string
	IsReady  bool
	Ships    []Ship
	Grid     [GridSize][GridSize]Cell
}
