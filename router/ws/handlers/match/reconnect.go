package match

import (
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func SendCurrentGameState(player *entities.GamePlayer) error {
	room := player.Room
	game := room.Game

	if room.Status == "setup" {
		mapState := UpdateMapState(room, player.Username)
		edgeState := UpdateEdgeState(room, player.Username)
		vertexState := UpdateVertexState(room, player.Username)
		currentRoundState := UpdateCurrentRoundPlayerState(room, player.Username)

		hydrateMsg := &types.WebSocketServerResponse{
			Type: "setup.hydrate",
			Payload: hydrateSetupMatchResponsePayload{
				DevHandCount:      game.NumberOfDevCardsByPlayer(),
				EdgeUpdate:        edgeState,
				Map:               game.GetBoard(),
				MapName:           game.MapName(),
				MapUpdate:         mapState,
				Players:           game.Players(),
				Ports:             game.Ports(),
				ResourceCount:     game.NumberOfResourcesByPlayer(),
				RoomStatus:        room.Status,
				RoundPlayerUpdate: currentRoundState,
				VertexUpdate:      vertexState,
			},
		}

		wsErr := player.WriteJSON(hydrateMsg)
		return wsErr
	}

	mapState := UpdateMapState(room, player.Username)
	edgeState := UpdateEdgeState(room, player.Username)
	vertexState := UpdateVertexState(room, player.Username)
	currentRoundState := UpdateCurrentRoundPlayerState(room, player.Username)
	diceState := UpdateDiceState(room, player.Username)
	handState := UpdatePlayerHand(room, player.Username)
	pointsState := UpdatePoints(room, player.Username)
	portsState := UpdatePortsState(room, player.Username)
	knightsUsageState := UpdateKnightUsage(room, player.Username)
	longestRoadState := UpdateLongestRoadSize(room, player.Username)
	devHandState := UpdatePlayerDevHand(room, player.Username)
	devHandPermissionsState := UpdatePlayerDevHandPermissions(room, player.Username)
	discardPhaseState := UpdateDiscardPhase(room, player.Username)
	passState := UpdatePass(room, player.Username)
	tradeState := UpdateTrade(room, player.Username)
	tradeOffersState := UpdateTradeOffers(room, player.Username)
	robberMovementState := UpdateRobberMovement(room, player.Username)
	robbablePlayersState := UpdateRobbablePlayers(room, player.Username)
	buyDevCardState := UpdateBuyDevelopmentCard(room, player.Username)
	yearOfPlentyState := UpdateYOP(room, player.Username)

	hydrateMsg := &types.WebSocketServerResponse{
		Type: "match.hydrate",
		Payload: hydrateOngoingMatchResponsePayload{
			BuyDevCardUpdate:         buyDevCardState,
			DevHandCount:             game.NumberOfDevCardsByPlayer(),
			DevHandUpdate:            devHandState,
			DevHandPermissionsUpdate: devHandPermissionsState,
			DiceUpdate:               diceState,
			DiscardUpdate:            discardPhaseState,
			EdgeUpdate:               edgeState,
			HandUpdate:               handState,
			KnightsUsageUpdate:       knightsUsageState,
			LongestRoadUpdate:        longestRoadState,
			Map:                      game.GetBoard(),
			MapName:                  game.MapName(),
			MapUpdate:                mapState,
			PassActionState:          passState,
			Players:                  game.Players(),
			PointsUpdate:             pointsState,
			Ports:                    game.Ports(),
			PortsUpdate:              portsState,
			ResourceCount:            game.NumberOfResourcesByPlayer(),
			RobbablePlayersUpdate:    robbablePlayersState,
			RobberUpdate:             robberMovementState,
			RoomStatus:               room.Status,
			RoundPlayerUpdate:        currentRoundState,
			TradeActionState:         tradeState,
			TradeOffersUpdate:        tradeOffersState,
			VertexUpdate:             vertexState,
			YearOfPlentyUpdate:       yearOfPlentyState,
		},
	}

	wsErr := player.WriteJSON(hydrateMsg)
	return wsErr
}
