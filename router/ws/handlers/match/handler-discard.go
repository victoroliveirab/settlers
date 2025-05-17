package match

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core/packages/round"
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
		wsErr := player.WriteJsonError(message.Type, err)
		return true, wsErr
	}

	room := player.Room
	game := room.Game
	err = game.DiscardPlayerCards(player.Username, payload.Resources)
	if err != nil {
		wsErr := player.WriteJsonError(message.Type, err)
		return true, wsErr
	}

	formattedResources := utils.FormatResources(payload.Resources)
	logs := []string{fmt.Sprintf("%s discarded %s", player.Username, formattedResources)}
	handleDiscardCardsResponse(room, logs)
	return true, nil
}

func handleDiscardCardsResponse(room *entities.Room, logs []string) {
	game := room.Game

	if game.RoundType() == round.MoveRobberDue7 {
		// last player that needed to discard has discarded
		room.StartSubRound(round.MoveRobberDue7)
		room.EnqueueBulkUpdate(
			UpdateCurrentRoundPlayerState,
			UpdatePlayerHand,
			UpdateResourceCount,
			UpdateDiscardPhase,
			UpdateRobberMovement,
			UpdateLogs(logs),
		)
	} else {
		// there are still players that need to discard
		room.EnqueueBulkUpdate(
			UpdateCurrentRoundPlayerState,
			UpdatePlayerHand,
			UpdateResourceCount,
			UpdateDiscardPhase,
			UpdateLogs(logs),
		)
	}
}
