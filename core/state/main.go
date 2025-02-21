package state

import (
	"fmt"
	"maps"
	"math"
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
	MoveRobberDueKnight
	PickRobbed
	BetweenTurns
	BuildRoad1Development
	BuildRoad2Development
	MonopolyPickResource
	YearOfPlentyPickResources
	DiscardPhase
)

var RoundTypeTranslation = [15]string{
	"SettlementSetup#1",
	"RoadSetup#1",
	"SettlementSetup#2",
	"RoadSetup#2",
	"FirstRound",
	"Regular",
	"MoveRobber(7)",
	"MoveRobber(Knight)",
	"ChooseRobbedPlayer",
	"BetweenRounds",
	"BuildRoadDevelopment(1)",
	"BuildRoadDevelopment(2)",
	"MonopolyPickResource",
	"YearOfPlentyPickResources",
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

type LongestRoad struct {
	PlayerID string
	Length   int
}

type MostKnights struct {
	PlayerID string
	Quantity int
}

type TradePlayerEntry struct {
	Status  string // "Open" | "Accepted" | "Declined"
	Blocked bool
}

type Trade struct {
	ID        int
	PlayerID  string
	Opponents map[string]*TradePlayerEntry
	Offer     map[string]int
	Request   map[string]int
	Status    string // "Open" | "Closed"
	Counters  []int
	ParentID  int
	Finalized bool
	Timestamp int64
}

type GameState struct {
	definition          coreMaps.MapDefinition
	tiles               []*coreT.MapBlock
	rand                *rand.Rand
	players             []coreT.Player
	logs                []StateLog
	maxCards            int
	maxSettlements      int
	maxCities           int
	maxRoads            int
	maxDevCardsPerRound int

	// points related
	targetPoint int
	points      map[string]int
	longestRoad LongestRoad
	mostKnights MostKnights

	// round related
	roundType                           int
	roundNumber                         int
	currentPlayerIndex                  int
	dice1                               int
	dice2                               int
	currentPlayerNumberOfPlayedDevCards int

	// discard related
	discardMap map[string]int

	// cards related
	playerResourceHandMap    map[string]map[string]int
	playerDevelopmentHandMap map[string]map[string]int
	developmentCards         []*coreT.DevelopmentCard
	developmentCardHeadIndex int

	// building related
	playerSettlementMap map[string][]int
	playerCityMap       map[string][]int
	playerRoadMap       map[string][]int
	playerLongestRoad   map[string][]int

	// book keeping
	cityMap       map[int]Building
	roadMap       map[int]Building
	settlementMap map[int]Building
	playersTrades map[int]*Trade
	playerTradeId int
}

type Params struct {
	MaxCards            int
	MaxDevCardsPerRound int
	MaxSettlements      int
	MaxCities           int
	MaxRoads            int
	TargetPoint         int
}

func (state *GameState) New(players []*coreT.Player, mapName string, seed int, params Params) error {
	state.rand = utils.RandNew(int64(seed))
	data, err := coreMaps.GenerateMap(mapName, state.rand)
	if err != nil {
		return err
	}

	state.definition = data.Definition
	state.tiles = data.Tiles
	state.developmentCards = data.DevelopmentCards
	state.developmentCardHeadIndex = 0
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
	state.maxSettlements = params.MaxSettlements
	state.maxCities = params.MaxCities
	state.maxRoads = params.MaxRoads
	state.maxDevCardsPerRound = params.MaxDevCardsPerRound

	state.targetPoint = params.TargetPoint
	state.points = make(map[string]int)

	state.roundType = SetupSettlement1
	state.roundNumber = 0
	state.currentPlayerIndex = 0
	state.currentPlayerNumberOfPlayedDevCards = 0
	state.dice1 = 0
	state.dice2 = 0
	state.discardMap = make(map[string]int)

	state.playerResourceHandMap = make(map[string]map[string]int)
	state.playerDevelopmentHandMap = make(map[string]map[string]int)

	state.playerSettlementMap = make(map[string][]int)
	state.playerCityMap = make(map[string][]int)
	state.playerRoadMap = make(map[string][]int)
	state.playerLongestRoad = make(map[string][]int)

	state.cityMap = make(map[int]Building)
	state.roadMap = make(map[int]Building)
	state.settlementMap = make(map[int]Building)

	state.playersTrades = make(map[int]*Trade)
	state.playerTradeId = 0

	for _, player := range players {
		state.discardMap[player.ID] = 0
		state.playerSettlementMap[player.ID] = make([]int, 0)
		state.playerCityMap[player.ID] = make([]int, 0)
		state.playerRoadMap[player.ID] = make([]int, 0)
		state.playerLongestRoad[player.ID] = make([]int, 0)
		state.playerResourceHandMap[player.ID] = make(map[string]int)
		state.playerDevelopmentHandMap[player.ID] = make(map[string]int)

		state.playerResourceHandMap[player.ID]["Lumber"] = 0
		state.playerResourceHandMap[player.ID]["Brick"] = 0
		state.playerResourceHandMap[player.ID]["Grain"] = 0
		state.playerResourceHandMap[player.ID]["Sheep"] = 0
		state.playerResourceHandMap[player.ID]["Ore"] = 0

		state.playerDevelopmentHandMap[player.ID]["Knight"] = 0
		state.playerDevelopmentHandMap[player.ID]["Victory Point"] = 0
		state.playerDevelopmentHandMap[player.ID]["Road Building"] = 0
		state.playerDevelopmentHandMap[player.ID]["Year of Plenty"] = 0
		state.playerDevelopmentHandMap[player.ID]["Monopoly"] = 0

		state.points[player.ID] = 0
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
func (state *GameState) Dice() [2]int {
	return [2]int{state.dice1, state.dice2}
}

func (state *GameState) RoundType() int {
	return state.roundType
}

func (state *GameState) DevelopmentHandByPlayer(playerID string) map[string]int {
	return state.playerDevelopmentHandMap[playerID]
}

func (state *GameState) ResourceHandByPlayer(playerID string) map[string]int {
	return state.playerResourceHandMap[playerID]
}

func (state *GameState) NumberOfCardsInHandByPlayer(playerID string) int {
	hand := state.ResourceHandByPlayer(playerID)
	sum := 0
	for _, count := range hand {
		sum += count
	}
	return sum
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

func (state *GameState) LongestRoadLengthByPlayer(playerID string) int {
	return len(state.playerLongestRoad[playerID])
}

func (state *GameState) RobbablePlayers(playerID string) ([]string, error) {
	keys := make([]string, 0)

	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot check robbable players during other player's round")
		return keys, err
	}

	if state.roundType != PickRobbed {
		err := fmt.Errorf("Cannot check robbable players outside PickRobbed round type")
		return keys, err
	}

	robbablePlayers := make(map[string]bool)
	for _, tile := range state.tiles {
		if tile.Blocked {
			vertices := state.definition.VerticesByTile[tile.ID]
			for _, vertexID := range vertices {
				settlement, hasSettlement := state.settlementMap[vertexID]
				city, hasCity := state.cityMap[vertexID]
				if hasSettlement {
					robbablePlayers[settlement.Owner] = true
				}
				if hasCity {
					robbablePlayers[city.Owner] = true
				}
			}
		}
	}
	for ownerID := range robbablePlayers {
		if ownerID != playerID {
			keys = append(keys, ownerID)
		}
	}
	return keys, nil
}

func (state *GameState) DiscardAmountByPlayer(playerID string) int {
	if state.roundType != DiscardPhase {
		return 0
	}
	total := 0
	for _, count := range state.playerResourceHandMap[playerID] {
		total += count
	}
	if total <= state.maxCards {
		return 0
	}
	return int(math.Floor(float64(total) / 2))
}

func (state *GameState) ActiveTradeOffers() []Trade {
	activeTrades := make([]Trade, 0)
	for _, trade := range state.playersTrades {
		if !trade.Finalized {
			activeTrades = append(activeTrades, *trade)
		}
	}
	return activeTrades
}
