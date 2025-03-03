package match

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/victoroliveirab/settlers/router/ws/entities"
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
		"dice":               game.Dice(),
	}
}

func diffResourceHands(before, after map[string]int) (map[string]int, error) {
	diff := map[string]int{}
	for resource, afterQuantity := range after {
		beforeQuantity, exists := before[resource]
		if !exists {
			return nil, fmt.Errorf("cannot compare two different shapes of hands - mismatch in key %s", resource)
		}
		diff[resource] = afterQuantity - beforeQuantity
	}

	return diff, nil
}

func hasDiff(diff map[string]int) bool {
	for _, quantity := range diff {
		if quantity != 0 {
			return true
		}
	}
	return false
}

func serializeHandDiff(hand map[string]int) string {
	var sb strings.Builder
	first := true
	for resource, quantity := range hand {
		if quantity != 0 {
			if !first {
				sb.WriteString(", ")
			}
			sb.WriteString(strconv.FormatInt(int64(quantity), 10))
			sb.WriteString(" ")
			sb.WriteString(resource)
		}
	}
	return sb.String()
}
