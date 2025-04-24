package core

import (
	"fmt"
	"math/rand"

	coreMaps "github.com/victoroliveirab/settlers/core/maps"
	"github.com/victoroliveirab/settlers/core/packages/board"
	"github.com/victoroliveirab/settlers/core/packages/development"
	"github.com/victoroliveirab/settlers/core/packages/player"
	"github.com/victoroliveirab/settlers/core/packages/round"
	"github.com/victoroliveirab/settlers/core/packages/statistics"
	"github.com/victoroliveirab/settlers/core/packages/trade"
	coreT "github.com/victoroliveirab/settlers/core/types"
	"github.com/victoroliveirab/settlers/utils"
)

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

type GameState struct {
	board               *board.Instance
	rand                *rand.Rand
	stats               *statistics.Instance
	logs                []StateLog
	maxCards            int
	maxSettlements      int
	maxCities           int
	maxRoads            int
	maxDevCardsPerRound int
	bankTradeAmount     int

	// player
	players       []coreT.Player
	playersStates map[string]*player.Instance

	// trade
	trade *trade.Instance

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
	round              *round.Instance
	currentPlayerIndex int

	// cards related
	development *development.Instance
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
	mapDefinitions, err := coreMaps.GetMapDefinitions(mapName)
	if err != nil {
		return err
	}

	state.board = board.New(mapName, mapDefinitions, randGenerator)
	state.stats = statistics.New(players)

	developmentCards := utils.MapToShuffledSlice[*coreT.DevelopmentCard](
		mapDefinitions.DevelopmentCards,
		func(el string) *coreT.DevelopmentCard { return &coreT.DevelopmentCard{Name: el} },
		randGenerator,
	)
	state.development = development.New(developmentCards)

	state.playersStates = make(map[string]*player.Instance)

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

	state.round = round.New()
	state.currentPlayerIndex = 0

	state.players = make([]coreT.Player, len(players))
	for i, playerDefinition := range players {
		state.players[i] = coreT.Player{
			ID:    playerDefinition.ID,
			Color: playerDefinition.Color,
		}
		state.playersStates[playerDefinition.ID] = player.New(playerDefinition, map[string]int{
			"Lumber": 0,
			"Brick":  0,
			"Sheep":  0,
			"Grain":  0,
			"Ore":    0,
		}, map[string][]*coreT.DevelopmentCard{})
	}

	state.trade = trade.New(state.stats)

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

func (state *GameState) MapName() string {
	return state.board.MapName
}

func (state *GameState) PortsByVertex() map[int]string {
	return state.board.Ports
}

func (state *GameState) PortsLocations() [][2]int {
	return state.board.Definition.PortsLocations
}

func (state *GameState) Ports() []coreT.Port {
	ports := make([]coreT.Port, 0)
	for _, location := range state.board.Definition.PortsLocations {
		ports = append(ports, coreT.Port{
			Type:     state.board.Ports[location[0]],
			Vertices: location,
		})
	}
	return ports
}

func (state *GameState) PortsByPlayer(playerID string) []string {
	return state.playersStates[playerID].GetPortTypes()
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
	return state.round.GetDice()
}

func (state *GameState) RoundType() round.Type {
	return state.round.GetRoundType()
}

func (state *GameState) DevelopmentHandByPlayer(playerID string) map[string]int {
	devHand := make(map[string]int)
	playerState := state.playersStates[playerID]
	for name, cards := range playerState.GetDevelopmentCards() {
		devHand[name] = len(cards)
	}
	return devHand
}

func (state *GameState) ResourceHandByPlayer(playerID string) map[string]int {
	playerState := state.playersStates[playerID]
	return playerState.GetResources()
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
	playerState := state.playersStates[playerID]
	return playerState.GetSettlements()
}

func (state *GameState) CitiesByPlayer(playerID string) []int {
	playerState := state.playersStates[playerID]
	return playerState.GetCities()
}

func (state *GameState) RoadsByPlayer(playerID string) []int {
	playerState := state.playersStates[playerID]
	return playerState.GetRoads()
}

func (state *GameState) LongestRoadLengths() map[string]int {
	longestRoadByPlayer := make(map[string]int)
	for _, player := range state.players {
		longestRoadByPlayer[player.ID] = state.LongestRoadLengthByPlayer(player.ID)
	}
	return longestRoadByPlayer
}

func (state *GameState) LongestRoadLengthByPlayer(playerID string) int {
	playerState := state.playersStates[playerID]
	return playerState.GetLongestRoadSize()
}

func (state *GameState) KnightUses() map[string]int {
	knightUses := make(map[string]int)
	for _, player := range state.players {
		playerState := state.playersStates[player.ID]
		knightUses[player.ID] = playerState.GetKnightCount()
	}
	return knightUses
}

func (state *GameState) RobbablePlayers(playerID string) ([]string, error) {
	keys := make([]string, 0)

	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot check robbable players during other player's round")
		return keys, err
	}

	if state.round.GetRoundType() != round.PickRobbed {
		err := fmt.Errorf("Cannot check robbable players outside PickRobbed round type")
		return keys, err
	}

	robbablePlayers := make(map[string]bool)
	settlements := state.board.GetSettlements()
	cities := state.board.GetCities()
	for _, tile := range state.board.GetTiles() {
		if tile.Blocked {
			vertices := state.board.Definition.VerticesByTile[tile.ID]
			for _, vertexID := range vertices {
				settlement, hasSettlement := settlements[vertexID]
				city, hasCity := cities[vertexID]
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
	if state.round.GetRoundType() != round.DiscardPhase {
		return 0
	}
	playerState := state.playersStates[playerID]
	if playerState.GetHasDiscardedThisTurn() {
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

func (state *GameState) Trades() []trade.Trade {
	return state.trade.Trades()
}

func (state *GameState) ActiveTradeOffers() []trade.Trade {
	return state.trade.ActiveTrades()
}

func (state *GameState) Round() int {
	return state.round.GetRoundNumber()
}

func (state *GameState) EndGame() {
	state.round.SetRoundType(round.GameOver)
	state.trade.CancelActiveTrades()
	state.stats.AddPointsRecord(state.points)
}
