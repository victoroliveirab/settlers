package match

import (
	"fmt"
	"maps"
	"math/rand"

	"github.com/victoroliveirab/settlers/core/packages/round"
	"github.com/victoroliveirab/settlers/logger"
	"github.com/victoroliveirab/settlers/router/ws/entities"
	wsUtils "github.com/victoroliveirab/settlers/router/ws/utils"
	"github.com/victoroliveirab/settlers/utils"
)

func OnSetupRoundTimeoutCurry(room *entities.Room) func() {
	return func() {
		game := room.Game
		currentRoundPlayer := game.CurrentRoundPlayer().ID
		currentRoundType := game.RoundType()

		logger.LogSystemMessage(fmt.Sprintf("onSetupRoundTimeout.%s", room.ID), fmt.Sprintf("handling timeout on %s for player %s", round.RoundTypeTranslation[currentRoundType], currentRoundPlayer))

		if currentRoundType == round.SetupSettlement1 || currentRoundType == round.SetupSettlement2 {
			availableSettlements, _ := game.AvailableVertices(currentRoundPlayer)
			vertexID := utils.SliceGetRandom(availableSettlements, room.Rand)
			game.BuildSettlement(currentRoundPlayer, vertexID)
			room.StartSubRound(game.RoundType())
			room.EnqueueBulkUpdate(
				UpdateCurrentRoundPlayerState,
				UpdateMapState,
				UpdateEdgeState,
				UpdateVertexState,
				UpdatePlayerHand,
				UpdatePortsState,
				UpdateBuyDevelopmentCard,
				UpdatePoints,
				UpdatePortsState,
				UpdateLongestRoadSize,
				UpdateLogs([]string{fmt.Sprintf("%s built a new settlement.", currentRoundPlayer)}),
			)
		} else {
			availableRoads, _ := game.AvailableEdges(currentRoundPlayer)
			edgeID := utils.SliceGetRandom(availableRoads, room.Rand)
			game.BuildRoad(currentRoundPlayer, edgeID)
			handleEdgeClickSetupResponse(room, []string{fmt.Sprintf("%s built a new road.", currentRoundPlayer)})
		}
	}
}

func OnRegularRoundTimeoutCurry(room *entities.Room) func() {
	return func() {
		game := room.Game
		currentRoundPlayer := game.CurrentRoundPlayer().ID
		logger.LogSystemMessage(fmt.Sprintf("onRegularRoundTimeout.%s", room.ID), fmt.Sprintf("handling timeout for player %s", currentRoundPlayer))

		game.EndRound(currentRoundPlayer)
		handleEndRoundResponse(room, currentRoundPlayer)
	}
}

func OnMoveRobberTimeoutCurry(room *entities.Room) func() {
	return func() {
		game := room.Game
		currentRoundPlayer := game.CurrentRoundPlayer().ID
		logger.LogSystemMessage(fmt.Sprintf("onMoveRobberTimeout.%s", room.ID), fmt.Sprintf("handling timeout for player %s", currentRoundPlayer))

		tileID := utils.SliceGetRandom(game.UnblockedTiles(), room.Rand)
		game.MoveRobber(currentRoundPlayer, tileID)
		handleMoveRobberResponse(room)
	}
}

func OnPickRobbedTimeoutCurry(room *entities.Room) func() {
	return func() {
		game := room.Game
		currentRoundPlayer := game.CurrentRoundPlayer().ID
		logger.LogSystemMessage(fmt.Sprintf("onPickRobbedTimeout.%s", room.ID), fmt.Sprintf("handling timeout for player %s", currentRoundPlayer))
		playersToRob, _ := game.RobbablePlayers(currentRoundPlayer)
		robbedPlayer := utils.SliceGetRandom(playersToRob, room.Rand)
		game.RobPlayer(currentRoundPlayer, robbedPlayer)
		handlePickRobbedResponse(room, currentRoundPlayer, robbedPlayer)
	}
}

