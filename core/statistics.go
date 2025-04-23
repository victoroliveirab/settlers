package core

import "github.com/victoroliveirab/settlers/core/packages/round"

type Statistics struct {
	GeneralDiceStats          map[int]int            `json:"generalDiceStats"`
	DiceStatsByPlayer         map[string]map[int]int `json:"diceStatsByPlayer"`
	LongestRoadEvolution      map[string][]int       `json:"longestRoadEvolution"`
	NumberOfRobberiesByPlayer map[string]int         `json:"numberOfRobberiesByPlayer"`
	PointsEvolution           map[string][]int       `json:"pointsEvolution"`
}

func (state *GameState) GetStatistics() Statistics {
	if state.round.GetRoundType() != round.GameOver {
		return Statistics{
			GeneralDiceStats:          state.stats.GetDiceHistory(),
			DiceStatsByPlayer:         state.stats.GetDiceHistoryByPlayer(),
			LongestRoadEvolution:      state.stats.GetLongestRoadEvolutionPerRound(),
			NumberOfRobberiesByPlayer: state.stats.GetNumberOfRobberiesByPlayer(),
		}
	}
	return Statistics{
		GeneralDiceStats:          state.stats.GetDiceHistory(),
		DiceStatsByPlayer:         state.stats.GetDiceHistoryByPlayer(),
		LongestRoadEvolution:      state.stats.GetLongestRoadEvolutionPerRound(),
		NumberOfRobberiesByPlayer: state.stats.GetNumberOfRobberiesByPlayer(),
		PointsEvolution:           state.stats.GetPointsEvolutionPerRound(),
	}
}
