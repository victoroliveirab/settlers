package prematch

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core/packages/round"
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/handlers/match"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func TryHandle(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	switch message.Type {
	case "room.update-capacity":
		requestPayload, err := utils.ParseJsonPayload[roomUpdateCapacityRequestPayload](message)
		if err != nil {
			wsErr := player.WriteJsonError(message.Type, err)
			return true, wsErr
		}

		room := player.Room
		err = room.UpdateSize(player, requestPayload.Capacity)
		if err != nil {
			wsErr := player.WriteJsonError(message.Type, err)
			return true, wsErr
		}

		room.EnqueueOutgoingMessage(BuildRoomMessage(room, fmt.Sprintf("%s.success", message.Type)), nil, nil)
		return true, nil
	case "room.update-param":
		requestPayload, err := utils.ParseJsonPayload[roomUpdateParamRequestPayload](message)
		if err != nil {
			wsErr := player.WriteJsonError(message.Type, err)
			return true, wsErr
		}

		room := player.Room
		key := requestPayload.Key
		value := requestPayload.Value
		err = room.UpdateParam(player, key, value)
		if err != nil {
			wsErr := player.WriteJsonError(message.Type, err)
			return true, wsErr
		}

		room.EnqueueOutgoingMessage(BuildRoomMessage(room, fmt.Sprintf("%s.success", message.Type)), nil, nil)
		return true, nil
	case "room.player-change-color":
		requestPayload, err := utils.ParseJsonPayload[roomPlayerChangeColorRequestPayload](message)
		if err != nil {
			wsErr := player.WriteJsonError(message.Type, err)
			return true, wsErr
		}

		room := player.Room
		color := requestPayload.Color
		err = room.ChangePlayerColor(player.ID, color)
		if err != nil {
			wsErr := player.WriteJsonError(message.Type, err)
			return true, wsErr
		}

		room.EnqueueOutgoingMessage(BuildRoomMessage(room, fmt.Sprintf("%s.success", message.Type)), nil, nil)
		return true, nil
	case "room.toggle-ready":
		requestPayload, err := utils.ParseJsonPayload[roomPlayerReadyRequestPayload](message)
		if err != nil {
			wsErr := player.WriteJsonError(message.Type, err)
			return true, wsErr
		}

		room := player.Room
		ready := requestPayload.Ready
		err = room.TogglePlayerReadyState(player.ID, ready)
		if err != nil {
			wsErr := player.WriteJsonError(message.Type, err)
			return true, wsErr
		}

		room.EnqueueOutgoingMessage(BuildRoomMessage(room, fmt.Sprintf("%s.success", message.Type)), nil, nil)
		return true, nil
	case "room.start-game":
		room := player.Room
		err := StartMatch(player, room)
		if err != nil {
			wsErr := player.WriteJsonError(message.Type, err)
			return true, wsErr
		}

		room.EnqueueOutgoingMessage(buildStartMatch(room), nil, func() {
			room.StartRound()
			room.StartSubRound(round.SetupSettlement1)
			room.EnqueueBulkUpdate(
				match.UpdateCurrentRoundPlayerState,
				match.UpdateVertexState,
				match.UpdateEdgeState,
				match.UpdateLogs([]string{"Setup phase starting."}),
			)
		})
		return true, nil
	default:
		return false, nil
	}
}
