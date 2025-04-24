package bookkeeping

import (
	"maps"

	coreT "github.com/victoroliveirab/settlers/core/types"
)

type Instance struct {
	dice                         map[int]int
	diceByPlayer                 map[string]map[int]int
	longestRoadEvolutionPerRound map[string][]int
	numberOfRobberiesByPlayer    map[string]int
	numberOfTimesRobbedByPlayer  map[string]int
	resourcesDiscardedByPlayer   map[string]map[string]int
	resourcesDrawnByPlayer       map[string]map[string]int
	resourcesBlockedByPlayer     map[string]map[string]int
	resourcesUsedByPlayer        map[string]map[string]int
	devCardsDrawnByPlayer        map[string]map[string]int
	pointsEvolutionPerRound      map[string][]int
	tradesByPlayer               map[string]map[string]int
}

var resources map[string]int = map[string]int{
	"Lumber": 0,
	"Brick":  0,
	"Sheep":  0,
	"Grain":  0,
	"Ore":    0,
}

func New(players []*coreT.Player) *Instance {
	dice := map[int]int{
		2:  0,
		3:  0,
		4:  0,
		5:  0,
		6:  0,
		7:  0,
		8:  0,
		9:  0,
		10: 0,
		11: 0,
		12: 0,
	}
	numberOfRobberiesByPlayer := make(map[string]int)
	numberOfTimesRobbedByPlayer := make(map[string]int)
	resourcesBlockedByPlayer := make(map[string]map[string]int)
	resourcesUsedByPlayer := make(map[string]map[string]int)
	diceStatsByPlayer := make(map[string]map[int]int)
	resourcesDiscardedByPlayer := make(map[string]map[string]int)
	resourcesDrawnByPlayer := make(map[string]map[string]int)
	devCardsDrawnByPlayer := make(map[string]map[string]int)
	longestRoadEvolutionPerRound := make(map[string][]int)
	pointsPerRound := make(map[string][]int)
	tradesByPlayer := make(map[string]map[string]int)

	for _, player := range players {
		playerID := player.ID
		numberOfRobberiesByPlayer[playerID] = 0
		numberOfTimesRobbedByPlayer[playerID] = 0
		resourcesBlockedByPlayer[playerID] = maps.Clone(resources)
		diceStatsByPlayer[playerID] = maps.Clone(dice)
		resourcesDiscardedByPlayer[playerID] = maps.Clone(resources)
		resourcesDrawnByPlayer[playerID] = maps.Clone(resources)
		resourcesUsedByPlayer[playerID] = maps.Clone(resources)
		devCardsDrawnByPlayer[playerID] = map[string]int{
			"Knight":         0,
			"Victory Point":  0,
			"Road Building":  0,
			"Year of Plenty": 0,
			"Monopoly":       0,
		}
		pointsPerRound[playerID] = make([]int, 0)
		longestRoadEvolutionPerRound[playerID] = make([]int, 0)
		tradesByPlayer[playerID] = map[string]int{
			"TotalStarted":          0,
			"TotalFinalized":        0,
			"ResourceTotalGiven":    0,
			"ResourceTotalReceived": 0,
		}
	}

	return &Instance{
		devCardsDrawnByPlayer:        devCardsDrawnByPlayer,
		dice:                         dice,
		diceByPlayer:                 diceStatsByPlayer,
		longestRoadEvolutionPerRound: longestRoadEvolutionPerRound,
		numberOfRobberiesByPlayer:    numberOfRobberiesByPlayer,
		numberOfTimesRobbedByPlayer:  numberOfTimesRobbedByPlayer,
		pointsEvolutionPerRound:      pointsPerRound,
		resourcesBlockedByPlayer:     resourcesBlockedByPlayer,
		resourcesDiscardedByPlayer:   resourcesDiscardedByPlayer,
		resourcesDrawnByPlayer:       resourcesDrawnByPlayer,
		resourcesUsedByPlayer:        resourcesUsedByPlayer,
		tradesByPlayer:               tradesByPlayer,
	}
}

func (s *Instance) AddDiceEntry(playerID string, sum int) {
	s.dice[sum]++
	s.diceByPlayer[playerID][sum]++
}

func (s *Instance) AddLongestRoadRecord(longestRoads map[string]int) {
	for playerID, score := range longestRoads {
		s.longestRoadEvolutionPerRound[playerID] = append(s.longestRoadEvolutionPerRound[playerID], score)
	}
}

func (s *Instance) AddResourceDiscarded(playerID, resource string, quantity int) {
	s.resourcesDiscardedByPlayer[playerID][resource] += quantity
}

func (s *Instance) AddResourceDrawn(playerID, resource string, quantity int) {
	s.resourcesDrawnByPlayer[playerID][resource] += quantity
}

// TODO: add "blockedBy" to the mix
func (s *Instance) AddResourcesBlocked(playerID, resource string, quantity int) {
	s.resourcesBlockedByPlayer[playerID][resource] += quantity
}

func (s *Instance) AddResourcesUsed(playerID, resource string, quantity int) {
	s.resourcesUsedByPlayer[playerID][resource] += quantity
}

func (s *Instance) AddDevCardDrawn(playerID, devCard string) {
	s.devCardsDrawnByPlayer[playerID][devCard]++
}

func (s *Instance) AddPointsRecord(points map[string]int) {
	for playerID, score := range points {
		s.pointsEvolutionPerRound[playerID] = append(s.pointsEvolutionPerRound[playerID], score)
	}
}

func (s *Instance) AddTradeStarted(playerID string) {
	s.tradesByPlayer[playerID]["TotalStarted"]++
}

func (s *Instance) AddTradeFinalized(playerID string) {
	s.tradesByPlayer[playerID]["TotalFinalized"]++
}

func (s *Instance) AddTradeResourceGiven(playerID string, quantity int) {
	s.tradesByPlayer[playerID]["ResourceTotalGiven"] += quantity
}

func (s *Instance) AddTradeResourceReceived(playerID string, quantity int) {
	s.tradesByPlayer[playerID]["ResourceTotalReceived"] += quantity
}

func (s *Instance) GetDiceHistory() map[int]int {
	return maps.Clone(s.dice)
}

func (s *Instance) GetDiceHistoryByPlayer() map[string]map[int]int {
	return maps.Clone(s.diceByPlayer)
}

func (s *Instance) GetLongestRoadEvolutionPerRound() map[string][]int {
	return maps.Clone(s.longestRoadEvolutionPerRound)
}

func (s *Instance) GetNumberOfRobberiesByPlayer() map[string]int {
	return maps.Clone(s.numberOfRobberiesByPlayer)
}

func (s *Instance) GetNumberOfTimesRobbedByPlayer() map[string]int {
	return maps.Clone(s.numberOfTimesRobbedByPlayer)
}

func (s *Instance) GetResourcesBlockedByPlayer() map[string]map[string]int {
	return maps.Clone(s.resourcesBlockedByPlayer)
}

func (s *Instance) GetResourcesUsedByPlayer() map[string]map[string]int {
	return maps.Clone(s.resourcesUsedByPlayer)
}

func (s *Instance) GetDevCardsDrawnByPlayer() map[string]map[string]int {
	return maps.Clone(s.devCardsDrawnByPlayer)
}

func (s *Instance) GetPointsEvolutionPerRound() map[string][]int {
	return maps.Clone(s.pointsEvolutionPerRound)
}

func (s *Instance) GetTradesByPlayer() map[string]map[string]int {
	return maps.Clone(s.tradesByPlayer)
}
