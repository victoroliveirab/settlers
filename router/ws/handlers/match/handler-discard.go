package match

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core"
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

type discardCardsPayload struct {
	resources map[string]int
}

func parseDiscardCardsPayload(payload map[string]interface{}) (*discardCardsPayload, error) {
	resourcesRaw, ok := payload["resources"]
	if !ok {
		return nil, fmt.Errorf("missing 'resources' key in payload")
	}

	resources := make(map[string]int)

	for key, value := range resourcesRaw.(map[string]interface{}) {
		num, ok := value.(float64)
		if !ok {
			err := fmt.Errorf("malformed data for key %s", key)
			return nil, err
		}
		resources[key] = int(num)
	}

	return &discardCardsPayload{resources: resources}, nil
}

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
