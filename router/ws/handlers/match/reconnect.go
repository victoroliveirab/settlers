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
				Map:               game.Map(),
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
	discardPhaseState := UpdateDiscardPhase(room, player.Username)
	passState := UpdatePass(room, player.Username)
	tradeState := UpdateTrade(room, player.Username)
	robberMovementState := UpdateRobberMovement(room, player.Username)

	hydrateMsg := &types.WebSocketServerResponse{
		Type: "match.hydrate",
		Payload: hydrateOngoingMatchResponsePayload{
			DiceUpdate:        diceState,
			DiscardUpdate:     discardPhaseState,
			EdgeUpdate:        edgeState,
			HandUpdate:        handState,
			Map:               game.Map(),
			MapUpdate:         mapState,
			PassUpdate:        passState,
			Players:           game.Players(),
			Ports:             game.Ports(),
			ResourceCount:     game.NumberOfResourcesByPlayer(),
			RobberUpdate:      robberMovementState,
			RoundPlayerUpdate: currentRoundState,
			TradeUpdate:       tradeState,
			VertexUpdate:      vertexState,
		},
	}

	wsErr := utils.WriteJson(player.Connection, player.ID, hydrateMsg)
	if wsErr != nil {
		return nil, wsErr
	}

	return player, nil
}
