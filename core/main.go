package core

import (
	"fmt"
	"maps"
	"math/rand"
	"sort"

	coreMaps "github.com/victoroliveirab/settlers/core/maps"
	coreT "github.com/victoroliveirab/settlers/core/types"
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
	GameOver
)

var RoundTypeTranslation = [16]string{
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
	"GameOver",
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

type ResponseStatus string

const (
	NoResponse ResponseStatus = "Open"
	Accepted   ResponseStatus = "Accepted"
	Declined   ResponseStatus = "Declined"
	Countered  ResponseStatus = "Countered"
)

type TradePlayerEntry struct {
	Status  ResponseStatus `json:"status"`
	Blocked bool           `json:"blocked"`
}

type TradeStatus string

const (
	TradeOpen      TradeStatus = "Open"
	TradeClosed    TradeStatus = "Closed"
	TradeFinalized TradeStatus = "Finalized"
)

type Trade struct {
	ID        int                          `json:"id"`
	Requester string                       `json:"requester"`
	Creator   string                       `json:"creator"`
	Responses map[string]*TradePlayerEntry `json:"responses"`
	Offer     map[string]int               `json:"offer"`
	Request   map[string]int               `json:"request"`
	Status    TradeStatus                  `json:"status"`
	ParentID  int                          `json:"parent"`
	Finalized bool                         `json:"finalized"`
	Timestamp int64                        `json:"timestamp"`
}

type GameState struct {
	definition          coreMaps.MapDefinition
	mapName             string
	tiles               []*coreT.MapBlock
	ports               map[int]string
	rand                *rand.Rand
	players             []coreT.Player
	logs                []StateLog
	maxCards            int
	maxSettlements      int
	maxCities           int
	maxRoads            int
	maxDevCardsPerRound int
	bankTradeAmount     int

	// cost related
	generalPortCost  int
	resourcePortCost int

	// points related
	targetPoint          int
	points               map[string]int
	longestRoad          LongestRoad
	mostKnights          MostKnights
	pointsPerCity        int
	pointsPerSettlement  int
	pointsPerMostKnights int
	pointsPerLongestRoad int
	mostKnightsMinimum   int
	longestRoadMinimum   int

	// round related
	roundType                           int
	roundNumber                         int
	currentPlayerIndex                  int
	dice1                               int
	dice2                               int
	currentPlayerNumberOfPlayedDevCards int

	// discard related
	discardMap                     map[string]int
	playerDiscardedCurrentRoundMap map[string]bool

	// cards related
	playerResourceHandMap        map[string]map[string]int
	playerDevelopmentHandMap     map[string]map[string][]*coreT.DevelopmentCard
	developmentCards             []*coreT.DevelopmentCard
	developmentCardHeadIndex     int
	playerDevelopmentCardUsedMap map[string]map[string]int

	// building related
	playerSettlementMap map[string][]int
	playerCityMap       map[string][]int
	playerRoadMap       map[string][]int
	playerPortMap       map[string][]int
	playerLongestRoad   map[string][]int

	// book keeping
	cityMap            map[int]Building
	roadMap            map[int]Building
	settlementMap      map[int]Building
	playersTrades      map[int]*Trade
	tradeParentToChild map[int][]int
	playerTradeId      int
}

type Params struct {
	Speed                int
	BankTradeAmount      int
	MaxCards             int
	MaxDevCardsPerRound  int
	MaxSettlements       int
	MaxCities            int
	MaxRoads             int
	TargetPoint          int
	PointsPerSettlement  int
	PointsPerCity        int
	PointsForMostKnights int
	PointsForLongestRoad int
	MostKnightsMinimum   int
	LongestRoadMinimum   int
}

func (state *GameState) New(players []*coreT.Player, mapName string, randGenerator *rand.Rand, params Params) error {
	state.rand = randGenerator
	data, err := coreMaps.GenerateMap(mapName, state.rand)
	if err != nil {
		return err
	}

	state.definition = data.Definition
	state.mapName = mapName
	state.tiles = data.Tiles
	state.ports = data.Ports
	state.developmentCards = data.DevelopmentCards
	state.developmentCardHeadIndex = 0
	state.playerDevelopmentCardUsedMap = make(map[string]map[string]int)
	state.players = make([]coreT.Player, len(players))
	for i, player := range players {
		state.players[i] = coreT.Player{
			ID:    player.ID,
			Color: player.Color,
		}
	}
	state.logs = make([]StateLog, 0)
	state.maxCards = params.MaxCards
	state.maxSettlements = params.MaxSettlements
	state.maxCities = params.MaxCities
	state.maxRoads = params.MaxRoads
	state.maxDevCardsPerRound = params.MaxDevCardsPerRound
	state.bankTradeAmount = params.BankTradeAmount
	state.pointsPerSettlement = params.PointsPerSettlement
	state.pointsPerCity = params.PointsPerCity
	state.pointsPerMostKnights = params.PointsForMostKnights
	state.pointsPerLongestRoad = params.PointsForLongestRoad
	state.mostKnightsMinimum = params.MostKnightsMinimum
	state.longestRoadMinimum = params.LongestRoadMinimum
	state.generalPortCost = 3
	state.resourcePortCost = 2

	state.targetPoint = params.TargetPoint
	state.points = make(map[string]int)

	state.roundType = SetupSettlement1
	state.roundNumber = 0
	state.currentPlayerIndex = 0
	state.currentPlayerNumberOfPlayedDevCards = 0
	state.dice1 = 0
	state.dice2 = 0
	state.discardMap = make(map[string]int)
	state.playerDiscardedCurrentRoundMap = make(map[string]bool)

	state.playerResourceHandMap = make(map[string]map[string]int)
	state.playerDevelopmentHandMap = make(map[string]map[string][]*coreT.DevelopmentCard)

	state.playerSettlementMap = make(map[string][]int)
	state.playerCityMap = make(map[string][]int)
	state.playerRoadMap = make(map[string][]int)
	state.playerPortMap = make(map[string][]int)
	state.playerLongestRoad = make(map[string][]int)

	state.cityMap = make(map[int]Building)
	state.roadMap = make(map[int]Building)
	state.settlementMap = make(map[int]Building)

	state.playersTrades = make(map[int]*Trade)
	state.tradeParentToChild = make(map[int][]int)
	state.playerTradeId = 0

	for _, player := range players {
		state.discardMap[player.ID] = 0
		state.playerSettlementMap[player.ID] = make([]int, 0)
		state.playerCityMap[player.ID] = make([]int, 0)
		state.playerRoadMap[player.ID] = make([]int, 0)
		state.playerPortMap[player.ID] = make([]int, 0)
		state.playerLongestRoad[player.ID] = make([]int, 0)
		state.playerResourceHandMap[player.ID] = make(map[string]int)
		state.playerDevelopmentHandMap[player.ID] = make(map[string][]*coreT.DevelopmentCard)
		state.playerDevelopmentCardUsedMap[player.ID] = map[string]int{
			"Knight":         0,
			"Monopoly":       0,
			"Road Building":  0,
			"Year of Plenty": 0,
		}

		state.playerResourceHandMap[player.ID]["Lumber"] = 0
		state.playerResourceHandMap[player.ID]["Brick"] = 0
		state.playerResourceHandMap[player.ID]["Grain"] = 0
		state.playerResourceHandMap[player.ID]["Sheep"] = 0
		state.playerResourceHandMap[player.ID]["Ore"] = 0

		state.playerDevelopmentHandMap[player.ID]["Knight"] = make([]*coreT.DevelopmentCard, 0)
		state.playerDevelopmentHandMap[player.ID]["Victory Point"] = make([]*coreT.DevelopmentCard, 0)
		state.playerDevelopmentHandMap[player.ID]["Road Building"] = make([]*coreT.DevelopmentCard, 0)
		state.playerDevelopmentHandMap[player.ID]["Year of Plenty"] = make([]*coreT.DevelopmentCard, 0)
		state.playerDevelopmentHandMap[player.ID]["Monopoly"] = make([]*coreT.DevelopmentCard, 0)

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

func (state *GameState) Map() []*coreT.MapBlock {
	// REFACTOR: return a copy
	return state.tiles
}

func (state *GameState) MapName() string {
	return state.mapName
}

func (state *GameState) PortsByVertex() map[int]string {
	return state.ports
}

func (state *GameState) PortsLocations() [][2]int {
	return state.definition.PortsLocations
}

func (state *GameState) Ports() []coreT.Port {
	ports := make([]coreT.Port, 0)
	for _, location := range state.definition.PortsLocations {
		ports = append(ports, coreT.Port{
			Type:     state.ports[location[0]],
			Vertices: location,
		})
	}
	return ports
}

func (state *GameState) PortsByPlayer(playerID string) []string {
	ports := make([]string, 0)
	for vertexID, kind := range state.ports {
		settlement, okSettlement := state.settlementMap[vertexID]
		city, okCity := state.cityMap[vertexID]
		if okSettlement && settlement.Owner == playerID {
			ports = append(ports, kind)
		} else if okCity && city.Owner == playerID {
			ports = append(ports, kind)
		}
	}
	return ports
}

func (state *GameState) Players() []coreT.Player {
	// REFACTOR: return a copy
	return state.players
}

func (state *GameState) CurrentRoundPlayer() coreT.Player {
	return state.players[state.currentPlayerIndex]
}

func (state *GameState) CurrentRoundPlayerIndex() int {
	return state.currentPlayerIndex
}

func (state *GameState) Dice() [2]int {
	return [2]int{state.dice1, state.dice2}
}

func (state *GameState) RoundType() int {
	return state.roundType
}

func (state *GameState) DevelopmentHandByPlayer(playerID string) map[string]int {
	devHand := make(map[string]int)
	for name, cards := range state.playerDevelopmentHandMap[playerID] {
		devHand[name] = len(cards)
	}
	return devHand
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

func (state *GameState) LongestRoadLengths() map[string]int {
	longestRoadByPlayer := make(map[string]int)
	for _, player := range state.players {
		longestRoadByPlayer[player.ID] = state.LongestRoadLengthByPlayer(player.ID)
	}
	return longestRoadByPlayer
}

func (state *GameState) LongestRoadLengthByPlayer(playerID string) int {
	return len(state.playerLongestRoad[playerID])
}

func (state *GameState) KnightUses() map[string]int {
	knightUses := make(map[string]int)
	for playerID, uses := range state.playerDevelopmentCardUsedMap {
		knightUses[playerID] = uses["Knight"]
	}
	return knightUses
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
	return state.discardAmountByPlayer(playerID)
}

func (state *GameState) DiscardAmounts() map[string]int {
	amounts := map[string]int{}
	for _, player := range state.players {
		amounts[player.ID] = state.DiscardAmountByPlayer(player.ID)
	}
	return amounts
}

func (state *GameState) Trades() []Trade {
	trades := make([]Trade, 0)
	for _, trade := range state.playersTrades {
		trades = append(trades, *trade)
	}

	sort.Slice(trades, func(i, j int) bool {
		return trades[i].ID < trades[j].ID
	})

	return trades
}

func (state *GameState) ActiveTradeOffers() []Trade {
	activeTrades := make([]Trade, 0)
	for _, trade := range state.playersTrades {
		if trade.Status == TradeOpen {
			activeTrades = append(activeTrades, *trade)
		}
	}

	sort.Slice(activeTrades, func(i, j int) bool {
		return activeTrades[i].ID < activeTrades[j].ID
	})

	return activeTrades
}

func (state *GameState) Round() int {
	return state.roundNumber
}
