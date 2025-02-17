package state

import (
	"maps"
	"math/rand"

	coreMaps "github.com/victoroliveirab/settlers/core/maps"
	coreT "github.com/victoroliveirab/settlers/core/types"
	"github.com/victoroliveirab/settlers/utils"
)

const (
	SetupSettlement1 int = iota
	SetupRoad1
	SetupSettlement2
	SetupRoad2
	FirstRound
	Regular
	MoveRobberDue7
	PickRobbed
	BetweenTurns
	DiscardPhase
)

var RoundTypeTranslation = [10]string{
	"SettlementSetup#1",
	"RoadSetup#1",
	"SettlementSetup#2",
	"RoadSetup#2",
	"FirstRound",
	"Regular",
	"MoveRobber(7)",
	"ChooseRobbedPlayer",
	"BetweenRounds",
	"DiscardPhase",
}

type Building struct {
	ID    int    `json:"id"`
	Owner string `json:"owner"`
}

type StateLog struct {
	Timestamp string
	Message   string
}

type GameState struct {
	definition     coreMaps.MapDefinition
	tiles          []*coreT.MapBlock
	rand           *rand.Rand
	players        []coreT.Player
	logs           []StateLog
	maxCards       int
	targetPoint    int
	maxSettlements int
	maxCities      int
	maxRoads       int

	// round related
	roundType          int
	roundNumber        int
	currentPlayerIndex int
	dice1              int
	dice2              int

	// cards related
	playerResourceHandMap    map[string]map[string]int
	playerDevelopmentHandMap map[string]map[string]int

	// building related
	playerSettlementMap map[string][]int
	playerCityMap       map[string][]int
	playerRoadMap       map[string][]int

	// book keeping
	cityMap       map[int]Building
	roadMap       map[int]Building
	settlementMap map[int]Building
}

type Params struct {
	MaxCards       int
	MaxSettlements int
	MaxCities      int
	MaxRoads       int
	TargetPoint    int
}

func (state *GameState) New(players []*coreT.Player, mapName string, seed int, params Params) error {
	state.rand = utils.RandNew(int64(seed))
	data, err := coreMaps.GenerateMap(mapName, state.rand)
	if err != nil {
		return err
	}

	state.definition = data.Definition
	state.tiles = data.Tiles
	state.players = make([]coreT.Player, len(players))
	for i, player := range players {
		state.players[i] = coreT.Player{
			ID:    player.ID,
			Name:  player.Name,
			Color: player.Color,
		}
	}
	state.logs = make([]StateLog, 0)
	state.maxCards = params.MaxCards
	state.targetPoint = params.TargetPoint
	state.maxSettlements = params.MaxSettlements
	state.maxCities = params.MaxCities
	state.maxRoads = params.MaxRoads

	state.roundType = SetupSettlement1
	state.roundNumber = 0
	state.currentPlayerIndex = 0
	state.dice1 = 0
	state.dice2 = 0

	state.playerResourceHandMap = make(map[string]map[string]int)
	state.playerDevelopmentHandMap = make(map[string]map[string]int)

	state.playerSettlementMap = make(map[string][]int)
	state.playerCityMap = make(map[string][]int)
	state.playerRoadMap = make(map[string][]int)

	state.cityMap = make(map[int]Building)
	state.roadMap = make(map[int]Building)
	state.settlementMap = make(map[int]Building)

	for _, player := range players {
		state.playerSettlementMap[player.ID] = make([]int, 0)
		state.playerCityMap[player.ID] = make([]int, 0)
		state.playerRoadMap[player.ID] = make([]int, 0)
		state.playerResourceHandMap[player.ID] = make(map[string]int)
		state.playerDevelopmentHandMap[player.ID] = make(map[string]int)

		state.playerResourceHandMap[player.ID]["Lumber"] = 0
		state.playerResourceHandMap[player.ID]["Brick"] = 0
		state.playerResourceHandMap[player.ID]["Grain"] = 0
		state.playerResourceHandMap[player.ID]["Sheep"] = 0
		state.playerResourceHandMap[player.ID]["Ore"] = 0
	}

	return nil
}

func (state *GameState) currentPlayer() *coreT.Player {
	return &state.players[state.currentPlayerIndex]
}

func (state *GameState) findPlayer(playerID string) *coreT.Player {
	for _, player := range state.players {
		if player.ID == playerID {
			return &player
		}
	}
	return nil
}

// Getters
func (state *GameState) RoundType() int {
	return state.roundType
}

func (state *GameState) ResourceHandByPlayer(playerID string) map[string]int {
	return state.playerResourceHandMap[playerID]
}

func (state *GameState) SettlementsByPlayer(playerID string) []int {
	return state.playerSettlementMap[playerID]
}

func (state *GameState) AllSettlements() map[int]Building {
	return maps.Clone(state.settlementMap)
}

func (state *GameState) CitiesByPlayer(playerID string) []int {
	return state.playerCityMap[playerID]
}

func (state *GameState) AllCities() map[int]Building {
	return maps.Clone(state.cityMap)
}

func (state *GameState) RoadsByPlayer(playerID string) []int {
	return state.playerRoadMap[playerID]
}

func (state *GameState) AllRoads() map[int]Building {
	return maps.Clone(state.roadMap)
}
