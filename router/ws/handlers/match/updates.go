package match

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core"
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func UpdateCurrentRoundPlayerState(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-round-player", room.Status)
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: currentRoundPlayerStateUpdateResponsePayload{
			Deadline:    room.RoundDeadline(),
			Player:      game.CurrentRoundPlayer().ID,
			ServerTime:  room.Now(),
			SubDeadline: room.SubRoundDeadline(),
		},
	}
}

func UpdateVertexState(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-vertices", room.Status)
	availableSettlementVertices, err := game.AvailableVertices(username)
	var availableCityVertices []int
	if err == nil {
		availableCityVertices = game.SettlementsByPlayer(username)
	} else {
		availableCityVertices = []int{}
	}
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: verticesStateUpdateResponsePayload{
			AvailableCityVertices:       availableCityVertices,
			AvailableSettlementVertices: availableSettlementVertices,
			Enabled:                     err == nil,
			Highlight:                   room.Status == "setup",
		},
	}
}

func UpdateEdgeState(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-edges", room.Status)
	availableEdges, err := game.AvailableEdges(username)
	// TODO: perhaps highlight during road building phases as well
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: edgesStateUpdateResponsePayload{
			AvailableEdges: availableEdges,
			Enabled:        err == nil,
			Highlight:      room.Status == "setup",
		},
	}
}

func UpdateDiceState(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-dice", room.Status)
	dice := game.Dice()
	diceHasValue := dice[0] > 0 && dice[1] > 0
	isPlayerRound := game.CurrentRoundPlayer().ID == username
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: diceStateUpdateResponsePayload{
			Dice:    game.Dice(),
			Enabled: !diceHasValue && isPlayerRound,
		},
	}
}

func UpdateMapState(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-map", room.Status)
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: mapStateUpdateResponsePayload{
			BlockedTiles: game.BlockedTiles(),
			Cities:       game.AllCities(),
			Roads:        game.AllRoads(),
			Settlements:  game.AllSettlements(),
		},
	}
}

func UpdatePortsState(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-ports", room.Status)
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: portStateUpdateResponsePayload{
			Ports: game.PortsByPlayer(username),
		},
	}
}

func UpdatePlayerHand(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-hand", room.Status)
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: handStateUpdateResponsePayload{
			Hand: game.ResourceHandByPlayer(username),
		},
	}
}

func UpdatePlayerDevHand(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-dev-hand", room.Status)
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: devHandStateUpdateResponsePayload{
			DevHand: game.DevelopmentHandByPlayer(username),
		},
	}
}

func UpdatePlayerDevHandPermissions(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-dev-hand-permissions", room.Status)
	devHand := game.DevelopmentHandByPlayer(username)
	permissions := make(map[string]bool)
	for devHandKind, quantity := range devHand {
		if quantity == 0 {
			permissions[devHandKind] = false
		} else {
			permissions[devHandKind] = game.IsDevCardPlayable(username, devHandKind) == nil
		}
	}
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: devHandPermissionsStateUpdateResponsePayload{
			DevHandPermissions: permissions,
		},
	}
}

func UpdateResourceCount(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-resource-count", room.Status)
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: resourceCountStateUpdateResponsePayload{
			ResourceCount: game.NumberOfResourcesByPlayer(),
		},
	}
}

func UpdateRobberMovement(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-robber-movement", room.Status)
	enabled := game.RoundType() == core.MoveRobberDue7 || game.RoundType() == core.MoveRobberDueKnight
	enabled = enabled && game.CurrentRoundPlayer().ID == username
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: moveRobberStateUpdateResponsePayload{
			AvailableTiles: game.UnblockedTiles(),
			Enabled:        enabled,
			Highlight:      enabled,
		},
	}
}

func UpdateDiscardPhase(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-discard-phase", room.Status)
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: discardPhaseStateUpdateResponsePayload{
			DiscardAmounts: game.DiscardAmounts(),
			Enabled:        game.DiscardAmountByPlayer(username) > 0,
		},
	}
}

func UpdatePass(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-pass", room.Status)
	dice := game.Dice()
	diceHasValue := dice[0] > 0 && dice[1] > 0
	isPlayerRound := game.CurrentRoundPlayer().ID == username
	isRoundStateToEnable := game.RoundType() == core.Regular
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: passStateUpdateResponsePayload{
			Enabled: diceHasValue && isPlayerRound && isRoundStateToEnable,
		},
	}
}

func UpdateTrade(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-trade", room.Status)
	dice := game.Dice()
	diceHasValue := dice[0] > 0 && dice[1] > 0
	isPlayerRound := game.CurrentRoundPlayer().ID == username
	isRoundStateToEnable := game.RoundType() == core.Regular
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: startTradeStateUpdateResponsePayload{
			Enabled: diceHasValue && isPlayerRound && isRoundStateToEnable,
		},
	}
}

func UpdateBuyDevelopmentCard(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-buy-dev-card", room.Status)
	enabled := game.IsBuyDevelopmentCardAvailable(username) == nil
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: buyDevCardStateUpdateResponsePayload{
			Enabled: enabled,
		},
	}
}

func UpdateRobbablePlayers(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-pick-robbed", room.Status)
	enabled := game.RoundType() == core.PickRobbed && game.CurrentRoundPlayer().ID == username
	robbablePlayers, _ := game.RobbablePlayers(username)
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: pickRobbedStateUpdate{
			Enabled: enabled,
			Options: robbablePlayers,
		},
	}
}

func UpdateTradeOffers(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-trade-offers", room.Status)
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: updateActiveTradeOffersStateUpdate{
			Offers: game.ActiveTradeOffers(),
		},
	}
}

func UpdatePoints(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-points", room.Status)
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: pointsStateUpdate{
			Points: game.Points(),
		},
	}
}

func UpdateLongestRoadSize(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-longest-road-size", room.Status)
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: longestRoadStateUpdate{
			LongestRoadSizeByPlayer: game.LongestRoadLengths(),
		},
	}
}

func UpdateKnightUsage(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-knight-usage", room.Status)
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: knightUsageStateUpdate{
			KnightUsesByPlayer: game.KnightUses(),
		},
	}
}

func UpdateMonopoly(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-monopoly", room.Status)
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: monopolyStateUpdate{
			Enabled: game.CurrentRoundPlayer().ID == username && game.RoundType() == core.MonopolyPickResource,
		},
	}
}

func UpdateYOP(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-year-of-plenty", room.Status)
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: yearOfPlentyStateUpdate{
			Enabled: game.CurrentRoundPlayer().ID == username && game.RoundType() == core.YearOfPlentyPickResources,
		},
	}
}

func UpdateLogs(logs []string) func(room *entities.Room, username string) *types.WebSocketServerResponse {
	return func(room *entities.Room, username string) *types.WebSocketServerResponse {
		messageType := fmt.Sprintf("%s.update-logs", room.Status)
		return &types.WebSocketServerResponse{
			Type:    types.ResponseType(messageType),
			Payload: logs,
		}
	}
}
