package core

import (
	"maps"
)

func (state *GameState) recountLongestRoad() bool {
	var longestRoad LongestRoad
	for _, player := range state.players {
		playerID := player.ID
		playerState := state.playersStates[playerID]
		playerLongestRoad := playerState.LongestRoadSegments
		playerLongestRoadSize := len(playerLongestRoad)
		// REFACTOR: this looks waaay more convoluted than I think it needs to be
		// By the time I'm writing this, this is the best I came up with
		if playerLongestRoadSize > longestRoad.Length {
			longestRoad = LongestRoad{
				PlayerID: playerID,
				Length:   playerLongestRoadSize,
			}
		} else if playerLongestRoadSize == longestRoad.Length && playerID == state.longestRoad.PlayerID {
			longestRoad = LongestRoad{
				PlayerID: playerID,
				Length:   playerLongestRoadSize,
			}
		}
	}

	// NOTE: this has the potential to reset longest road if a settlement blocks a previous largest
	// road of 5. I don't know what should be done rulewise, so I chose to revoke the title
	if longestRoad.Length < state.longestRoadMinimum {
		changed := state.longestRoad.PlayerID != ""
		state.longestRoad = LongestRoad{
			PlayerID: "",
			Length:   0,
		}
		return changed
	}

	if longestRoad.PlayerID != state.longestRoad.PlayerID {
		state.longestRoad.PlayerID = longestRoad.PlayerID
		return true
	}
	return false
}

func (state *GameState) recountKnights() bool {
	changed := false
	for _, player := range state.players {
		playerState := state.playersStates[player.ID]
		knightsUsed := playerState.UsedDevelopmentCards["Knight"]
		if knightsUsed > state.mostKnights.Quantity && knightsUsed >= state.mostKnightsMinimum {
			state.mostKnights.PlayerID = player.ID
			state.mostKnights.Quantity = knightsUsed
			changed = true
		}
	}
	return changed
}

func (state *GameState) updatePoints() {
	var victoryPlayer string
	for _, player := range state.players {
		playerID := player.ID
		playerState := state.playersStates[playerID]
		sum := 0
		sum += state.pointsPerSettlement * len(playerState.Settlements)
		sum += state.pointsPerCity * len(playerState.Cities)
		sum += len(playerState.DevelopmentCards["Victory Point"])

		if state.mostKnights.PlayerID == playerID {
			sum += state.pointsPerMostKnights
		}
		if state.longestRoad.PlayerID == playerID {
			sum += state.pointsPerLongestRoad
		}
		state.points[playerID] = sum
		if sum >= state.targetPoint {
			victoryPlayer = playerID
		}
	}
	if victoryPlayer != "" {
		state.roundType = GameOver
	}
}

func (state *GameState) Points() map[string]int {
	return maps.Clone(state.points)
}

func (state *GameState) PublicPoints() map[string]int {
	points := state.Points()
	for player := range points {
		numberOfVictoryPoints := state.DevelopmentHandByPlayer(player)["Victory Point"]
		points[player] -= numberOfVictoryPoints
	}
	return points
}
