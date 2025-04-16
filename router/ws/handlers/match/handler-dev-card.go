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

type monopolyPickRequestPayload struct {
	Resource string `json:"resource"`
}

type yearOfPlentyPickRequestPayload struct {
	Resource1 string `json:"resource1"`
	Resource2 string `json:"resource2"`
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
		UpdateCurrentRoundPlayerState,
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
		room.StartSubRound(core.MoveRobberDueKnight)
		room.EnqueueBulkUpdate(
			UpdateCurrentRoundPlayerState,
			UpdateRobberMovement,
			UpdatePlayerDevHand,
			UpdatePlayerDevHandPermissions,
			UpdatePass,
			UpdateTrade,
			UpdatePoints,
			UpdateLogs([]string{fmt.Sprintf("%s used Knight card", player.Username)}),
		)
	} else if game.RoundType() == core.MonopolyPickResource {
		room.StartSubRound(core.MonopolyPickResource)
		room.EnqueueBulkUpdate(
			UpdateCurrentRoundPlayerState,
			UpdatePlayerDevHand,
			UpdatePlayerDevHandPermissions,
			UpdatePass,
			UpdateTrade,
			UpdateMonopoly,
			UpdatePoints,
			UpdateLogs([]string{fmt.Sprintf("%s used Monopoly card", player.Username)}),
		)
	} else if game.RoundType() == core.BuildRoad1Development || game.RoundType() == core.BuildRoad2Development {
		room.StartSubRound(game.RoundType())
		room.EnqueueBulkUpdate(
			UpdateCurrentRoundPlayerState,
			UpdateEdgeState,
			UpdatePlayerDevHand,
			UpdatePlayerDevHandPermissions,
			UpdatePass,
			UpdateTrade,
			UpdatePoints,
			UpdateLogs([]string{fmt.Sprintf("%s used Road Building card", player.Username)}),
		)
	} else if game.RoundType() == core.YearOfPlentyPickResources {
		room.StartSubRound(core.YearOfPlentyPickResources)
		room.EnqueueBulkUpdate(
			UpdateCurrentRoundPlayerState,
			UpdatePlayerDevHand,
			UpdatePlayerDevHandPermissions,
			UpdateYOP,
			UpdatePoints,
			UpdateLogs([]string{fmt.Sprintf("%s used Year of Plenty card", player.Username)}),
		)
	} else {
		room.EnqueueBulkUpdate(
			UpdateCurrentRoundPlayerState,
			UpdatePass,
			UpdateTrade,
			UpdatePlayerDevHand,
			UpdatePlayerDevHandPermissions,
			UpdatePoints,
		)
	}

	return true, nil
}

func handleMonopolyResource(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	payload, err := utils.ParseJsonPayload[monopolyPickRequestPayload](message)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	room := player.Room
	game := room.Game
	resourceCountBefore := game.NumberOfResourcesByPlayer()[player.Username]

	err = game.PickMonopolyResource(player.Username, payload.Resource)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	handleMonopolyResourceResponse(room, payload.Resource, resourceCountBefore)
	return true, nil
}

func handleMonopolyResourceResponse(room *entities.Room, resourceStolen string, resourceCountBefore int) {
	game := room.Game
	currentRoundPlayer := game.CurrentRoundPlayer().ID

	resourceCountAfter := game.NumberOfResourcesByPlayer()[currentRoundPlayer]
	resourceDiff := resourceCountAfter - resourceCountBefore

	logs := []string{fmt.Sprintf("%s stole %d [res]%s[/res] from the other players", currentRoundPlayer, resourceDiff, resourceStolen)}
	room.ResumeRound()
	room.EnqueueBulkUpdate(
		UpdateCurrentRoundPlayerState,
		UpdateResourceCount,
		UpdatePlayerHand,
		UpdateMonopoly,
		UpdateLogs(logs),
	)
}

func handlePickYearOfPlentyResources(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	payload, err := utils.ParseJsonPayload[yearOfPlentyPickRequestPayload](message)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	room := player.Room
	game := room.Game
	err = game.PickYearOfPlentyResources(player.Username, payload.Resource1, payload.Resource2)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	handlePickYearOfPlentyResourcesResponse(room, payload.Resource1, payload.Resource2)
	return true, nil
}

func handlePickYearOfPlentyResourcesResponse(room *entities.Room, resource1, resource2 string) {
	game := room.Game
	currentRoundPlayer := game.CurrentRoundPlayer().ID

	logs := []string{fmt.Sprintf("%s picked [res]%s[/res] and [res]%s[/res]", currentRoundPlayer, resource1, resource2)}
	room.ResumeRound()
	room.EnqueueBulkUpdate(
		UpdateCurrentRoundPlayerState,
		UpdateResourceCount,
		UpdatePlayerHand,
		UpdateYOP,
		UpdateLogs(logs),
	)
}
