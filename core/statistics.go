package core

import (
	"github.com/victoroliveirab/settlers/core/packages/round"
	"github.com/victoroliveirab/settlers/core/packages/summary"
)

func (state *GameState) GetReport() summary.ReportOutput {
	report := state.summary.GetReport(summary.ReportInput{
		LargestArmyOwner: state.mostKnights.PlayerID,
		LongestRoadOwner: state.longestRoad.PlayerID,
		Points:           state.Points(),
	})
	if state.round.GetRoundType() != round.GameOver {
		report.Statistics.PointsEvolution = nil
		report.PointsDistribution = nil
	}
	return report
}
