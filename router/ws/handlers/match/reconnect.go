package match

import (
	"fmt"

	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func ReconnectPlayer(room *entities.Room, playerID int64, conn *types.WebSocketConnection, onDisconnect func(player *entities.GamePlayer)) (*entities.GamePlayer, error) {
	player, err := room.ReconnectPlayer(playerID, conn, func(player *entities.GamePlayer) {
		room.RemovePlayer(playerID)
	})
	if err != nil {
		wsErr := utils.WriteJsonError(conn, playerID, "match.ReconnectPlayer", err)
		return nil, wsErr
	}

	game := room.Game
	if game == nil {
		err := fmt.Errorf("Game not assigned to the room", room.ID)
		return nil, err
	}

	if room.Status == "setup" {
		mapState := UpdateMapState(room, player.Username)
		edgeState := UpdateEdgeState(room, player.Username)
		vertexState := UpdateVertexState(room, player.Username)
		currentRoundState := UpdateCurrentRoundPlayerState(room, player.Username)

		hydrateMsg := &types.WebSocketServerResponse{
			Type: "setup.hydrate",
			Payload: hydrateSetupMatchResponsePayload{
				EdgeUpdate:        edgeState,
				Map:               game.GetBoard(),
				MapName:           game.MapName(),
				MapUpdate:         mapState,
				Players:           game.Players(),
				Ports:             game.Ports(),
				ResourceCount:     game.NumberOfResourcesByPlayer(),
				RoundPlayerUpdate: currentRoundState,
				VertexUpdate:      vertexState,
			},
		}

		wsErr := utils.WriteJson(player.Connection, player.ID, hydrateMsg)
		if wsErr != nil {
			return nil, wsErr
		}

		return player, nil
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
			RoundPlayerUpdate:        currentRoundState,
			TradeActionState:         tradeState,
			TradeOffersUpdate:        tradeOffersState,
			VertexUpdate:             vertexState,
			YearOfPlentyUpdate:       yearOfPlentyState,
		},
	}

	wsErr := utils.WriteJson(player.Connection, player.ID, hydrateMsg)
	if wsErr != nil {
		return nil, wsErr
	}

	return player, nil
}
