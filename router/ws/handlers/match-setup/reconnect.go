package matchsetup

import (
	"fmt"

	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/handlers/match"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func ReconnectPlayer(room *entities.Room, playerID int64, conn *types.WebSocketConnection, onDisconnect func(player *entities.GamePlayer)) (*entities.GamePlayer, error) {
	player, err := room.ReconnectPlayer(playerID, conn, func(player *entities.GamePlayer) {
		room.RemovePlayer(playerID)
	})

	if err != nil {
		wsErr := utils.WriteJsonError(conn, playerID, "matchsetup.ReconnectPlayer", err)
		return nil, wsErr
	}

	game := room.Game
	if game == nil {
		err := fmt.Errorf("Game not assigned to the room", room.ID)
		return nil, err
	}

	mapState := match.UpdateMapState(room, player.Username)
	edgeState := match.UpdateEdgeState(room, player.Username)
	vertexState := match.UpdateVertexState(room, player.Username)
	currentRoundState := match.UpdateCurrentRoundPlayerState(room, player.Username)

	hydrateMsg := &types.WebSocketServerResponse{
		Type: "setup.hydrate",
		Payload: hydrateResponsePayload{
			EdgeUpdate:        edgeState,
			Map:               game.Map(),
			MapUpdate:         mapState,
			Players:           game.Players(),
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
