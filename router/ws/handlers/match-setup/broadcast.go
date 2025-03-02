package matchsetup

import (
	"fmt"
	"strings"

	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func buildSettlementSetupBuildSuccessBroadcast(builderID string, vertexID int, logs []string) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Type: "setup.settlement-build.success",
		Payload: map[string]interface{}{
			"settlement": map[string]interface{}{
				"id":    vertexID,
				"owner": builderID,
			},
			"logs": logs,
		},
	}
}

func buildRoadSetupBuildSuccessBroadcast(builderID string, edgeID int, logs []string) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Type: "setup.road-build.success",
		Payload: map[string]interface{}{
			"road": map[string]interface{}{
				"id":    edgeID,
				"owner": builderID,
			},
			"logs": logs,
		},
	}
}

func buildSetupPlayerRoundChangedBroadcast(room *entities.Room) *types.WebSocketMessage {
	playerRound := room.Game.CurrentRoundPlayer()
	return &types.WebSocketMessage{
		Type: "setup.player-round-changed",
		Payload: map[string]interface{}{
			"currentRoundPlayer": playerRound.ID,
		},
	}
}

func buildSetupPhaseOverBroadcast(room *entities.Room) *types.WebSocketMessage {
	logs := make([]string, 0)
	hands := make(map[string]map[string]int)

	game := room.Game
	for _, player := range game.Players() {
		playerHand := game.ResourceHandByPlayer(player.ID)
		hands[player.ID] = playerHand

		receivedResource := false
		var builder strings.Builder
		builder.WriteString(fmt.Sprintf("%s received: ", player.ID))

		for resource, quantity := range playerHand {
			if quantity > 0 {
				receivedResource = true
				builder.WriteString(fmt.Sprintf("%d %s ", quantity, resource))
			}
		}

		if receivedResource {
			logEntry := strings.TrimSuffix(builder.String(), " ")
			logs = append(logs, logEntry)
		} else {
			logs = append(logs, fmt.Sprintf("%s received: nothing", player.ID))
		}
	}

	logs = append(logs, "Game starting. May the best settler win!")

	return &types.WebSocketMessage{
		Type: "setup.end",
		Payload: map[string]interface{}{
			"hands": hands,
			"logs":  logs,
		},
	}
}
