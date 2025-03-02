package reconnect

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core"
	"github.com/victoroliveirab/settlers/router/ws/entities"
	matchsetup "github.com/victoroliveirab/settlers/router/ws/handlers/match-setup"
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
		err := matchsetup.SendBuildSetupSettlementRequest(player)
		return err
	} else if game.RoundType() == core.SetupRoad1 || game.RoundType() == core.SetupRoad2 {
		err := matchsetup.SendBuildSetupRoadRequest(player)
		return err
	} else if game.RoundType() == core.FirstRound || game.RoundType() == core.Regular {
		return nil // hydrate will take care of it
	}

	err = fmt.Errorf("Cannot reconnect player#%s: not known round type %s", player.Username, core.RoundTypeTranslation[game.RoundType()])
	return err
}
