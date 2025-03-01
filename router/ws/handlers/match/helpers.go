package match

import (
	"github.com/victoroliveirab/settlers/router/ws/entities"
)

func generateGameStateDump(player *entities.GamePlayer) map[string]interface{} {
	game := player.Room.Game
	currentRoundPlayer := game.CurrentRoundPlayer()
	return map[string]interface{}{
		"map":                game.Map(),
		"settlements":        game.AllSettlements(),
		"cities":             game.AllCities(),
		"roads":              game.AllRoads(),
		"players":            game.Players(),
		"currentRoundPlayer": currentRoundPlayer.ID,
		"roundType":          game.RoundType(),
		"hand":               game.ResourceHandByPlayer(player.Username),
		"resourceCount":      game.NumberOfResourcesByPlayer(),
		"dice":               game.Dice(),
	}
}