func OnBetweenTurnsTimeoutCurry(room *entities.Room) func() {
	return func() {
		game := room.Game
		currentRoundPlayer := game.CurrentRoundPlayer().ID
		logger.LogSystemMessage(fmt.Sprintf("onBetweenTurnsTimeout.%s", room.ID), fmt.Sprintf("handling timeout for player %s", currentRoundPlayer))
		game.RollDice(currentRoundPlayer)
		prevResourceHands := map[string]map[string]int{}
		for _, player := range game.Players() {
			prevResourceHands[player.ID] = maps.Clone(game.ResourceHandByPlayer(player.ID))
		}

		handleDiceRollResponse(room, prevResourceHands)
	}
}

func OnBuildRoadDevelopmentTimeoutCurry(room *entities.Room) func() {
	return func() {
		game := room.Game
		currentRoundPlayer := game.CurrentRoundPlayer().ID
		var numberOfRoadsToBuild int
		if game.RoundType() == round.BuildRoad1Development {
			numberOfRoadsToBuild = 2
		} else {
			numberOfRoadsToBuild = 1
		}

		logs := make([]string, 0)
		for i := 0; i < numberOfRoadsToBuild; i++ {
			availableRoads, err := game.AvailableEdges(currentRoundPlayer)
			if err != nil {
				// May have reached max number of roads
				break
			}
			edgeID := utils.SliceGetRandom(availableRoads, room.Rand)
			game.BuildRoad(currentRoundPlayer, edgeID)
			logs = append(logs, fmt.Sprintf("%s has built a new road.", currentRoundPlayer))
		}

		// Force it back to regular round count
		handleEdgeClickMatchResponse(room, round.BuildRoad2Development, logs)
	}
}

func OnMonopolyPickResourceTimeoutCurry(room *entities.Room) func() {
	return func() {
		game := room.Game
		currentRoundPlayer := game.CurrentRoundPlayer().ID
		resourceCountBefore := game.NumberOfResourcesByPlayer()[currentRoundPlayer]
		resource := utils.SliceGetRandom([]string{"Lumber", "Brick", "Sheep", "Grain", "Ore"}, room.Rand)

		game.PickMonopolyResource(currentRoundPlayer, resource)
		handleMonopolyResourceResponse(room, resource, resourceCountBefore)
	}
}

func OnYearOfPlentyPickResourcesTimeoutCurry(room *entities.Room) func() {
	return func() {
		game := room.Game
		currentRoundPlayer := game.CurrentRoundPlayer().ID
		resource1 := utils.SliceGetRandom([]string{"Lumber", "Brick", "Sheep", "Grain", "Ore"}, room.Rand)
		resource2 := utils.SliceGetRandom([]string{"Lumber", "Brick", "Sheep", "Grain", "Ore"}, room.Rand)
		game.PickYearOfPlentyResources(currentRoundPlayer, resource1, resource2)
		handlePickYearOfPlentyResourcesResponse(room, resource1, resource2)
	}
}

func distributeResources(input map[string]int, quantity int, randGenerator *rand.Rand) map[string]int {
	result := make(map[string]int)
	remaining := quantity

	// Copy keys into a slice for randomized iteration
	keys := make([]string, 0, len(input))
	for k := range input {
		keys = append(keys, k)
	}

	// Keep distributing 1 unit at a time randomly
	for remaining > 0 {
		// Shuffle keys each round to avoid bias
		randGenerator.Shuffle(len(keys), func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })

		for _, k := range keys {
			if remaining == 0 {
				break
			}
			if result[k] < input[k] {
				result[k]++
				remaining--
			}
		}
	}

	return result
}

func OnDiscardPhaseTimeoutCurry(room *entities.Room) func() {
	return func() {
		game := room.Game
		discardAmounts := game.DiscardAmounts()

		logs := make([]string, 0)
		for player, amount := range discardAmounts {
			if amount == 0 {
				continue
			}
			discardMap := distributeResources(game.ResourceHandByPlayer(player), amount, room.Rand)
			game.DiscardPlayerCards(player, discardMap)
			formattedResources := wsUtils.FormatResources(discardMap)
			logs = append(logs, fmt.Sprintf("%s discarded %s", player, formattedResources))
		}
		handleDiscardCardsResponse(room, logs)
	}
}
