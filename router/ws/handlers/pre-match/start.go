package prematch

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core"
	coreT "github.com/victoroliveirab/settlers/core/types"
	"github.com/victoroliveirab/settlers/logger"
	"github.com/victoroliveirab/settlers/router/ws/entities"
	// matchsetup "github.com/victoroliveirab/settlers/router/ws/handlers/match-setup"
)

func metaEntriesToParams(entries []entities.RoomParamsMetaEntry) *core.Params {
	params := core.Params{}
	valueMap := map[string]*int{
		"BankTradeAmount":      &params.BankTradeAmount,
		"MaxCards":             &params.MaxCards,
		"MaxDevCardsPerRound":  &params.MaxDevCardsPerRound,
		"MaxSettlements":       &params.MaxSettlements,
		"MaxCities":            &params.MaxCities,
		"MaxRoads":             &params.MaxRoads,
		"TargetPoint":          &params.TargetPoint,
		"PointsPerSettlement":  &params.PointsPerSettlement,
		"PointsPerCity":        &params.PointsPerCity,
		"PointsForMostKnights": &params.PointsForMostKnights,
		"PointsForLongestRoad": &params.PointsForLongestRoad,
		"MostKnightsMinimum":   &params.MostKnightsMinimum,
		"LongestRoadMinimum":   &params.LongestRoadMinimum,
	}

	for _, entry := range entries {
		if ptr, ok := valueMap[entry.Key]; ok {
			*ptr = entry.Value
		}
	}

	return &params
}

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

	params := metaEntriesToParams(room.Params())
	err := gameState.New(players, room.MapName, 42, *params)
	if err != nil {
		return err
	}

	logger.LogSystemMessage("StartMatch", fmt.Sprintf("%s %s %v", room.ID, room.MapName, params))

	room.Game = gameState
	room.Status = "setup"
	return nil
}
