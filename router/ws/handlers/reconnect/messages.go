package reconnect

import (
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func generateGameStateDump(player *entities.GamePlayer) map[string]interface{} {
	game := player.Room.Game
	currentRoundPlayer := game.CurrentRoundPlayer()
	// TODO: add ongoing trade offers
	return map[string]interface{}{
		"map":                game.Map(),
		"settlements":        game.AllSettlements(),
		"cities":             game.AllCities(),
		"roads":              game.AllRoads(),
		"players":            game.Players(),
		"currentRoundPlayer": currentRoundPlayer.ID,
		"hand":               game.ResourceHandByPlayer(player.Username),
		"resourceCount":      game.NumberOfResourcesByPlayer(),
		"devCount":           game.NumberOfDevCardsByPlayer(),
		"dice":               game.Dice(),
	}
}

func sendHydratePlayer(player *entities.GamePlayer) error {
	return utils.WriteJson(player.Connection, player.ID, &types.WebSocketMessage{
		Type: "hydrate",
		Payload: map[string]interface{}{
			"state": generateGameStateDump(player),
		},
	})
}
