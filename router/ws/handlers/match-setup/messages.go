package matchsetup

import (
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func sendHydratePlayer(player *entities.GamePlayer) error {
	game := player.Room.Game
	currentRoundPlayer := game.CurrentRoundPlayer()
	return utils.WriteJson(player.Connection, player.ID, &types.WebSocketServerResponse{
		Type: "setup.hydrate",
		Payload: map[string]interface{}{
			"state": map[string]interface{}{
				"map":                game.Map(),
				"settlements":        game.AllSettlements(),
				"cities":             game.AllCities(),
				"roads":              game.AllRoads(),
				"players":            game.Players(),
				"currentRoundPlayer": currentRoundPlayer.ID,
			},
		},
	})
}
