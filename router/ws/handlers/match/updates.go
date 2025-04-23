package match

import (
	"fmt"

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
	isRoadBuilding := game.IsRoadBuilding()
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: edgesStateUpdateResponsePayload{
			AvailableEdges: availableEdges,
			Enabled:        err == nil,
			Highlight:      isRoadBuilding,
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
			Cities:       game.GetAllCities(),
			Roads:        game.GetAllRoads(),
			Settlements:  game.GetAllSettlements(),
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
			permissions[devHandKind] = game.IsDevCardPlayable(username, devHandKind)
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
	isRobberTurn := game.IsRobberTurn()
	isPlayerTurn := game.IsPlayerTurn(username)
	enabled := isRobberTurn && isPlayerTurn
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
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: passStateUpdateResponsePayload{
			Enabled: game.IsPassTurnAllowed(username),
		},
	}
}

func UpdateTrade(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-trade", room.Status)
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: startTradeStateUpdateResponsePayload{
			Enabled: game.IsStartTradeAllowed(username),
		},
	}
}

func UpdateBuyDevelopmentCard(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-buy-dev-card", room.Status)
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: buyDevCardStateUpdateResponsePayload{
			Enabled: game.IsBuyDevCardAllowed(username),
		},
	}
}

func UpdateRobbablePlayers(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-pick-robbed", room.Status)
	enabled := game.IsPickingRobbedAllowed(username)
	var robbablePlayers []string
	if enabled {
		robbablePlayers, _ = game.RobbablePlayers(username)
	}
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
			Points: game.PublicPoints(),
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
			Enabled: game.IsPickingMonopolyAllowed(username),
		},
	}
}

func UpdateYOP(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-year-of-plenty", room.Status)
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: yearOfPlentyStateUpdate{
			Enabled: game.IsPickingYOPAllowed(username),
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
