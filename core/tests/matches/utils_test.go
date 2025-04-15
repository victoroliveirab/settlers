package matches_test

import (
	"fmt"

	testUtils "github.com/victoroliveirab/settlers/core"
)

type GameStateStub struct {
	game                    *testUtils.GameState
	expectedHand            func(expected map[string]string)
	expectedDevHand         func(expected map[string]string)
	expectedSettlements     func(expected map[string][]int)
	expectedCities          func(expected map[string][]int)
	expectedRoads           func(expected map[string][]int)
	expectedLongestRoadSize func(expected map[string]int)
	expectedKnightsUsed     func(expected map[string]int)
	expectedPoints          func(expected map[string]int)
	expectedDice            func(expected int)
}

func createGameStateStub(opts ...testUtils.GameStateOption) *GameStateStub {
	game := testUtils.CreateTestGame(opts...)
	return &GameStateStub{
		game: game,
		expectedHand: func(expected map[string]string) {
			entryOrder := []byte{'L', 'B', 'S', 'G', 'O'}
			entryMap := map[byte]string{
				'L': "Lumber",
				'B': "Brick",
				'S': "Sheep",
				'G': "Grain",
				'O': "Ore",
			}
			players := game.Players()
			for _, player := range players {
				playerID := player.ID
				expectedHand, exists := expected[playerID]
				if !exists {
					panic(fmt.Errorf("not found expected hand for player %s", playerID))
				}
				actualHand := game.ResourceHandByPlayer(playerID)
				index := 0
				quantity := 0
				for _, entry := range entryOrder {
					quantity = 0
					for {
						if index >= len(expectedHand) {
							break
						}
						// fmt.Printf("index:%d,quantity:%d,resource:%s,currentByte:%s\n", index, quantity, string(resource), string(expectedHand[index]))
						if expectedHand[index] == entry {
							quantity++
							index++
						} else {
							break
						}
					}
					entryName := entryMap[entry]
					if actualHand[entryName] != quantity {
						panic(fmt.Errorf("expected player %s to have %d %s, found %d", playerID, quantity, entryName, actualHand[entryName]))
					}
					// fmt.Printf("player:%s, checked %s and has correct quantity, reseting quantity variable\n", playerID, resourceName)
				}
			}
		},
		expectedDevHand: func(expected map[string]string) {
			entryOrder := []byte{'K', 'V', 'R', 'Y', 'M'}
			entryMap := map[byte]string{
				'K': "Knight",
				'V': "Victory Point",
				'R': "Road Building",
				'Y': "Year of Plenty",
				'M': "Monopoly",
			}
			players := game.Players()
			for _, player := range players {
				playerID := player.ID
				expectedDevHand, exists := expected[playerID]
				if !exists {
					panic(fmt.Errorf("not found expected devhand for player %s", playerID))
				}
				actualHand := game.DevelopmentHandByPlayer(playerID)
				index := 0
				quantity := 0
				for _, entry := range entryOrder {
					quantity = 0
					for {
						if index >= len(expectedDevHand) {
							break
						}
						if expectedDevHand[index] == entry {
							quantity++
							index++
						} else {
							break
						}
					}
					entryName := entryMap[entry]
					if actualHand[entryName] != quantity {
						panic(fmt.Errorf("expected player %s to have %d %s, found %d", playerID, quantity, entryName, actualHand[entryName]))
					}
				}
			}
		},
		expectedSettlements: func(expected map[string][]int) {
			players := game.Players()
			for _, player := range players {
				playerID := player.ID
				expectedSettlements, exists := expected[playerID]
				if !exists {
					panic(fmt.Errorf("not found expected settlements for player %s", playerID))
				}
				actualSettlements := game.SettlementsByPlayer(playerID)
				if len(expectedSettlements) != len(actualSettlements) {
					panic(fmt.Errorf("expected settlements to be %v, got %v", expectedSettlements, actualSettlements))
				}
				for i, vertexID := range actualSettlements {
					if expectedSettlements[i] != vertexID {
						panic(fmt.Errorf("expected to have settlement#%d, but doesn't", expectedSettlements[i]))
					}
				}
			}
		},
		expectedCities: func(expected map[string][]int) {
			players := game.Players()
			for _, player := range players {
				playerID := player.ID
				expectedCities, exists := expected[playerID]
				if !exists {
					panic(fmt.Errorf("not found expected cities for player %s", playerID))
				}
				actualCities := game.CitiesByPlayer(playerID)
				if len(expectedCities) != len(actualCities) {
					panic(fmt.Errorf("expected cities to be %v, got %v", expectedCities, actualCities))
				}
				for i, vertexID := range actualCities {
					if expectedCities[i] != vertexID {
						panic(fmt.Errorf("expected to have city#%d, but doesn't", expectedCities[i]))
					}
				}
			}
		},
		expectedRoads: func(expected map[string][]int) {
			players := game.Players()
			for _, player := range players {
				playerID := player.ID
				expectedRoads, exists := expected[playerID]
				if !exists {
					panic(fmt.Errorf("not found expected roads for player %s", playerID))
				}
				actualRoads := game.RoadsByPlayer(playerID)
				if len(expectedRoads) != len(actualRoads) {
					panic(fmt.Errorf("expected roads to be %v, got %v", expectedRoads, actualRoads))
				}
				for i, vertexID := range actualRoads {
					if expectedRoads[i] != vertexID {
						panic(fmt.Errorf("expected to have road#%d, but doesn't", expectedRoads[i]))
					}
				}
			}
		},
		expectedLongestRoadSize: func(expected map[string]int) {
			players := game.Players()
			for _, player := range players {
				playerID := player.ID
				expectedLongestRoadSize, exists := expected[playerID]
				if !exists {
					panic(fmt.Errorf("not found expected longest road size for player %s", playerID))
				}
				actualLongestRoadSize := game.LongestRoadLengthByPlayer(playerID)
				if expectedLongestRoadSize != actualLongestRoadSize {
					panic(fmt.Errorf("expected %s to have %d longest road length, actually has %d", playerID, expectedLongestRoadSize, actualLongestRoadSize))
				}
			}
		},
		expectedKnightsUsed: func(expected map[string]int) {
			players := game.Players()
			for _, player := range players {
				playerID := player.ID
				expectedKnightsUsed, exists := expected[playerID]
				if !exists {
					panic(fmt.Errorf("not found expected number of knights used for player %s", playerID))
				}
				actualKnightsUsed := game.NumberOfKnightsUsedByPlayer(playerID)
				if expectedKnightsUsed != actualKnightsUsed {
					panic(fmt.Errorf("expected %s to have %d knights used, actually has %d", playerID, expectedKnightsUsed, actualKnightsUsed))
				}
			}
		},
		expectedPoints: func(expected map[string]int) {
			players := game.Players()
			points := game.Points()
			for _, player := range players {
				playerID := player.ID
				expectedPoints, exists := expected[playerID]
				if !exists {
					panic(fmt.Errorf("not found expected points for player %s", playerID))
				}
				actualPoints := points[playerID]
				if expectedPoints != actualPoints {
					panic(fmt.Errorf("expected %s to have %d points, actually has %d", playerID, expectedPoints, actualPoints))
				}
			}
		},
		expectedDice: func(expected int) {
			dice := game.Dice()
			if expected != dice[0]+dice[1] {
				panic(fmt.Errorf("expected dice to be %d, but actually got %d + %d", expected, dice[0], dice[1]))
			}
		},
	}
}
