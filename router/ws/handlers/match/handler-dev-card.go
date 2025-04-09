package match

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core"
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

type devCardClickRequestPayload struct {
	Kind string `json:"kind"`
}

func handleBuyDevCard(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	room := player.Room
	game := room.Game

	err := game.BuyDevelopmentCard(player.Username)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	room.EnqueueBulkUpdate(
		UpdateResourceCount,
		UpdatePlayerHand,
		UpdatePlayerDevHand,
		UpdatePlayerDevHandPermissions,
		UpdateBuyDevelopmentCard,
		UpdateLogs([]string{fmt.Sprintf("%s bought a [dev][/dev] card", player.Username)}),
	)
	return true, nil
}

func handleDevCardClick(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	payload, err := utils.ParseJsonPayload[devCardClickRequestPayload](message)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	room := player.Room
	game := room.Game
	err = game.UseDevelopmentCard(player.Username, payload.Kind)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	if game.RoundType() == core.GameOver {
		panic("HANDLE ME")
	} else if game.RoundType() == core.MoveRobberDueKnight {
		room.EnqueueBulkUpdate(
			UpdateRobberMovement,
			UpdatePlayerDevHand,
			UpdatePlayerDevHandPermissions,
			UpdatePass,
			UpdateTrade,
			UpdateLogs([]string{fmt.Sprintf("%s used Knight card", player.Username)}),
		)
	} else if game.RoundType() == core.MonopolyPickResource {
		panic("HANDLE ME")
	} else if game.RoundType() == core.BuildRoad1Development || game.RoundType() == core.BuildRoad2Development {
		room.EnqueueBulkUpdate(
			UpdateEdgeState,
			UpdatePlayerDevHand,
			UpdatePlayerDevHandPermissions,
			UpdatePass,
			UpdateTrade,
			UpdateLogs([]string{fmt.Sprintf("%s used Road Building card", player.Username)}),
		)
	} else if game.RoundType() == core.YearOfPlentyPickResources {
		panic("HANDLE ME")
	} else {
		room.EnqueueBulkUpdate(
			UpdatePass,
			UpdateTrade,
			UpdatePlayerDevHand,
			UpdatePlayerDevHandPermissions,
		)
	}

	return true, nil
}
