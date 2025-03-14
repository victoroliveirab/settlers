package prematch

import (
	"fmt"

	"github.com/victoroliveirab/settlers/router/ws/entities"
	matchsetup "github.com/victoroliveirab/settlers/router/ws/handlers/match-setup"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func TryHandle(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	switch message.Type {
	case "room.update-param":
		requestPayload, err := utils.ParseJsonPayload[roomUpdateParamRequestPayload](message)
		if err != nil {
			wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
			return true, wsErr
		}

		room := player.Room
		key := requestPayload.Key
		value := requestPayload.Value
		err = room.UpdateParam(player, key, value)
		if err != nil {
			wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
			return true, wsErr
		}

		room.EnqueueOutgoingMessage(BuildRoomMessage(room, fmt.Sprintf("%s.success", message.Type)), nil, nil)
		return true, nil
	case "room.player-change-color":
		requestPayload, err := utils.ParseJsonPayload[roomPlayerChangeColorRequestPayload](message)
		if err != nil {
			wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
			return true, wsErr
		}

		room := player.Room
		color := requestPayload.Color
		err = room.ChangePlayerColor(player.ID, color)
		if err != nil {
			wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
			return true, wsErr
		}

		room.EnqueueOutgoingMessage(BuildRoomMessage(room, fmt.Sprintf("%s.success", message.Type)), nil, nil)
		return true, nil
	case "room.toggle-ready":
		requestPayload, err := utils.ParseJsonPayload[roomPlayerReadyRequestPayload](message)
		if err != nil {
			wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
			return true, wsErr
		}

		room := player.Room
		ready := requestPayload.Ready
		err = room.TogglePlayerReadyState(player.ID, ready)
		if err != nil {
			wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
			return true, wsErr
		}

		room.EnqueueOutgoingMessage(BuildRoomMessage(room, fmt.Sprintf("%s.success", message.Type)), nil, nil)
		return true, nil
	case "room.start-game":
		room := player.Room
		err := StartMatch(player, room)
		if err != nil {
			wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
			return true, wsErr
		}

		room.EnqueueOutgoingMessage(buildStartMatch(room), nil, func() {
			room.EnqueueBulkUpdate(
				matchsetup.UpdateCurrentRoundPlayerState,
				matchsetup.UpdateVertexState,
				matchsetup.UpdateEdgeState,
				matchsetup.UpdateLogs([]string{"Setup phase starting."}),
			)
		})
		return true, nil
	default:
		return false, nil
	}
}
