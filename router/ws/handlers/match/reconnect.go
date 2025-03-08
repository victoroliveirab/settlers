package match

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core"
	"github.com/victoroliveirab/settlers/router/ws/entities"
)

func ReconnectPlayer(player *entities.GamePlayer) error {
	room := player.Room
	if room == nil {
		err := fmt.Errorf("Room not assigned to the player %s", player.ID)
		return err
	}

	game := room.Game
	if game == nil {
		err := fmt.Errorf("Game not assigned to the room", room.ID)
		return err
	}

	err := sendHydratePlayer(player)
	if err != nil {
		return err
	}

	// Not player's turn, nothing to do
	if game.CurrentRoundPlayer().ID != player.Username {
		return nil
	}

	if game.RoundType() == core.FirstRound || game.RoundType() == core.Regular || game.RoundType() == core.BetweenTurns {
		return nil // hydrate will take care of it
	}

	err = fmt.Errorf("Cannot reconnect player#%s during match setup: not known round type %s", player.Username, core.RoundTypeTranslation[game.RoundType()])
	return err
}
