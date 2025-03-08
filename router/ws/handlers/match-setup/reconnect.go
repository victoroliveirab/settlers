package matchsetup

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

	if game.RoundType() == core.SetupSettlement1 || game.RoundType() == core.SetupSettlement2 {
		err := SendBuildSetupSettlementRequest(player)
		return player, err
	} else if game.RoundType() == core.SetupRoad1 || game.RoundType() == core.SetupRoad2 {
		err := SendBuildSetupRoadRequest(player)
		return player, err
	}

	err = fmt.Errorf("Cannot reconnect player#%s during match setup: not known round type %s", player.Username, core.RoundTypeTranslation[game.RoundType()])
	return player, err
}
