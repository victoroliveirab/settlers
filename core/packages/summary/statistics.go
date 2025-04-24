package summary

func (s *Instance) getStatistics(input ReportInput) Statistics {
	return Statistics{
		GeneralDiceStats:          s.bookKeeping.GetDiceHistory(),
		DiceStatsByPlayer:         s.bookKeeping.GetDiceHistoryByPlayer(),
		LongestRoadEvolution:      s.bookKeeping.GetLongestRoadEvolutionPerRound(),
		NumberOfRobberiesByPlayer: s.bookKeeping.GetNumberOfRobberiesByPlayer(),
		PointsEvolution:           s.bookKeeping.GetPointsEvolutionPerRound(),
	}
}
