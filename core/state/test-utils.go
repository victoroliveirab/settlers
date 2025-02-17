package state

import (
	"fmt"
	"math/rand"
	"strconv"

	mapsdefinitions "github.com/victoroliveirab/settlers/core/maps"
	coreT "github.com/victoroliveirab/settlers/core/types"
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

func GetAllRounds() []int {
	return []int{
		SetupSettlement1,
		SetupRoad1,
		SetupSettlement2,
		SetupRoad2,
		FirstRound,
		Regular,
		MoveRobberDue7,
		PickRobbed,
		BetweenTurns,
		DiscardPhase,
	}
}

func CreateTestGame(opts ...GameStateOption) *GameState {
	mapsdefinitions.LoadMap("base4")
	var game GameState
	players := make([]*coreT.Player, 4)
	for i := 0; i < 4; i++ {
		players[i] = &coreT.Player{
			ID:    strconv.FormatInt(int64(i+1), 10),
			Name:  fmt.Sprintf("Player %d", i+1),
			Color: "color",
		}
	}

	game.New(players, "base4", 42, Params{
		MaxCards:       7,
		MaxSettlements: 5,
		MaxCities:      4,
		MaxRoads:       20,
	})

	for _, opt := range opts {
		opt(&game)
	}

	return &game
}

func MockWithRoundType(roundType int) GameStateOption {
	return func(gs *GameState) {
		gs.roundType = roundType
	}
}

func MockWithRoundNumber(roundNumber int) GameStateOption {
	return func(gs *GameState) {
		gs.roundNumber = roundNumber
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
		gs.playerResourceHandMap = resourcesByPlayer
	}
}

func MockWithSettlementsByPlayer(settlementsByPlayer map[string][]int) GameStateOption {
	return func(gs *GameState) {
		gs.playerSettlementMap = settlementsByPlayer
		gs.settlementMap = make(map[int]Building)
		for playerID, settlements := range settlementsByPlayer {
			for _, vertexID := range settlements {
				gs.settlementMap[vertexID] = Building{
					ID:    vertexID,
					Owner: playerID,
				}
			}
		}
	}
}

func MockWithCitiesByPlayer(citiesByPlayer map[string][]int) GameStateOption {
	return func(gs *GameState) {
		gs.playerCityMap = citiesByPlayer
		gs.cityMap = make(map[int]Building)
		for playerID, cities := range citiesByPlayer {
			for _, vertexID := range cities {
				gs.cityMap[vertexID] = Building{
					ID:    vertexID,
					Owner: playerID,
				}
			}
		}
	}
}

func MockWithRoadsByPlayer(roadsByPlayer map[string][]int) GameStateOption {
	return func(gs *GameState) {
		gs.playerRoadMap = roadsByPlayer
		gs.roadMap = make(map[int]Building)
		for playerID, roads := range roadsByPlayer {
			for _, edgeID := range roads {
				gs.roadMap[edgeID] = Building{
					ID:    edgeID,
					Owner: playerID,
				}
			}
		}
	}
}

func MockWithBlockedTile(tileID int) GameStateOption {
	return func(gs *GameState) {
		for _, tile := range gs.tiles {
			if tile.ID == tileID {
				tile.Blocked = true
			} else {
				tile.Blocked = false
			}
		}
	}
}

func MockWithRand(r *rand.Rand) GameStateOption {
	return func(gs *GameState) {
		gs.rand = r
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
