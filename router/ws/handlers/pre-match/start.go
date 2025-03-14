package prematch

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core"
	coreT "github.com/victoroliveirab/settlers/core/types"
	"github.com/victoroliveirab/settlers/logger"
	"github.com/victoroliveirab/settlers/router/ws/entities"
	// matchsetup "github.com/victoroliveirab/settlers/router/ws/handlers/match-setup"
)

func StartMatch(player *entities.GamePlayer, room *entities.Room) error {
	if room.Owner != player.Username {
		err := fmt.Errorf("cannot start match on room %s: not room owner", room.ID)
		return err
	}

	gameState := &core.GameState{}
	players := make([]*coreT.Player, room.Capacity)
	for i, entry := range room.Participants {
		players[i] = &coreT.Player{
			ID:    entry.Player.Username,
			Color: entry.Player.Color,
		}
	}

	// REFACTOR: think in a better way of making this more straight forward (if you care)
	// May the lord forgive my soul for this code
	roomParams := room.Params()
	params := core.Params{}
	for _, param := range roomParams {
		switch param.Key {
		case "bankTradeAmount":
			params.BankTradeAmount = param.Value
			break
		case "maxCards":
			params.MaxCards = param.Value
			break
		case "maxSettlements":
			params.MaxSettlements = param.Value
			break
		case "maxCities":
			params.MaxCities = param.Value
			break
		case "maxRoads":
			params.MaxRoads = param.Value
			break
		case "maxDevCardsPerRound":
			params.MaxDevCardsPerRound = param.Value
			break
		case "targetPoint":
			params.TargetPoint = param.Value
			break
		case "pointsPerSettlement":
			params.PointsPerSettlement = param.Value
			break
		case "pointsPerCity":
			params.PointsPerCity = param.Value
			break
		case "pointsForMostKnights":
			params.PointsForMostKnights = param.Value
			break
		case "pointsForLongestRoad":
			params.PointsForLongestRoad = param.Value
			break
		case "longestRoadMinimum":
			params.LongestRoadMinimum = param.Value
			break
		case "mostKnightsMinimum":
			params.MostKnightsMinimum = param.Value
			break
		default:
			fmt.Println("unknown key:", param.Key)
		}
	}

	err := gameState.New(players, room.MapName, 42, params)
	if err != nil {
		return err
	}

	logger.LogSystemMessage("StartMatch", fmt.Sprintf("%s %s %v", room.ID, room.MapName, params))

	room.Game = gameState
	room.Status = "setup"

	// room.EnqueueBroadcastMessage(buildStartGameBroadcast(room, []string{"Setup phase starting."}), []int64{}, func() {
	// 	var firstPlayer *entities.GamePlayer
	// 	for _, participant := range room.Participants {
	// 		if participant.Player.Username == room.Owner {
	// 			firstPlayer = participant.Player
	// 			break
	// 		}
	// 	}
	// 	matchsetup.SendBuildSetupSettlementRequest(firstPlayer)
	// })
	return nil
}
