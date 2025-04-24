//go:build test

package core

import (
	"math/rand"
	"strconv"

	mapsdefinitions "github.com/victoroliveirab/settlers/core/maps"
	"github.com/victoroliveirab/settlers/core/packages/round"
	coreT "github.com/victoroliveirab/settlers/core/types"
	"github.com/victoroliveirab/settlers/utils"
)

type GameStateMock struct {
	RoundType           int
	RoundNumber         int
	CurrentPlayerID     string
	ResourcesByPlayer   map[string]map[string]int
	SettlementsByPlayer map[string][]int
	RoadsByPlayer       map[string][]int
	CitiesByPlayer      map[string][]int
}

type GameStateOption func(*GameState)

func CreateTestGame(opts ...GameStateOption) *GameState {
	mapsdefinitions.LoadMap("base4")
	var game GameState
	players := make([]*coreT.Player, 4)
	for i := 0; i < 4; i++ {
		players[i] = &coreT.Player{
			ID: strconv.FormatInt(int64(i+1), 10),
			Color: coreT.PlayerColor{
				Background: "bg",
				Foreground: "fg",
			},
		}
	}

	randGenerator := utils.RandNew(42)

	game.New(players, "base4", randGenerator, Params{
		BankTradeAmount:      4,
		MaxCards:             7,
		MaxSettlements:       5,
		MaxCities:            4,
		MaxRoads:             20,
		MaxDevCardsPerRound:  1,
		TargetPoint:          10,
		PointsPerSettlement:  1,
		PointsPerCity:        2,
		PointsForMostKnights: 2,
		PointsForLongestRoad: 2,
		LongestRoadMinimum:   5,
		MostKnightsMinimum:   3,
	})

	for _, opt := range opts {
		opt(&game)
	}

	return &game
}

func CreateTestGameWithRand(randGenerator *rand.Rand, opts ...GameStateOption) *GameState {
	mapsdefinitions.LoadMap("base4")
	var game GameState
	players := make([]*coreT.Player, 4)
	for i := 0; i < 4; i++ {
		players[i] = &coreT.Player{
			ID: strconv.FormatInt(int64(i+1), 10),
			Color: coreT.PlayerColor{
				Background: "bg",
				Foreground: "fg",
			},
		}
	}

	game.New(players, "base4", randGenerator, Params{
		BankTradeAmount:      4,
		MaxCards:             7,
		MaxSettlements:       5,
		MaxCities:            4,
		MaxRoads:             20,
		MaxDevCardsPerRound:  1,
		TargetPoint:          10,
		PointsPerSettlement:  1,
		PointsPerCity:        2,
		PointsForMostKnights: 2,
		PointsForLongestRoad: 2,
		LongestRoadMinimum:   5,
		MostKnightsMinimum:   3,
	})

	for _, opt := range opts {
		opt(&game)
	}

	return &game
}

func MockWithRoundType(roundType round.Type) GameStateOption {
	return func(gs *GameState) {
		gs.round.SetRoundType(roundType)
		// Hack to have a dice set since the round type is already regular
		if roundType == round.Regular {
			gs.round.SetDice(1, 1)
		}
	}
}

func MockWithRoundNumber(roundNumber int) GameStateOption {
	return func(gs *GameState) {
		gs.round.SetRoundNumber(roundNumber)
	}
}

func MockWithCurrentRoundPlayer(playerID string) GameStateOption {
	return func(gs *GameState) {
		for i, player := range gs.players {
			if player.ID == playerID {
				gs.currentPlayerIndex = i
				return
			}
		}
	}
}

func MockWithResourcesByPlayer(resourcesByPlayer map[string]map[string]int) GameStateOption {
	return func(gs *GameState) {
		for _, player := range gs.players {
			playerState := gs.playersStates[player.ID]
			playerState.SetResources(resourcesByPlayer[player.ID])
		}
	}
}

