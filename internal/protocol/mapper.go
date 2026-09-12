package protocol

import (
	"errors"
	"strings"

	"github.com/jauharfz/battleship/internal/game"
)

func (dto ShipPlacementDTO) ToDomain() (game.ShipPlacement, error) {
	var dir game.Direction
	switch strings.ToLower(dto.Direction) {
	case "horizontal":
		dir = game.Horizontal
	case "vertical":
		dir = game.Vertical
	default:
		return game.ShipPlacement{}, errors.New("invalid direction")
	}

	return game.ShipPlacement{
		Name:      game.ShipName(dto.Name),
		Type:      game.ShipType(dto.Type),
		X:         dto.X,
		Y:         dto.Y,
		Direction: dir,
	}, nil
}
