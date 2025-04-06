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
			Player: game.CurrentRoundPlayer().ID,
		},
	}
}

func UpdateVertexState(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-vertices", room.Status)
	availableVertices, err := game.AvailableVertices(username)
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: verticesStateUpdateResponsePayload{
			AvailableCityVertices:       game.SettlementsByPlayer(username),
			AvailableSettlementVertices: availableVertices,
			Enabled:                     err == nil,
			Highlight:                   room.Status == "setup",
		},
	}
}

func UpdateEdgeState(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-edges", room.Status)
	availableEdges, err := game.AvailableEdges(username)
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
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: passStateUpdateResponsePayload{
			Enabled: diceHasValue && isPlayerRound,
		},
	}
}

func UpdateTrade(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-trade", room.Status)
	dice := game.Dice()
	diceHasValue := dice[0] > 0 && dice[1] > 0
	isPlayerRound := game.CurrentRoundPlayer().ID == username
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: startTradeStateUpdateResponsePayload{
			Enabled: diceHasValue && isPlayerRound,
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

func UpdateLogs(logs []string) func(room *entities.Room, username string) *types.WebSocketServerResponse {
	return func(room *entities.Room, username string) *types.WebSocketServerResponse {
		messageType := fmt.Sprintf("%s.update-logs", room.Status)
		return &types.WebSocketServerResponse{
			Type:    types.ResponseType(messageType),
			Payload: logs,
		}
	}
}
