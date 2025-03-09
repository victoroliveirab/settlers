package match

import (
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func SendPlayerRound(room *entities.Room, player *entities.GamePlayer) error {
	availableVertices, _ := room.Game.AvailableVertices(player.Username)
	availableEdges, _ := room.Game.AvailableEdges(player.Username)
	return utils.WriteJson(player.Connection, player.ID, &types.WebSocketMessage{
		Type: "game.your-round",
		Payload: map[string]interface{}{
			"availableEdges":     availableEdges,
			"availableVertices":  availableVertices,
			"currentRoundPlayer": room.Game.CurrentRoundPlayer().ID,
			"roundType":          room.Game.RoundType(),
			"round":              room.Game.Round(),
		},
	})
}

func sendHydratePlayer(player *entities.GamePlayer) error {
	game := player.Room.Game
	currentRoundPlayer := game.CurrentRoundPlayer()
	availableVertices, _ := game.AvailableVertices(player.Username)
	availableEdges, _ := game.AvailableEdges(player.Username)
	return utils.WriteJson(player.Connection, player.ID, &types.WebSocketMessage{
		Type: "game.hydrate",
		Payload: map[string]interface{}{
			"state": map[string]interface{}{
				"availableEdges":     availableEdges,
				"availableVertices":  availableVertices,
				"map":                game.Map(),
				"settlements":        game.AllSettlements(),
				"cities":             game.AllCities(),
				"roads":              game.AllRoads(),
				"players":            game.Players(),
				"roundType":          game.RoundType(),
				"round":              game.Round(),
				"currentRoundPlayer": currentRoundPlayer.ID,
				"hand":               game.ResourceHandByPlayer(player.Username),
				"resourceCount":      game.NumberOfResourcesByPlayer(),
				"devCount":           game.NumberOfDevCardsByPlayer(),
				"dice":               game.Dice(),
			},
		},
	})
}
