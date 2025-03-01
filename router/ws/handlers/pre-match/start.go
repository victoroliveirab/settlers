package prematch

import (
	"github.com/victoroliveirab/settlers/core"
	coreT "github.com/victoroliveirab/settlers/core/types"
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/handlers/match"
)

func StartMatch(room *entities.Room) error {
	gameState := &core.GameState{}
	players := make([]*coreT.Player, room.Capacity)
	for i, entry := range room.Participants {
		players[i] = &coreT.Player{
			ID:    entry.Player.Username,
			Color: entry.Player.Color,
		}
	}

	// FIXME: params should come from room
	err := gameState.New(players, room.MapName, 42, core.Params{
		BankTradeAmount:      4,
		MaxCards:             7,
		MaxDevCardsPerRound:  1,
		MaxSettlements:       5,
		MaxCities:            4,
		MaxRoads:             20,
		TargetPoint:          10,
		PointsPerSettlement:  1,
		PointsPerCity:        2,
		PointsForMostKnights: 2,
		PointsForLongestRoad: 2,
		MostKnightsMinimum:   3,
		LongestRoadMinimum:   5,
	})
	if err != nil {
		return err
	}

	room.Game = gameState

	room.EnqueueBroadcastMessage(buildStartGameBroadcast(room, []string{"Setup phase starting."}), []int64{}, func() {
		var firstPlayer *entities.GamePlayer
		for _, participant := range room.Participants {
			if participant.Player.Username == room.Owner {
				firstPlayer = participant.Player
				break
			}
		}
		match.SendBuildSetupSettlementRequest(firstPlayer)
	})
	return nil
}
