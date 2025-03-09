package match

import (
	"fmt"

	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func handleNewRoad(player *entities.GamePlayer, message *types.WebSocketMessage) (bool, error) {
	payload, err := parseRoadBuildPayload(message.Payload)
	if err != nil {
		wsErr := sendNewRoadError(player.Connection, player.ID, err)
		return true, wsErr
	}

	edgeID := payload.edgeID
	room := player.Room
	game := room.Game
	err = game.BuildRoad(player.Username, edgeID)
	if err != nil {
		wsErr := sendNewRoadError(player.Connection, player.ID, err)
		return true, wsErr
	}

	logs := []string{fmt.Sprintf("%s just built a road.", player.Username)}
	err = sendNewRoadSuccess(player, edgeID, logs)
	room.EnqueueBroadcastMessage(buildNewRoadBroadcast(player.Username, edgeID, logs), []int64{player.ID}, nil)
	return true, err

}
