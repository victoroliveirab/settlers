package core

import (
	"github.com/victoroliveirab/settlers/core/packages/round"
	"github.com/victoroliveirab/settlers/core/packages/summary"
)

func (state *GameState) GetStatistics() summary.Statistics {
	report := state.summary.GetReport(summary.ReportInput{
		LargestArmyOwner: state.mostKnights.PlayerID,
		LongestRoadOwner: state.longestRoad.PlayerID,
		Points:           state.Points(),
	})
	stats := report.Statistics
	if state.round.GetRoundType() != round.GameOver {
		stats.PointsEvolution = nil
	}
	return stats
}
