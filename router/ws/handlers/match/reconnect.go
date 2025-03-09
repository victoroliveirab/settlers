package match

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core"
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func ReconnectPlayer(room *entities.Room, playerID int64, conn *types.WebSocketConnection, onDisconnect func(player *entities.GamePlayer)) (*entities.GamePlayer, error) {
	player, err := room.ReconnectPlayer(playerID, conn, func(player *entities.GamePlayer) {
		room.RemovePlayer(playerID)
	})

	game := room.Game
	if game == nil {
		err := fmt.Errorf("Game not assigned to the room", room.ID)
		return nil, err
	}

	err = sendHydratePlayer(player)
	if err != nil {
		return nil, err
	}

	// Not player's turn, nothing to do
	if game.CurrentRoundPlayer().ID != player.Username {
		return player, nil
	}

	if game.RoundType() == core.FirstRound || game.RoundType() == core.Regular || game.RoundType() == core.BetweenTurns {
		return player, nil // hydrate will take care of it
	}

	err = fmt.Errorf("Cannot reconnect player#%s during match: not known round type %s", player.Username, core.RoundTypeTranslation[game.RoundType()])
	return nil, err
}
