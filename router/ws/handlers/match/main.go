package match

import (
	"fmt"

	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func TryHandle(player *entities.GamePlayer, message *types.WebSocketMessage) (bool, error) {
	switch message.Type {
	case "setup.new-settlement":
		payload, err := parseSettlementBuildPayload(message.Payload)
		if err != nil {
			wsErr := sendSettlementSetupBuildError(player.Connection, player.ID, err)
			return true, wsErr
		}

		vertexID := payload.vertexID
		room := player.Room
		game := room.Game
		err = game.BuildSettlement(player.Username, vertexID)
		if err != nil {
			wsErr := sendSettlementSetupBuildError(player.Connection, player.ID, err)
			return true, wsErr
		}

		room.EnqueueBroadcastMessage(buildSettlementSetupBuildSuccessBroadcast(player.Username, vertexID, []string{fmt.Sprintf("%s just built a settlement.", player.Username)}), []int64{}, func() {
			err := SendBuildSetupRoadRequest(player)
			if err != nil {
				// TODO: handle player disconnect
				fmt.Println(err)
			}
		})
		return true, nil
	case "setup.new-road":
		payload, err := parseRoadBuildPayload(message.Payload)
		if err != nil {
			wsErr := sendRoadSetupBuildError(player.Connection, player.ID, err)
			return true, wsErr
		}

		edgeID := payload.edgeID
		room := player.Room
		game := room.Game
		err = game.BuildRoad(player.Username, edgeID)
		if err != nil {
			wsErr := sendRoadSetupBuildError(player.Connection, player.ID, err)
			return true, wsErr
		}

		room.EnqueueBroadcastMessage(buildRoadSetupBuildSuccessBroadcast(player.Username, edgeID, []string{fmt.Sprintf("%s just built a road.", player.Username)}), []int64{}, func() {
			nextRoundPlayer := game.CurrentRoundPlayer()
			for _, participant := range room.Participants {
				if participant.Player != nil && participant.Player.Username == nextRoundPlayer.ID {
					err := SendBuildSetupSettlementRequest(participant.Player)
					if err != nil {
						// TODO handle this err properly
						fmt.Println("err")
						fmt.Println(err)
					}
					break
				}
			}
		})
		return true, nil
	default:
		return false, nil
	}
}
