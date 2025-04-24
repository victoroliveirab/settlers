package prematch

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core"
	"github.com/victoroliveirab/settlers/core/packages/round"
	coreT "github.com/victoroliveirab/settlers/core/types"
	"github.com/victoroliveirab/settlers/logger"
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/handlers/match"
)

func metaEntriesToParams(entries []entities.RoomParamsMetaEntry) *core.Params {
	params := core.Params{}
	valueMap := map[string]*int{
		"speed":                &params.Speed,
		"bankTradeAmount":      &params.BankTradeAmount,
		"maxCards":             &params.MaxCards,
		"maxDevCardsPerRound":  &params.MaxDevCardsPerRound,
		"maxSettlements":       &params.MaxSettlements,
		"maxCities":            &params.MaxCities,
		"maxRoads":             &params.MaxRoads,
		"targetPoint":          &params.TargetPoint,
		"pointsPerSettlement":  &params.PointsPerSettlement,
		"pointsPerCity":        &params.PointsPerCity,
		"pointsForMostKnights": &params.PointsForMostKnights,
		"pointsForLongestRoad": &params.PointsForLongestRoad,
		"mostKnightsMinimum":   &params.MostKnightsMinimum,
		"longestRoadMinimum":   &params.LongestRoadMinimum,
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
			Color: *entry.Player.Color,
		}
	}

	params := metaEntriesToParams(room.Params())
	err := gameState.New(players, room.MapName, room.Rand, *params)
	if err != nil {
		return err
	}

	logger.LogSystemMessage("StartMatch", fmt.Sprintf("%s %s %v", room.ID, room.MapName, params))

	room.Game = gameState
	room.ProgressStatus()

	onSetupRoundTimeout := match.OnSetupRoundTimeoutCurry(room)
	onRegularRoundTimeout := match.OnRegularRoundTimeoutCurry(room)

	onBetweenTurnsTimeout := match.OnBetweenTurnsTimeoutCurry(room)
	onMoveRobberTimeout := match.OnMoveRobberTimeoutCurry(room)
	onPickRobbedTimeout := match.OnPickRobbedTimeoutCurry(room)
	onBuildRoadDevelopmentTimeout := match.OnBuildRoadDevelopmentTimeoutCurry(room)
	onMonopolyPickResourceTimeout := match.OnMonopolyPickResourceTimeoutCurry(room)
	onYearOfPlentyPickResourcesTimeout := match.OnYearOfPlentyPickResourcesTimeoutCurry(room)
	onDiscardPhaseTimeout := match.OnDiscardPhaseTimeoutCurry(room)

	room.CreateRoundManager(onRegularRoundTimeout, map[round.Type]func(){
		round.SetupSettlement1:          onSetupRoundTimeout,
		round.SetupRoad1:                onSetupRoundTimeout,
		round.SetupSettlement2:          onSetupRoundTimeout,
		round.SetupRoad2:                onSetupRoundTimeout,
		round.FirstRound:                onBetweenTurnsTimeout,
		round.MoveRobberDue7:            onMoveRobberTimeout,
		round.MoveRobberDueKnight:       onMoveRobberTimeout,
		round.PickRobbed:                onPickRobbedTimeout,
		round.BetweenTurns:              onBetweenTurnsTimeout,
		round.BuildRoad1Development:     onBuildRoadDevelopmentTimeout,
		round.BuildRoad2Development:     onBuildRoadDevelopmentTimeout,
		round.MonopolyPickResource:      onMonopolyPickResourceTimeout,
		round.YearOfPlentyPickResources: onYearOfPlentyPickResourcesTimeout,
		round.DiscardPhase:              onDiscardPhaseTimeout,
	})
	return nil
}
