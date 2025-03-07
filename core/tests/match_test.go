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
	// TODO: add expectedLongestRoads

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
		err = game.FinalizeCounterTrade("4", "1", counterTradeID)
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
		expectedPoints(map[string]int{
			"1": 2,
			"2": 3,
			"3": 3,
			"4": 2,
		})
		// TODO: add expected knights used

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
	})
}
