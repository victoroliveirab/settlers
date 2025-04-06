package tests

import (
	"fmt"
	"testing"

	testUtils "github.com/victoroliveirab/settlers/core"
)

func TestFullGame(t *testing.T) {
	game := testUtils.CreateTestGame()

	expectedHand := func(expected map[string]string) {
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
	}
	expectedDevHand := func(expected map[string]string) {
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
	}
	expectedSettlements := func(expected map[string][]int) {
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
	}
	expectedCities := func(expected map[string][]int) {
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
	}
	expectedRoads := func(expected map[string][]int) {
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
	}
	expectedLongestRoadSize := func(expected map[string]int) {
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
	}
	expectedKnightsUsed := func(expected map[string]int) {
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
	}
	expectedPoints := func(expected map[string]int) {
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
	}

	t.Run("full run", func(t *testing.T) {
		game.BuildSettlement("1", 2)
		game.BuildRoad("1", 2)

		game.BuildSettlement("2", 28)
		game.BuildRoad("2", 37)

		game.BuildSettlement("3", 45)
		game.BuildRoad("3", 60)

		game.BuildSettlement("4", 22)
		game.BuildRoad("4", 28)

		game.BuildSettlement("4", 24)
		game.BuildRoad("4", 29)

		game.BuildSettlement("3", 8)
		game.BuildRoad("3", 9)

		game.BuildSettlement("2", 33)
		game.BuildRoad("2", 43)

		game.BuildSettlement("1", 4)
		game.BuildRoad("1", 3)

		// End of setup phase
		if game.RoundType() != testUtils.FirstRound {
			t.Errorf("expected to be at first round after setup end, round type is actually %s", testUtils.RoundTypeTranslation[game.RoundType()])
		}

		expectedHand(map[string]string{
			"1": "SO",
			"2": "LBO",
			"3": "LSG",
			"4": "GO",
		})
		expectedRoads(map[string][]int{
			"1": {2, 3},
			"2": {37, 43},
			"3": {60, 9},
			"4": {28, 29},
		})
		expectedSettlements(map[string][]int{
			"1": {2, 4},
			"2": {28, 33},
			"3": {45, 8},
			"4": {22, 24},
		})
		expectedCities(map[string][]int{
			"1": {},
			"2": {},
			"3": {},
			"4": {},
		})
		expectedLongestRoadSize(map[string]int{
			"1": 2,
			"2": 1,
			"3": 1,
			"4": 2,
		})
		expectedKnightsUsed(map[string]int{
			"1": 0,
			"2": 0,
			"3": 0,
			"4": 0,
		})
		expectedPoints(map[string]int{
			"1": 2,
			"2": 2,
			"3": 2,
			"4": 2,
		})

		// Start of the game
		game.RollDice("1") // 8
		game.EndRound("1")
		expectedHand(map[string]string{
			"1": "SO",
			"2": "LBGOO",
			"3": "LSG",
			"4": "GOO",
		})

		game.RollDice("2") // 5
		game.BuildRoad("2", 44)
		game.EndRound("2")
		expectedHand(map[string]string{
			"1": "SO",
			"2": "GGOO",
			"3": "LSGG",
			"4": "GGOO",
		})
		expectedRoads(map[string][]int{
			"1": {2, 3},
			"2": {37, 43, 44},
			"3": {60, 9},
			"4": {28, 29},
		})
		expectedLongestRoadSize(map[string]int{
			"1": 2,
			"2": 2,
			"3": 1,
			"4": 2,
		})
		expectedPoints(map[string]int{
			"1": 2,
			"2": 2,
			"3": 2,
			"4": 2,
		})

		game.RollDice("3") // 11
		game.EndRound("3")
		expectedHand(map[string]string{
			"1": "SOOO",
			"2": "GGOO",
			"3": "LLSGG",
			"4": "GGOO",
		})

		game.RollDice("4") // 9
		tradeID, err := game.MakeTradeOffer("4", map[string]int{"Grain": 1}, map[string]int{"Sheep": 1}, []string{})
		if err != nil {
			t.Error(err)
			return
		}
		counterTradeID, _ := game.MakeCounterTradeOffer("1", tradeID, map[string]int{"Sheep": 1}, map[string]int{"Grain": 1, "Ore": 1})
		game.RejectTradeOffer("2", tradeID)
		game.AcceptTradeOffer("3", counterTradeID)
		err = game.FinalizeTrade("4", "1", counterTradeID)
		if err != nil {
			t.Error(err)
			return
		}
		err = game.BuyDevelopmentCard("4")
		if err != nil {
			t.Error(err)
			return
		}
		game.EndRound("4")
		expectedHand(map[string]string{
			"1": "GOOOO",
			"2": "GGOO",
			"3": "LLLSGG",
			"4": "",
		})
		expectedDevHand(map[string]string{
			"1": "",
			"2": "",
			"3": "",
			"4": "K",
		})
		expectedPoints(map[string]int{
			"1": 2,
			"2": 2,
			"3": 2,
			"4": 2,
		})

		game.RollDice("1") // 10
		game.EndRound("1")
		expectedHand(map[string]string{
			"1": "GOOOO",
			"2": "BGGOO",
			"3": "LLLBSGG",
			"4": "O",
		})
		expectedDevHand(map[string]string{
			"1": "",
			"2": "",
			"3": "",
			"4": "K",
		})

		game.RollDice("2") // 5
		game.MakeTradeOffer("2", map[string]int{"Brick": 1}, map[string]int{"Sheep": 1}, []string{})
		tradeID, _ = game.MakeTradeOffer("2", map[string]int{"Brick": 1}, map[string]int{"Ore": 1}, []string{})
		game.AcceptTradeOffer("1", tradeID)
		game.MakeCounterTradeOffer("4", tradeID, map[string]int{"Sheep": 1}, map[string]int{"Lumber": 1, "Brick": 1})
		game.RejectTradeOffer("3", tradeID)
		game.FinalizeTrade("2", "1", tradeID)
		err = game.BuildCity("2", 28)
		if err != nil {
			t.Error(err)
			return
		}
		game.EndRound("2")
		expectedHand(map[string]string{
			"1": "BGOOO",
			"2": "G",
			"3": "LLLBSGGG",
			"4": "GO",
		})
		expectedDevHand(map[string]string{
			"1": "",
			"2": "",
			"3": "",
			"4": "K",
		})
		expectedSettlements(map[string][]int{
			"1": {2, 4},
			"2": {33},
			"3": {45, 8},
			"4": {22, 24},
		})
		expectedCities(map[string][]int{
			"1": {},
			"2": {28},
			"3": {},
			"4": {},
		})
		expectedPoints(map[string]int{
			"1": 2,
			"2": 3,
			"3": 2,
			"4": 2,
		})

		game.RollDice("3") // 8
		game.BuildRoad("3", 72)
		tradeID, _ = game.MakeTradeOffer("3", map[string]int{"Lumber": 1, "Grain": 1}, map[string]int{"Brick": 1}, []string{})
		game.AcceptTradeOffer("1", tradeID)
		game.FinalizeTrade("3", "1", tradeID)
		err = game.BuildSettlement("3", 54)
		if err != nil {
			t.Error(err)
			return
		}
		game.EndRound("3")
		expectedHand(map[string]string{
			"1": "LGGOOO",
			"2": "GGGO",
			"3": "G",
			"4": "GOO",
		})
		expectedDevHand(map[string]string{
			"1": "",
			"2": "",
			"3": "",
			"4": "K",
		})
		expectedSettlements(map[string][]int{
			"1": {2, 4},
			"2": {33},
			"3": {45, 8, 54},
			"4": {22, 24},
		})
		expectedRoads(map[string][]int{
			"1": {2, 3},
			"2": {37, 43, 44},
			"3": {60, 9, 72},
			"4": {28, 29},
		})
		expectedLongestRoadSize(map[string]int{
			"1": 2,
			"2": 2,
			"3": 2,
			"4": 2,
		})
		expectedPoints(map[string]int{
			"1": 2,
			"2": 3,
			"3": 3,
			"4": 2,
		})

		err = game.UseKnight("4")
		if err != nil {
			t.Error(err)
			return
		}
		err = game.MoveRobber("4", 6)
		if err != nil {
			t.Error(err)
			return
		}
		game.RobPlayer("4", "3")
		expectedHand(map[string]string{
			"1": "LGGOOO",
			"2": "GGGO",
			"3": "",
			"4": "GGOO",
		})
		expectedDevHand(map[string]string{
			"1": "",
			"2": "",
			"3": "",
			"4": "",
		})
		expectedKnightsUsed(map[string]int{
			"1": 0,
			"2": 0,
			"3": 0,
			"4": 1,
		})
		expectedPoints(map[string]int{
			"1": 2,
			"2": 3,
			"3": 3,
			"4": 2,
		})

		game.RollDice("4") // 7
		err = game.MoveRobber("4", 1)
		if err != nil {
			t.Error(err)
			return
		}
		game.RobPlayer("4", "1")
		game.EndRound("4")
		expectedHand(map[string]string{
			"1": "GGOOO",
			"2": "GGGO",
			"3": "",
			"4": "LGGOO",
		})

		game.RollDice("1") // 2
		err = game.BuildCity("1", 4)
		if err != nil {
			t.Error(err)
			return
		}
		game.EndRound("1")
		expectedHand(map[string]string{
			"1": "",
			"2": "GGGO",
			"3": "",
			"4": "LGGOO",
		})
		expectedSettlements(map[string][]int{
			"1": {2},
			"2": {33},
			"3": {45, 8, 54},
			"4": {22, 24},
		})
		expectedCities(map[string][]int{
			"1": {4},
			"2": {28},
			"3": {},
			"4": {},
		})
		expectedPoints(map[string]int{
			"1": 3,
			"2": 3,
			"3": 3,
			"4": 2,
		})

		game.RollDice("2") // 2
		game.EndRound("2")

		game.RollDice("3") // 9
		game.EndRound("3")
		expectedHand(map[string]string{
			"1": "",
			"2": "GGGO",
			"3": "L",
			"4": "LGGOO",
		})

		game.RollDice("4") // 4
		game.EndRound("4")
		expectedHand(map[string]string{
			"1": "",
			"2": "LGGGO",
			"3": "L",
			"4": "LGGOO",
		})
		expectedDevHand(map[string]string{
			"1": "",
			"2": "",
			"3": "",
			"4": "",
		})

		game.RollDice("1") // 7
		game.MoveRobber("1", 15)
		game.RobPlayer("1", "3")
		game.EndRound("1")
		expectedHand(map[string]string{
			"1": "L",
			"2": "LGGGO",
			"3": "",
			"4": "LGGOO",
		})

		game.RollDice("2") // 5
		game.MakeBankTrade("2", "Grain", "Sheep")
		game.BuyDevelopmentCard("2")
		game.EndRound("2")
		expectedHand(map[string]string{
			"1": "L",
			"2": "L",
			"3": "G",
			"4": "LGGGOO",
		})
		expectedDevHand(map[string]string{
			"1": "",
			"2": "Y",
			"3": "",
			"4": "",
		})

		game.RollDice("3") // 10
		game.EndRound("3")
		expectedHand(map[string]string{
			"1": "L",
			"2": "L",
			"3": "G",
			"4": "LGGGOOO",
		})

		game.RollDice("4") // 9
		game.BuildCity("4", 22)
		game.EndRound("4")
		expectedHand(map[string]string{
			"1": "L",
			"2": "L",
			"3": "LG",
			"4": "LG",
		})
		expectedSettlements(map[string][]int{
			"1": {2},
			"2": {33},
			"3": {45, 8, 54},
			"4": {24},
		})
		expectedCities(map[string][]int{
			"1": {4},
			"2": {28},
			"3": {},
			"4": {22},
		})
		expectedPoints(map[string]int{
			"1": 3,
			"2": 3,
			"3": 3,
			"4": 3,
		})

		game.RollDice("1") // 6
		game.EndRound("1")
		expectedHand(map[string]string{
			"1": "L",
			"2": "LSS",
			"3": "LLLG",
			"4": "LG",
		})

		game.RollDice("2") // 7
		game.MoveRobber("2", 19)
		game.RobPlayer("2", "3")
		game.UseYearOfPlenty("2")
		game.PickYearOfPlentyResources("2", "Brick", "Grain")
		game.BuildSettlement("2", 35)
		game.EndRound("2")
		expectedHand(map[string]string{
			"1": "L",
			"2": "LS",
			"3": "LLG",
			"4": "LG",
		})
		expectedSettlements(map[string][]int{
			"1": {2},
			"2": {33, 35},
			"3": {45, 8, 54},
			"4": {24},
		})
		expectedCities(map[string][]int{
			"1": {4},
			"2": {28},
			"3": {},
			"4": {22},
		})
		expectedPoints(map[string]int{
			"1": 3,
			"2": 4,
			"3": 3,
			"4": 3,
		})
		expectedDevHand(map[string]string{
			"1": "",
			"2": "",
			"3": "",
			"4": "",
		})

		game.RollDice("3") // 4
		game.EndRound("3")
		expectedHand(map[string]string{
			"1": "L",
			"2": "LLS",
			"3": "LLG",
			"4": "LG",
		})

		game.RollDice("4") // 3
		game.MakeBankTrade("4", "Grain", "Brick")
		game.BuildRoad("4", 26)
		game.EndRound("4")
		expectedHand(map[string]string{
			"1": "LSS",
			"2": "LLS",
			"3": "LLSG",
			"4": "",
		})
		expectedRoads(map[string][]int{
			"1": {2, 3},
			"2": {37, 43, 44},
			"3": {60, 9, 72},
			"4": {28, 29, 26},
		})
		expectedLongestRoadSize(map[string]int{
			"1": 2,
			"2": 2,
			"3": 2,
			"4": 3,
		})

		game.RollDice("1") // 10
		game.EndRound("1")
		expectedHand(map[string]string{
			"1": "LSS",
			"2": "LLBSO",
			"3": "LLBSG",
			"4": "O",
		})

		game.RollDice("2") // 7
		game.MoveRobber("2", 3)
		game.RobPlayer("2", "3")
		game.BuildRoad("2", 52)
		game.EndRound("2")
		expectedHand(map[string]string{
			"1": "LSS",
			"2": "LSSO",
			"3": "LLBG",
			"4": "O",
		})
		expectedRoads(map[string][]int{
			"1": {2, 3},
			"2": {37, 43, 44, 52},
			"3": {60, 9, 72},
			"4": {28, 29, 26},
		})
		expectedLongestRoadSize(map[string]int{
			"1": 2,
			"2": 2,
			"3": 2,
			"4": 3,
		})

		game.RollDice("3") // 5
		game.BuildRoad("3", 16)
		game.EndRound("3")
		expectedHand(map[string]string{
			"1": "LSS",
			"2": "LSSGGO",
			"3": "LGG",
			"4": "GGO",
		})
		expectedRoads(map[string][]int{
			"1": {2, 3},
			"2": {37, 43, 44, 52},
			"3": {60, 9, 72, 16},
			"4": {28, 29, 26},
		})
		expectedLongestRoadSize(map[string]int{
			"1": 2,
			"2": 2,
			"3": 2,
			"4": 3,
		})

		game.RollDice("4") // 9
		tradeID, _ = game.MakeTradeOffer("4", map[string]int{"Grain": 1}, map[string]int{"Sheep": 1}, []string{})
		counterTradeID, _ = game.MakeCounterTradeOffer("1", tradeID, map[string]int{"Sheep": 1}, map[string]int{"Grain": 1, "Ore": 1})
		game.RejectTradeOffer("4", counterTradeID)
		game.AcceptTradeOffer("1", tradeID)
		game.FinalizeTrade("4", "1", tradeID)
		game.BuyDevelopmentCard("4")
		game.EndRound("4")
		expectedHand(map[string]string{
			"1": "LSG",
			"2": "LSSGGO",
			"3": "LGG",
			"4": "",
		})
		expectedDevHand(map[string]string{
			"1": "",
			"2": "",
			"3": "",
			"4": "K",
		})

		game.RollDice("1") // 11
		game.BuyDevelopmentCard("1")
		game.MakePortTrade("1", 4, "Ore", "Brick")
		game.BuildRoad("1", 19)
		game.EndRound("1")
		expectedHand(map[string]string{
			"1": "",
			"2": "LLSSGGO",
			"3": "LLGG",
			"4": "",
		})
		expectedDevHand(map[string]string{
			"1": "K",
			"2": "",
			"3": "",
			"4": "K",
		})
		expectedRoads(map[string][]int{
			"1": {2, 3, 19},
			"2": {37, 43, 44, 52},
			"3": {60, 9, 72, 16},
			"4": {28, 29, 26},
		})
		expectedLongestRoadSize(map[string]int{
			"1": 3,
			"2": 2,
			"3": 2,
			"4": 3,
		})

		game.RollDice("2") // 3
		game.BuyDevelopmentCard("2")
		game.EndRound("2")
		expectedHand(map[string]string{
			"1": "SS",
			"2": "LLSG",
			"3": "LLSGG",
			"4": "GGG",
		})
		expectedDevHand(map[string]string{
			"1": "K",
			"2": "Y",
			"3": "",
			"4": "K",
		})

		game.RollDice("3") // 12
		game.MakePortTrade("3", 54, "Lumber", "Ore")
		game.BuyDevelopmentCard("3")
		err = game.UseKnight("3")
		if err == nil {
			t.Error("Shouldn't be able to use knight the round it bought")
			return
		}
		game.EndRound("3")
		expectedHand(map[string]string{
			"1": "BSS",
			"2": "LLSG",
			"3": "G",
			"4": "GGG",
		})
		expectedDevHand(map[string]string{
			"1": "K",
			"2": "Y",
			"3": "K",
			"4": "K",
		})
		expectedKnightsUsed(map[string]int{
			"1": 0,
			"2": 0,
			"3": 0,
			"4": 1,
		})

		game.UseKnight("4")
		game.MoveRobber("4", 19)
		game.RobPlayer("4", "3")
		game.RollDice("4") // 4
		game.EndRound("4")
		expectedHand(map[string]string{
			"1": "BSS",
			"2": "LLLSG",
			"3": "",
			"4": "GGGG",
		})
		expectedDevHand(map[string]string{
			"1": "K",
			"2": "Y",
			"3": "K",
			"4": "",
		})
		expectedKnightsUsed(map[string]int{
			"1": 0,
			"2": 0,
			"3": 0,
			"4": 2,
		})
		expectedPoints(map[string]int{
			"1": 3,
			"2": 4,
			"3": 3,
			"4": 3,
		})

		game.RollDice("1") // 12
		tradeID, _ = game.MakeTradeOffer("1", map[string]int{"Sheep": 1}, map[string]int{"Lumber": 1}, []string{})
		game.AcceptTradeOffer("2", tradeID)
		game.FinalizeTrade("1", "2", tradeID)
		game.BuildRoad("1", 22)
		game.EndRound("1")
		expectedHand(map[string]string{
			"1": "BS",
			"2": "LLSSG",
			"3": "",
			"4": "GGGG",
		})
		expectedRoads(map[string][]int{
			"1": {2, 3, 19, 22},
			"2": {37, 43, 44, 52},
			"3": {60, 9, 72, 16},
			"4": {28, 29, 26},
		})
		expectedLongestRoadSize(map[string]int{
			"1": 4,
			"2": 2,
			"3": 2,
			"4": 3,
		})

		expectedPoints(map[string]int{
			"1": 3,
			"2": 4,
			"3": 3,
			"4": 3,
		})

		game.RollDice("2") // 2
		game.UseYearOfPlenty("2")
		game.PickYearOfPlentyResources("2", "Brick", "Brick")
		game.BuildRoad("2", 51)
		game.BuildSettlement("2", 41)
		game.EndRound("2")
		expectedHand(map[string]string{
			"1": "BS",
			"2": "S",
			"3": "",
			"4": "GGGG",
		})
		expectedDevHand(map[string]string{
			"1": "K",
			"2": "",
			"3": "K",
			"4": "",
		})
		expectedRoads(map[string][]int{
			"1": {2, 3, 19, 22},
			"2": {37, 43, 44, 52, 51},
			"3": {60, 9, 72, 16},
			"4": {28, 29, 26},
		})
		expectedLongestRoadSize(map[string]int{
			"1": 4,
			"2": 3,
			"3": 2,
			"4": 3,
		})
		expectedSettlements(map[string][]int{
			"1": {2},
			"2": {33, 35, 41},
			"3": {45, 8, 54},
			"4": {24},
		})
		expectedPoints(map[string]int{
			"1": 3,
			"2": 5,
			"3": 3,
			"4": 3,
		})

		game.UseKnight("3")
		game.MoveRobber("3", 11)
		game.RobPlayer("3", "2")
		game.RollDice("3") // 5
		game.EndRound("3")
		expectedHand(map[string]string{
			"1": "BS",
			"2": "GG",
			"3": "SG",
			"4": "GGGGGG",
		})
		expectedKnightsUsed(map[string]int{
			"1": 0,
			"2": 0,
			"3": 1,
			"4": 2,
		})
		expectedDevHand(map[string]string{
			"1": "K",
			"2": "",
			"3": "",
			"4": "",
		})

		game.RollDice("4") // 7
		game.MoveRobber("4", 13)
		game.RobPlayer("4", "2")
		tradeID, _ = game.MakeTradeOffer("4", map[string]int{"Grain": 1}, map[string]int{"Sheep": 1}, []string{})
		counterTradeID, _ = game.MakeCounterTradeOffer("3", tradeID, map[string]int{"Sheep": 1}, map[string]int{"Grain": 2})
		game.FinalizeTrade("4", "3", counterTradeID)
		game.MakeBankTrade("4", "Grain", "Ore")
		game.BuyDevelopmentCard("4")
		game.EndRound("4")
		expectedHand(map[string]string{
			"1": "BS",
			"2": "G",
			"3": "GGG",
			"4": "",
		})
		expectedKnightsUsed(map[string]int{
			"1": 0,
			"2": 0,
			"3": 1,
			"4": 2,
		})
		expectedDevHand(map[string]string{
			"1": "K",
			"2": "",
			"3": "",
			"4": "K",
		})

		game.RollDice("1") // 4
		game.EndRound("1")
		expectedHand(map[string]string{
			"1": "BS",
			"2": "LBG",
			"3": "GGG",
			"4": "",
		})

		game.RollDice("2") // 6
		game.BuildRoad("2", 50)
		game.EndRound("2")
		expectedHand(map[string]string{
			"1": "BS",
			"2": "G",
			"3": "LLGGG",
			"4": "",
		})
		expectedRoads(map[string][]int{
			"1": {2, 3, 19, 22},
			"2": {37, 43, 44, 52, 51, 50},
			"3": {60, 9, 72, 16},
			"4": {28, 29, 26},
		})
		expectedLongestRoadSize(map[string]int{
			"1": 4,
			"2": 4,
			"3": 2,
			"4": 3,
		})

		game.RollDice("3") // 10
		tradeID, _ = game.MakeTradeOffer("3", map[string]int{"Grain": 1}, map[string]int{"Sheep": 1}, []string{})
		counterTradeID, _ = game.MakeCounterTradeOffer("1", tradeID, map[string]int{"Sheep": 1}, map[string]int{"Grain": 1, "Lumber": 1})
		err = game.FinalizeTrade("3", "1", counterTradeID)
		game.BuildSettlement("3", 14)
		game.EndRound("3")
		expectedHand(map[string]string{
			"1": "LBG",
			"2": "BGO",
			"3": "G",
			"4": "O",
		})
		expectedSettlements(map[string][]int{
			"1": {2},
			"2": {33, 35, 41},
			"3": {45, 8, 54, 14},
			"4": {24},
		})
		expectedPoints(map[string]int{
			"1": 3,
			"2": 5,
			"3": 4,
			"4": 3,
		})

		game.UseKnight("4")
		game.MoveRobber("4", 19)
		game.RobPlayer("4", "3")
		game.RollDice("4") // 3
		game.EndRound("4")
		expectedHand(map[string]string{
			"1": "LBSSG",
			"2": "BGO",
			"3": "S",
			"4": "GGGGO",
		})
		expectedPoints(map[string]int{
			"1": 3,
			"2": 5,
			"3": 4,
			"4": 5,
		})
		expectedKnightsUsed(map[string]int{
			"1": 0,
			"2": 0,
			"3": 1,
			"4": 3,
		})

		game.RollDice("1") // 4
		game.BuildSettlement("1", 19)
		game.EndRound("1")
		expectedHand(map[string]string{
			"1": "S",
			"2": "LBBGO",
			"3": "S",
			"4": "GGGGO",
		})
		expectedSettlements(map[string][]int{
			"1": {2, 19},
			"2": {33, 35, 41},
			"3": {45, 8, 54, 14},
			"4": {24},
		})
		expectedPoints(map[string]int{
			"1": 4,
			"2": 5,
			"3": 4,
			"4": 5,
		})

		game.RollDice("2") // 5
		tradeID, _ = game.MakeTradeOffer("2", map[string]int{"Grain": 1, "Ore": 1}, map[string]int{"Sheep": 1}, []string{})
		counterTradeID, _ = game.MakeCounterTradeOffer("3", tradeID, map[string]int{"Sheep": 1}, map[string]int{"Grain": 1, "Ore": 1, "Brick": 1})
		game.FinalizeTrade("2", "3", counterTradeID)
		game.BuildSettlement("2", 39)
		game.EndRound("2")
		expectedHand(map[string]string{
			"1": "SG",
			"2": "G",
			"3": "BGGO",
			"4": "GGGGGGO",
		})
		expectedSettlements(map[string][]int{
			"1": {2, 19},
			"2": {33, 35, 41, 39},
			"3": {45, 8, 54, 14},
			"4": {24},
		})
		expectedPoints(map[string]int{
			"1": 4,
			"2": 6,
			"3": 4,
			"4": 5,
		})

		game.RollDice("3") // 6
		tradeID, _ = game.MakeTradeOffer("3", map[string]int{"Grain": 1}, map[string]int{"Sheep": 1}, []string{})
		counterTradeID, _ = game.MakeCounterTradeOffer("2", tradeID, map[string]int{"Sheep": 1}, map[string]int{"Grain": 1, "Brick": 1})
		game.FinalizeTrade("3", "2", counterTradeID)
		game.BuyDevelopmentCard("3")
		game.EndRound("3")
		expectedHand(map[string]string{
			"1": "SG",
			"2": "BSSSGG",
			"3": "",
			"4": "GGGGGGO",
		})
		expectedDevHand(map[string]string{
			"1": "K",
			"2": "",
			"3": "K",
			"4": "",
		})

		game.RollDice("4") // 6
		tradeID, _ = game.MakeTradeOffer("4", map[string]int{"Grain": 1}, map[string]int{"Sheep": 1}, []string{})
		counterTradeID, _ = game.MakeCounterTradeOffer("2", tradeID, map[string]int{"Sheep": 1}, map[string]int{"Grain": 2})
		game.FinalizeTrade("4", "2", counterTradeID)
		game.BuyDevelopmentCard("4")
		game.EndRound("4")
		expectedHand(map[string]string{
			"1": "SG",
			"2": "BSSSSSSGGGG",
			"3": "",
			"4": "GGG",
		})
		expectedDevHand(map[string]string{
			"1": "K",
			"2": "",
			"3": "K",
			"4": "R",
		})

		game.RollDice("1") // 7
		if game.RoundType() != testUtils.DiscardPhase {
			t.Errorf("expected round type to be %s, but got %s", testUtils.RoundTypeTranslation[testUtils.DiscardPhase], testUtils.RoundTypeTranslation[game.RoundType()])
			return
		}
		err = game.DiscardPlayerCards("2", map[string]int{
			"Sheep": 4,
			"Grain": 1,
		})
		if err != nil {
			t.Errorf("expected to have player#2 discard 5 cards just fine, but got error %s", err.Error())
			return
		}
		if game.RoundType() != testUtils.MoveRobberDue7 {
			t.Errorf("expected round type to be %s, but got %s", testUtils.RoundTypeTranslation[testUtils.MoveRobberDue7], testUtils.RoundTypeTranslation[game.RoundType()])
			return
		}
		game.MoveRobber("1", 8)
		game.RobPlayer("1", "2")
		game.EndRound("1")
		expectedHand(map[string]string{
			"1": "SSG",
			"2": "BSGGG",
			"3": "",
			"4": "GGG",
		})

		game.RollDice("2") // 8
		game.MakePortTrade("2", 39, "Grain", "Lumber")
		game.BuildRoad("2", 53)
		game.EndRound("2")
		expectedHand(map[string]string{
			"1": "SSG",
			"2": "SGOO",
			"3": "",
			"4": "GGGOO",
		})
		expectedRoads(map[string][]int{
			"1": {2, 3, 19, 22},
			"2": {37, 43, 44, 52, 51, 50, 53},
			"3": {60, 9, 72, 16},
			"4": {28, 29, 26},
		})
		expectedLongestRoadSize(map[string]int{
			"1": 4,
			"2": 5,
			"3": 2,
			"4": 3,
		})
		expectedPoints(map[string]int{
			"1": 4,
			"2": 8,
			"3": 4,
			"4": 5,
		})

		game.RollDice("3") // 9
		game.EndRound("3")
		expectedHand(map[string]string{
			"1": "SSG",
			"2": "SGOO",
			"3": "LL",
			"4": "GGGOO",
		})

		err = game.UseRoadBuilding("4")
		if err == nil {
			t.Error("expected player not be able to use road building between rounds, but used just fine")
			return
		}
		game.RollDice("4") // 6
		game.UseRoadBuilding("4")
		game.PickRoadBuildingSpot("4", 30)
		game.PickRoadBuildingSpot("4", 31)
		tradeID, _ = game.MakeTradeOffer("4", map[string]int{"Grain": 1, "Ore": 1}, map[string]int{"Sheep": 1}, []string{})
		game.AcceptTradeOffer("2", tradeID)
		game.AcceptTradeOffer("1", tradeID)
		game.FinalizeTrade("4", "1", tradeID)
		game.BuyDevelopmentCard("4")
		game.EndRound("4")
		expectedHand(map[string]string{
			"1": "SGGO",
			"2": "SSSSSGOO",
			"3": "LLLL",
			"4": "G",
		})
		expectedDevHand(map[string]string{
			"1": "K",
			"2": "",
			"3": "K",
			"4": "K",
		})

		game.RollDice("1") // 8
		game.BuyDevelopmentCard("1")
		game.EndRound("1")
		expectedHand(map[string]string{
			"1": "G",
			"2": "SSSSSGOOOO",
			"3": "LLLL",
			"4": "GOO",
		})
		expectedDevHand(map[string]string{
			"1": "KV",
			"2": "",
			"3": "K",
			"4": "K",
		})
		expectedPoints(map[string]int{
			"1": 5,
			"2": 8,
			"3": 4,
			"4": 5,
		})

		game.RollDice("2") // 5
		game.BuyDevelopmentCard("2")
		game.BuildCity("2", 41)
		game.EndRound("2")
		expectedHand(map[string]string{
			"1": "GG",
			"2": "SSSS",
			"3": "LLLLG",
			"4": "GGGOO",
		})
		expectedSettlements(map[string][]int{
			"1": {2, 19},
			"2": {33, 35, 39},
			"3": {45, 8, 54, 14},
			"4": {24},
		})
		expectedCities(map[string][]int{
			"1": {4},
			"2": {28, 41},
			"3": {},
			"4": {22},
		})
		expectedDevHand(map[string]string{
			"1": "KV",
			"2": "K",
			"3": "K",
			"4": "K",
		})
		expectedPoints(map[string]int{
			"1": 5,
			"2": 9,
			"3": 4,
			"4": 5,
		})

		game.RollDice("3") // 6
		game.MakePortTrade("3", 54, "Lumber", "Sheep")
		game.MakePortTrade("3", 54, "Lumber", "Ore")
		game.BuyDevelopmentCard("3")
		err = game.UseKnight("3")
		if err != nil {
			t.Error(err)
			return
		}
		game.MoveRobber("3", 13)
		game.RobPlayer("3", "2")
		game.EndRound("3")
		expectedHand(map[string]string{
			"1": "GG",
			"2": "SSSSSSSS",
			"3": "LLS",
			"4": "GGGOO",
		})
		expectedDevHand(map[string]string{
			"1": "KV",
			"2": "K",
			"3": "K",
			"4": "K",
		})

		game.UseKnight("4")
		game.MoveRobber("4", 8)
		game.RobPlayer("4", "2")
		game.RollDice("4") // 7
		game.MoveRobber("4", 13)
		game.RobPlayer("4", "2")
		game.EndRound("4")
		expectedHand(map[string]string{
			"1": "GG",
			"2": "SSSSSS",
			"3": "LLS",
			"4": "SSGGGOO",
		})
		expectedDevHand(map[string]string{
			"1": "KV",
			"2": "K",
			"3": "K",
			"4": "",
		})

		game.RollDice("1") // 11
		tradeID, _ = game.MakeTradeOffer("1", map[string]int{"Grain": 1}, map[string]int{"Lumber": 2}, []string{"2"})
		game.AcceptTradeOffer("3", tradeID)
		game.FinalizeTrade("1", "3", tradeID)
		tradeID, _ = game.MakeTradeOffer("1", map[string]int{"Grain": 1}, map[string]int{"Ore": 2}, []string{"2"})
		game.AcceptTradeOffer("4", tradeID)
		game.FinalizeTrade("1", "4", tradeID)
		game.MakePortTrade("1", 4, "Ore", "Brick")
		game.MakePortTrade("1", 4, "Ore", "Brick")
		game.BuildRoad("1", 4)
		game.BuildRoad("1", 5)
		game.EndRound("1")
		expectedHand(map[string]string{
			"1": "O",
			"2": "LSSSSSS",
			"3": "LSG",
			"4": "SSGGGG",
		})
		expectedRoads(map[string][]int{
			"1": {2, 3, 19, 22, 4, 5},
			"2": {37, 43, 44, 52, 51, 50, 53},
			"3": {60, 9, 72, 16},
			"4": {28, 29, 26, 30, 31},
		})
		expectedLongestRoadSize(map[string]int{
			"1": 6,
			"2": 5,
			"3": 2,
			"4": 5,
		})
		expectedPoints(map[string]int{
			"1": 7,
			"2": 7,
			"3": 4,
			"4": 5,
		})

		game.RollDice("2") // 8
		game.BuyDevelopmentCard("2")
		game.BuyDevelopmentCard("2")
		game.MakeBankTrade("2", "Sheep", "Brick")
		game.BuildRoad("2", 32)
		game.EndRound("2")
		expectedHand(map[string]string{
			"1": "O",
			"2": "",
			"3": "LSG",
			"4": "SSGGGGOO",
		})
		expectedRoads(map[string][]int{
			"1": {2, 3, 19, 22, 4, 5},
			"2": {37, 43, 44, 52, 51, 50, 53, 32},
			"3": {60, 9, 72, 16},
			"4": {28, 29, 26, 30, 31},
		})
		expectedLongestRoadSize(map[string]int{
			"1": 6,
			"2": 6,
			"3": 2,
			"4": 5,
		})
		expectedPoints(map[string]int{
			"1": 7,
			"2": 7,
			"3": 4,
			"4": 5,
		})
		expectedDevHand(map[string]string{
			"1": "KV",
			"2": "KMM",
			"3": "K",
			"4": "",
		})

		game.RollDice("3") // 4
		game.UseKnight("3")
		game.MoveRobber("3", 8)
		game.RobPlayer("3", "2")
		tradeID, _ = game.MakeTradeOffer("3", map[string]int{"Brick": 1}, map[string]int{"Ore": 1}, []string{"2"})
		game.AcceptTradeOffer("1", tradeID)
		game.FinalizeTrade("3", "1", tradeID)
		game.BuyDevelopmentCard("3")
		game.EndRound("3")
		expectedHand(map[string]string{
			"1": "LB",
			"2": "LB",
			"3": "L",
			"4": "SSGGGGOO",
		})
		expectedDevHand(map[string]string{
			"1": "KV",
			"2": "KMM",
			"3": "K",
			"4": "",
		})
		expectedKnightsUsed(map[string]int{
			"1": 0,
			"2": 0,
			"3": 3,
			"4": 4,
		})

		game.RollDice("4") // 5
		game.BuyDevelopmentCard("4")
		game.BuyDevelopmentCard("4")
		game.EndRound("4")
		expectedHand(map[string]string{
			"1": "LBG",
			"2": "LBGG",
			"3": "LG",
			"4": "GGGG",
		})
		expectedDevHand(map[string]string{
			"1": "KV",
			"2": "KMM",
			"3": "K",
			"4": "KV",
		})
		expectedPoints(map[string]int{
			"1": 7,
			"2": 7,
			"3": 4,
			"4": 6,
		})

		game.UseKnight("1")
		game.MoveRobber("1", 13)
		game.RobPlayer("1", "2")
		game.RollDice("1") // 9
		game.BuildRoad("1", 6)
		game.EndRound("1")
		expectedHand(map[string]string{
			"1": "GG",
			"2": "LBG",
			"3": "LLLG",
			"4": "GGGG",
		})
		expectedDevHand(map[string]string{
			"1": "V",
			"2": "KMM",
			"3": "K",
			"4": "KV",
		})
		expectedRoads(map[string][]int{
			"1": {2, 3, 19, 22, 4, 5, 6},
			"2": {37, 43, 44, 52, 51, 50, 53, 32},
			"3": {60, 9, 72, 16},
			"4": {28, 29, 26, 30, 31},
		})
		expectedLongestRoadSize(map[string]int{
			"1": 7,
			"2": 6,
			"3": 2,
			"4": 5,
		})
		expectedPoints(map[string]int{
			"1": 7,
			"2": 7,
			"3": 4,
			"4": 6,
		})

		game.RollDice("2") // 8
		game.UseMonopoly("2")
		game.PickMonopolyResource("2", "Grain")
		game.MakePortTrade("2", 39, "Grain", "Ore")
		game.BuildCity("2", 35)
		game.MakePortTrade("2", 39, "Grain", "Lumber")
		game.MakePortTrade("2", 39, "Grain", "Brick")
		game.BuildRoad("2", 36)
		game.BuildRoad("2", 35)
		game.EndRound("2")
		expectedHand(map[string]string{
			"1": "",
			"2": "GG",
			"3": "LLL",
			"4": "OO",
		})
		expectedRoads(map[string][]int{
			"1": {2, 3, 19, 22, 4, 5, 6},
			"2": {37, 43, 44, 52, 51, 50, 53, 32, 36, 35},
			"3": {60, 9, 72, 16},
			"4": {28, 29, 26, 30, 31},
		})
		expectedLongestRoadSize(map[string]int{
			"1": 7,
			"2": 8,
			"3": 2,
			"4": 5,
		})
		expectedPoints(map[string]int{
			"1": 5,
			"2": 10,
			"3": 4,
			"4": 6,
		})
		if game.RoundType() != testUtils.GameOver {
			t.Errorf("expected round type to be %s, but got %s", testUtils.RoundTypeTranslation[testUtils.GameOver], testUtils.RoundTypeTranslation[game.RoundType()])
		}
	})
}
