package match

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core"
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func handleDiscardCards(player *entities.GamePlayer, message *types.WebSocketMessage) (bool, error) {
	payload, err := parseDiscardCardsPayload(message.Payload)
	if err != nil {
		wsErr := sendDiscardCardsError(player.Connection, player.ID, err)
		return true, wsErr
	}

	room := player.Room
	game := room.Game
	err = game.DiscardPlayerCards(player.Username, payload.resources)
	if err != nil {
		wsErr := sendDiscardCardsError(player.Connection, player.ID, err)
		return true, wsErr
	}

	formattedResources := utils.FormatResources(payload.resources)
	logs := []string{fmt.Sprintf("%s discarded %s", player.Username, formattedResources)}

	err = sendDiscardCardsSuccess(room, player)
	if err != nil {
		return true, err
	}

	room.EnqueueBroadcastMessage(buildDiscardedCardsBroadcast(room, logs), []int64{}, func() {
		if game.RoundType() == core.MoveRobberDue7 {
			room.EnqueueBroadcastMessage(buildMoveRobberDueTo7Broadcast(room), []int64{}, nil)
		}
	})
	return true, nil

}
