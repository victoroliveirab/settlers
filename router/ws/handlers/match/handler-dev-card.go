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
		room.EnqueueBulkUpdate(
			UpdateRobberMovement,
			UpdatePlayerDevHand,
			UpdatePlayerDevHandPermissions,
			UpdatePass,
			UpdateTrade,
			UpdateLogs([]string{fmt.Sprintf("%s used Knight card", player.Username)}),
		)
	} else if game.RoundType() == core.MonopolyPickResource {
		room.EnqueueBulkUpdate(
			UpdatePlayerDevHand,
			UpdatePlayerDevHandPermissions,
			UpdatePass,
			UpdateTrade,
			UpdateMonopoly,
			UpdateLogs([]string{fmt.Sprintf("%s used Monopoly card", player.Username)}),
		)
	} else if game.RoundType() == core.BuildRoad1Development || game.RoundType() == core.BuildRoad2Development {
		room.EnqueueBulkUpdate(
			UpdateEdgeState,
			UpdatePlayerDevHand,
			UpdatePlayerDevHandPermissions,
			UpdatePass,
			UpdateTrade,
			UpdateLogs([]string{fmt.Sprintf("%s used Road Building card", player.Username)}),
		)
	} else if game.RoundType() == core.YearOfPlentyPickResources {
		room.EnqueueBulkUpdate(
			UpdatePlayerDevHand,
			UpdatePlayerDevHandPermissions,
			UpdateYOP,
			UpdateLogs([]string{fmt.Sprintf("%s used Year of Plenty card", player.Username)}),
		)
	} else {
		room.EnqueueBulkUpdate(
			UpdatePass,
			UpdateTrade,
			UpdatePlayerDevHand,
			UpdatePlayerDevHandPermissions,
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

	resourceCountAfter := game.NumberOfResourcesByPlayer()[player.Username]
	resourceDiff := resourceCountAfter - resourceCountBefore

	logs := []string{fmt.Sprintf("%s stole %d [res]%s[/res] from the other players", player.Username, resourceDiff, payload.Resource)}

	room.EnqueueBulkUpdate(
		UpdateResourceCount,
		UpdatePlayerHand,
		UpdateMonopoly,
		UpdateLogs(logs),
	)

	return true, nil
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

	logs := []string{fmt.Sprintf("%s picked [res]%s[/res] and [res]%s[/res]", player.Username, payload.Resource1, payload.Resource2)}

	room.EnqueueBulkUpdate(
		UpdateResourceCount,
		UpdatePlayerHand,
		UpdateYOP,
		UpdateLogs(logs),
	)

	return true, nil
}
