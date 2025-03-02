package match

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core"
	"github.com/victoroliveirab/settlers/router/ws/entities"
)

func TryReconnectPlayer(player *entities.GamePlayer) error {
	err := player.Room.ReconnectPlayer(player)
	if err != nil {
		return err
	}

	err = sendHydratePlayer(player)
	if err != nil {
		return err
	}
	room := player.Room
	room.EnqueueBroadcastMessage(buildPlayerReconnectedBroadcast(player), []int64{player.ID}, nil)

	game := room.Game
	// Not player's turn, nothing to do
	if game.CurrentRoundPlayer().ID != player.Username {
		return nil
	}

	if game.RoundType() == core.SetupSettlement1 || game.RoundType() == core.SetupSettlement2 {
		err := SendBuildSetupSettlementRequest(player)
		return err
	} else if game.RoundType() == core.SetupRoad1 || game.RoundType() == core.SetupRoad2 {
		err := SendBuildSetupRoadRequest(player)
		return err
	}

	err = fmt.Errorf("Cannot reconnect player#%s: not known round type %s", player.Username, core.RoundTypeTranslation[game.RoundType()])
	return err
}
