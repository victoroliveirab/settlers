package summary

import (
	"github.com/victoroliveirab/settlers/core/packages/book-keeping"
	"github.com/victoroliveirab/settlers/core/packages/player"
	coreT "github.com/victoroliveirab/settlers/core/types"
)

type PlayerPointDistribution struct {
	Total         int `json:"total"`
	Settlements   int `json:"settlements"`
	Cities        int `json:"cities"`
	VictoryPoints int `json:"victoryPoints"`
	LargestArmy   int `json:"largestArmy"`
	LongestRoad   int `json:"longestRoad"`
}

type Statistics struct {
	GeneralDiceStats          map[int]int            `json:"generalDiceStats"`
	DiceStatsByPlayer         map[string]map[int]int `json:"diceStatsByPlayer"`
	LongestRoadEvolution      map[string][]int       `json:"longestRoadEvolution"`
	NumberOfRobberiesByPlayer map[string]int         `json:"numberOfRobberiesByPlayer"`
	PointsEvolution           map[string][]int       `json:"pointsEvolution"`
}

type ReportInput struct {
	LargestArmyOwner string
	LongestRoadOwner string
	Points           map[string]int
}

type ReportOutput struct {
	PointsDistribution map[string]PlayerPointDistribution `json:"pointsDistribution"`
	Statistics         Statistics                         `json:"statistics"`
}

type Instance struct {
	playersStates map[string]*player.Instance
	settings      coreT.Settings
	bookKeeping   *bookkeeping.Instance
}

func New(
	players map[string]*player.Instance,
	settings coreT.Settings,
	bookKeeping *bookkeeping.Instance,
) *Instance {
	return &Instance{
		bookKeeping:   bookKeeping,
		playersStates: players,
		settings:      settings,
	}
}

func (s *Instance) GetReport(input ReportInput) ReportOutput {
	pointsDistribution := s.getPlayerPointDistribution(input)
	statistics := s.getStatistics(input)
	return ReportOutput{
		PointsDistribution: pointsDistribution,
		Statistics:         statistics,
	}
}
