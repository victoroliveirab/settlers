package summary

func (s *Instance) getPlayerPointDistribution(input ReportInput) map[string]PlayerPointDistribution {
	pointsDistribution := make(map[string]PlayerPointDistribution)
	for playerID, player := range s.playersStates {
		playerPointDistribution := PlayerPointDistribution{}
		playerPointDistribution.Total = input.Points[playerID]
		playerPointDistribution.Settlements = player.GetNumberOfSettlements() * s.settings.PointsPerSettlement
		playerPointDistribution.Cities = player.GetNumberOfCities() * s.settings.PointsPerCity
		playerPointDistribution.VictoryPoints = player.GetNumberOfVictoryPoints()
		playerPointDistribution.LongestRoad = 0
		if playerID == input.LongestRoadOwner {
			playerPointDistribution.LongestRoad = s.settings.PointsForLongestRoad
		}
		playerPointDistribution.LargestArmy = 0
		if playerID == input.LargestArmyOwner {
			playerPointDistribution.LargestArmy = s.settings.PointsForMostKnights
		}
	}
	return pointsDistribution
}
