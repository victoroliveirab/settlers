package match

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core"
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

type discardCardsRequestPayload struct {
	Resources map[string]int `json:"resources"`
}

func handleDiscardCards(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	payload, err := utils.ParseJsonPayload[discardCardsRequestPayload](message)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	room := player.Room
	game := room.Game
	err = game.DiscardPlayerCards(player.Username, payload.Resources)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	formattedResources := utils.FormatResources(payload.Resources)
	logs := []string{fmt.Sprintf("%s discarded %s", player.Username, formattedResources)}

	if game.RoundType() == core.MoveRobberDue7 {
		room.EnqueueBulkUpdate(
			UpdatePlayerHand,
			UpdateResourceCount,
			UpdateDiscardPhase,
			UpdateRobberMovement,
			UpdateLogs(logs),
		)
	} else {
		room.EnqueueBulkUpdate(
			UpdatePlayerHand,
			UpdateResourceCount,
			UpdateDiscardPhase,
			UpdateLogs(logs),
		)
	}

	return true, nil
}