func MockWithDevelopmentsByPlayer(developmentCardsByPlayer map[string]map[string][]*coreT.DevelopmentCard) GameStateOption {
	return func(gs *GameState) {
		for _, player := range gs.players {
			playerState := gs.playersStates[player.ID]
			playerState.SetDevelopmentCards(developmentCardsByPlayer[player.ID])
		}
	}
}

func MockWithSettlementsByPlayer(settlementsByPlayer map[string][]int) GameStateOption {
	return func(gs *GameState) {
		for _, player := range gs.players {
			playerState := gs.playersStates[player.ID]
			for _, vertexID := range settlementsByPlayer[player.ID] {
				playerState.AddSettlement(vertexID)
				gs.board.AddSettlement(player.ID, vertexID)
			}
		}
	}
}

func MockWithCitiesByPlayer(citiesByPlayer map[string][]int) GameStateOption {
	return func(gs *GameState) {
		for _, player := range gs.players {
			playerState := gs.playersStates[player.ID]
			for _, vertexID := range citiesByPlayer[player.ID] {
				playerState.AddCity(vertexID)
				gs.board.AddCity(player.ID, vertexID)
			}
		}
	}
}

// NOTE: not prepared to be used in conjuction with MockWithSettlementsByPlayer and MockWithCitiesByPlayer
// If ports are needed during test, only use this mock function
func MockWithPortsByPlayer(portsByPlayer map[string][]string) GameStateOption {
	return func(gs *GameState) {
		for playerID, ports := range portsByPlayer {
			playerState := gs.playersStates[playerID]
			for _, port := range ports {
				vertexID := -1
				for candidateVertexID, portType := range gs.board.Ports {
					if portType == port {
						vertexID = candidateVertexID
					}
				}
				if vertexID == -1 {
					panic(portsByPlayer)
				}
				gs.board.AddSettlement(playerID, vertexID)
				playerState.AddSettlement(vertexID)
				playerState.AddPort(vertexID, gs.board.Ports[vertexID])
			}
		}
	}
}

func MockWithRoadsByPlayer(roadsByPlayer map[string][]int) GameStateOption {
	return func(gs *GameState) {
		for _, player := range gs.players {
			playerState := gs.playersStates[player.ID]
			for _, edgeID := range roadsByPlayer[player.ID] {
				playerState.AddRoad(edgeID)
				gs.board.AddRoad(player.ID, edgeID)
			}
			gs.computeLongestRoad(player.ID)
		}
	}
}

func MockWithBlockedTile(tileID int) GameStateOption {
	return func(gs *GameState) {
		for i, tile := range gs.board.GetTiles() {
			if tile.ID == tileID {
				gs.board.BlockTileByIndex(i)
			} else {
				gs.board.UnblockTileByIndex(i)
			}
		}
	}
}

func MockWithUsedDevelopmentCardsByPlayer(developmentCardsUsedByPlayer map[string]map[string]int) GameStateOption {
	return func(gs *GameState) {
		for playerID, devCards := range developmentCardsUsedByPlayer {
			gs.playersStates[playerID].SetUsedDevelopmentCards(devCards)
		}
	}
}

func MockWithPoints() GameStateOption {
	return func(gs *GameState) {
		for _, player := range gs.players {
			gs.computeLongestRoad(player.ID)
		}
		gs.recountLongestRoad()
		gs.recountKnights()
		gs.updatePoints()
	}
}

func MockWithRand(r *rand.Rand) GameStateOption {
	return func(gs *GameState) {
		gs.rand = r
	}
}

func MockWithNextDevelopmentCard(name string) GameStateOption {
	return func(gs *GameState) {
		gs.development.SetCardByIndex(0, name)
	}
}

func StubRand(desiredSum int) *rand.Rand {
	seedByDesiredSum := map[int]int64{
		2:  56,
		3:  16,
		4:  7,
		5:  15,
		6:  2,
		7:  4,
		8:  10,
		9:  13,
		10: 1,
		11: 3,
		12: 42,
	}
	stub := rand.NewSource(seedByDesiredSum[desiredSum])
	r := rand.New(stub)
	return r
}
