package game

type Player struct {
	PlayerID string
	IsReady  bool
	Ships    []*Ship
	Grid     *[GridSize][GridSize]Cell
}

func NewPlayer(playerID string) *Player {
	return &Player{
		PlayerID: playerID,
	}
}

// ClientSide Validate each Ship Position
// ServerSide ReValidate All Ship Position

// check len, ship assign exactly one, overlap betweeb ships
func (p *Player) SetShips(ships []ShipPlacement) error {
	Ships, Grid, err := ValidateAllShipPosition(ships)
	if err != nil {
		return err
	}
	p.Grid = Grid
	p.Ships = Ships
	p.IsReady = true
	return nil
}
