package tests

import (
	"testing"

	testUtils "github.com/victoroliveirab/settlers/core"
)

func TestPointsMultipleConfigurations(t *testing.T) {
	createGame := func(settlementMap, cityMap, roadMap map[string][]int, developmentCardsByPlayer, usedDevelopmentCardsByPlayer map[string]map[string]int) *testUtils.GameState {
		game := testUtils.CreateTestGame(
			testUtils.MockWithRoundType(testUtils.Regular),
			testUtils.MockWithSettlementsByPlayer(settlementMap),
			testUtils.MockWithCitiesByPlayer(cityMap),
			testUtils.MockWithRoadsByPlayer(roadMap),
			testUtils.MockWithDevelopmentsByPlayer(developmentCardsByPlayer),
			testUtils.MockWithUsedDevelopmentCardsByPlayer(usedDevelopmentCardsByPlayer),
			testUtils.MockWithPoints(),
		)
		return game
	}

	var tests = []struct {
		description          string
		developmentCards     map[string]map[string]int
		developmentCardsUsed map[string]map[string]int
		cityMap              map[string][]int
		roadMap              map[string][]int
		settlementMap        map[string][]int
		expectedResult       map[string]int
	}{
		{
			description: "no longest road and no most knights achieved yet",
			developmentCards: map[string]map[string]int{
				"1": {
					"Victory Point": 2,
				},
			},
			developmentCardsUsed: map[string]map[string]int{
				"2": {
					"Knight": 2,
				},
			},
			cityMap: map[string][]int{
				"3": {42},
			},
			roadMap: map[string][]int{
				"1": {1, 2, 3, 4},
				"2": {27, 30},
				"3": {54, 69},
				"4": {32, 33},
			},
			settlementMap: map[string][]int{
				"1": {1, 3},
				"2": {22, 24},
				"3": {44},
				"4": {26, 28},
			},
			expectedResult: map[string]int{
				"1": 4,
				"2": 2,
				"3": 3,
				"4": 2,
			},
		},
		{
			description: "no longest road yet, but most knights achieved",
			developmentCards: map[string]map[string]int{
				"1": {
					"Victory Point": 2,
				},
			},
			developmentCardsUsed: map[string]map[string]int{
				"2": {
					"Knight": 3,
				},
			},
			cityMap: map[string][]int{
				"3": {42},
			},
			roadMap: map[string][]int{
				"1": {1, 2, 3, 4},
				"2": {27, 30},
				"3": {54, 69},
				"4": {32, 33},
			},
			settlementMap: map[string][]int{
				"1": {1, 3},
				"2": {22, 24},
				"3": {44},
				"4": {26, 28},
			},
			expectedResult: map[string]int{
				"1": 4,
				"2": 4,
				"3": 3,
				"4": 2,
			},
		},
		{
			description: "longest road achieved, but no most knights yet",
			developmentCards: map[string]map[string]int{
				"1": {
					"Victory Point": 2,
				},
			},
			developmentCardsUsed: map[string]map[string]int{
				"2": {
					"Knight": 2,
				},
			},
			cityMap: map[string][]int{
				"3": {42},
			},
			roadMap: map[string][]int{
				"1": {1, 2, 3, 4, 5},
				"2": {27, 30},
				"3": {54, 69},
				"4": {32, 33},
			},
			settlementMap: map[string][]int{
				"1": {1, 3},
				"2": {22, 24},
				"3": {44},
				"4": {26, 28},
			},
			expectedResult: map[string]int{
				"1": 6,
				"2": 2,
				"3": 3,
				"4": 2,
			},
		},
		{
			description: "longest road and most knights achieved",
			developmentCards: map[string]map[string]int{
				"1": {
					"Victory Point": 2,
				},
			},
			developmentCardsUsed: map[string]map[string]int{
				"2": {
					"Knight": 3,
				},
			},
			cityMap: map[string][]int{
				"3": {42},
			},
			roadMap: map[string][]int{
				"1": {1, 2, 3, 4, 5},
				"2": {27, 30},
				"3": {54, 69},
				"4": {32, 33},
			},
			settlementMap: map[string][]int{
				"1": {1, 3},
				"2": {22, 24},
				"3": {44},
				"4": {26, 28},
			},
			expectedResult: map[string]int{
				"1": 6,
				"2": 4,
				"3": 3,
				"4": 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			game := createGame(tt.settlementMap, tt.cityMap, tt.roadMap, tt.developmentCards, tt.developmentCardsUsed)
			points := game.Points()
			for playerID, sum := range points {
				if sum != tt.expectedResult[playerID] {
					t.Errorf("expected player %s to have %d points, actually has %d", playerID, tt.expectedResult[playerID], sum)
				}
			}
		})
	}
}

func TestPointsIncreaseOnSettlementBuild(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 10,
				"Brick":  10,
				"Sheep":  10,
				"Grain":  10,
				"Ore":    10,
			},
		}),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {1, 3},
			"2": {22, 24},
			"3": {44},
			"4": {26, 28},
		}),
		testUtils.MockWithCitiesByPlayer(map[string][]int{
			"3": {42},
		}),
		testUtils.MockWithRoadsByPlayer(map[string][]int{
			"1": {1, 2, 3, 4, 5},
			"2": {27, 30},
			"3": {54, 69},
			"4": {32, 33},
		}),
		testUtils.MockWithDevelopmentsByPlayer(map[string]map[string]int{
			"1": {
				"Victory Point": 2,
			},
		}),
		testUtils.MockWithUsedDevelopmentCardsByPlayer(map[string]map[string]int{
			"2": {
				"Knight": 3,
			},
		}),
		testUtils.MockWithPoints(),
	)

	t.Run("points - increase after building a settlement", func(t *testing.T) {
		points := game.Points()
		if points["1"] != 6 {
			t.Errorf("expected player#1 to have 6 points before settlement build, but actually got %d", points["1"])
		}
		err := game.BuildSettlement("1", 5)
		if err != nil {
			t.Errorf("expected to build a settlement just fine, but actually got error %s", err.Error())
		}

		points = game.Points()
		if points["1"] != 7 {
			t.Errorf("expected player#1 to have 7 points after settlement build, but actually got %d", points["1"])
		}
	})
}

func TestPointsIncreaseOnCityBuild(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 10,
				"Brick":  10,
				"Sheep":  10,
				"Grain":  10,
				"Ore":    10,
			},
		}),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {1, 3},
			"2": {22, 24},
			"3": {44},
			"4": {26, 28},
		}),
		testUtils.MockWithCitiesByPlayer(map[string][]int{
			"3": {42},
		}),
		testUtils.MockWithRoadsByPlayer(map[string][]int{
			"1": {1, 2, 3, 4, 5},
			"2": {27, 30},
			"3": {54, 69},
			"4": {32, 33},
		}),
		testUtils.MockWithDevelopmentsByPlayer(map[string]map[string]int{
			"1": {
				"Victory Point": 2,
			},
		}),
		testUtils.MockWithUsedDevelopmentCardsByPlayer(map[string]map[string]int{
			"2": {
				"Knight": 3,
			},
		}),
		testUtils.MockWithPoints(),
	)

	t.Run("points - increase after building a city", func(t *testing.T) {
		points := game.Points()
		if points["1"] != 6 {
			t.Errorf("expected player#1 to have 6 points before city build, but actually got %d", points["1"])
		}
		err := game.BuildCity("1", 1)
		if err != nil {
			t.Errorf("expected to build a city just fine, but actually got error %s", err.Error())
		}

		points = game.Points()
		if points["1"] != 7 {
			t.Errorf("expected player#1 to have 7 points after settlement build, but actually got %d", points["1"])
		}
	})
}

func TestPointsIncreaseOnLongestRoadAchieved(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 10,
				"Brick":  10,
				"Sheep":  10,
				"Grain":  10,
				"Ore":    10,
			},
		}),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {1, 3},
			"2": {22, 24},
			"3": {44},
			"4": {26, 28},
		}),
		testUtils.MockWithCitiesByPlayer(map[string][]int{
			"3": {42},
		}),
		testUtils.MockWithRoadsByPlayer(map[string][]int{
			"1": {1, 2, 3, 4},
			"2": {27, 30},
			"3": {54, 69},
			"4": {32, 33},
		}),
		testUtils.MockWithDevelopmentsByPlayer(map[string]map[string]int{
			"1": {
				"Victory Point": 2,
			},
		}),
		testUtils.MockWithUsedDevelopmentCardsByPlayer(map[string]map[string]int{
			"2": {
				"Knight": 3,
			},
		}),
		testUtils.MockWithPoints(),
	)

	t.Run("points - increase after Longest Road achieved", func(t *testing.T) {
		points := game.Points()
		if points["1"] != 4 {
			t.Errorf("expected player#1 to have 4 points before longest road, but actually got %d", points["1"])
		}
		err := game.BuildRoad("1", 5)
		if err != nil {
			t.Errorf("expected to build a road just fine, but actually got error %s", err.Error())
		}

		points = game.Points()
		if points["1"] != 6 {
			t.Errorf("expected player#1 to have 6 points after longest road achieved, but actually got %d", points["1"])
		}
	})
}

func TestPointsIncreaseOnLongestRoadStolen(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 10,
				"Brick":  10,
				"Sheep":  10,
				"Grain":  10,
				"Ore":    10,
			},
		}),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {1, 3},
			"2": {22, 24},
			"3": {44},
			"4": {26, 28},
		}),
		testUtils.MockWithCitiesByPlayer(map[string][]int{
			"3": {42},
		}),
		testUtils.MockWithRoadsByPlayer(map[string][]int{
			"1": {1, 2, 3, 4},
			"2": {27, 28, 29, 30, 31},
			"3": {54, 69},
			"4": {32, 33},
		}),
		testUtils.MockWithDevelopmentsByPlayer(map[string]map[string]int{
			"1": {
				"Victory Point": 2,
			},
		}),
		testUtils.MockWithUsedDevelopmentCardsByPlayer(map[string]map[string]int{
			"2": {
				"Knight": 3,
			},
		}),
		testUtils.MockWithPoints(),
	)

	t.Run("points - increase after Longest Road stolen", func(t *testing.T) {
		points := game.Points()
		if points["1"] != 4 {
			t.Errorf("expected player#1 to have 4 points before stealing longest road, but actually got %d", points["1"])
		}
		if points["2"] != 6 {
			t.Errorf("expected player#2 to have 6 points before having longest road stolen, but actually got %d", points["2"])
		}
		err := game.BuildRoad("1", 5)
		if err != nil {
			t.Errorf("expected to build a road just fine, but actually got error %s", err.Error())
		}

		points = game.Points()
		if points["1"] != 4 {
			t.Errorf("expected player#1 to have 4 points before stealing longest road, but actually got %d", points["1"])
		}
		if points["2"] != 6 {
			t.Errorf("expected player#2 to have 6 points before having longest road stolen, but actually got %d", points["2"])
		}

		err = game.BuildRoad("1", 6)
		if err != nil {
			t.Errorf("expected to build a road just fine, but actually got error %s", err.Error())
		}

		points = game.Points()
		if points["1"] != 6 {
			t.Errorf("expected player#1 to have 6 points after stealing longest road, but actually got %d", points["1"])
		}
		if points["2"] != 4 {
			t.Errorf("expected player#2 to have 4 points after having longest road stolen, but actually got %d", points["2"])
		}
	})
}

func TestPointsDecreaseOnLongestRoadRevoked(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 10,
				"Brick":  10,
				"Sheep":  10,
				"Grain":  10,
				"Ore":    10,
			},
		}),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {1, 23},
			"2": {32, 44},
			"3": {26},
			"4": {9, 13},
		}),
		testUtils.MockWithCitiesByPlayer(map[string][]int{
			"3": {39},
		}),
		testUtils.MockWithRoadsByPlayer(map[string][]int{
			"1": {1, 29, 30, 44, 45},
			"2": {41, 43, 57, 58, 59},
			"3": {32, 53},
			"4": {15, 16},
		}),
		testUtils.MockWithDevelopmentsByPlayer(map[string]map[string]int{
			"1": {
				"Victory Point": 2,
			},
		}),
		testUtils.MockWithUsedDevelopmentCardsByPlayer(map[string]map[string]int{
			"2": {
				"Knight": 3,
			},
		}),
		testUtils.MockWithPoints(),
	)

	t.Run("points - revoke longest road on interrupt", func(t *testing.T) {
		points := game.Points()
		t.Log(points)
		if points["1"] != 4 {
			t.Errorf("expected player#1 to have 4 points before revoking longest road, but actually got %d", points["1"])
		}
		if points["2"] != 6 {
			t.Errorf("expected player#2 to have 6 points before having longest road revoked, but actually got %d", points["2"])
		}
		err := game.BuildSettlement("1", 34)
		if err != nil {
			t.Errorf("expected to build a settlement just fine, but actually got error %s", err.Error())
		}

		points = game.Points()
		if points["1"] != 5 {
			t.Errorf("expected player#1 to have 5 points after building settlement, but actually got %d", points["1"])
		}
		if points["2"] != 4 {
			t.Errorf("expected player#2 to have 4 points after having longest road revoked, but actually got %d", points["2"])
		}
	})
}

func TestPointsIncreaseOnFirstMostKnights(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 10,
				"Brick":  10,
				"Sheep":  10,
				"Grain":  10,
				"Ore":    10,
			},
		}),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {1, 3},
			"2": {22, 24},
			"3": {44},
			"4": {26, 28},
		}),
		testUtils.MockWithCitiesByPlayer(map[string][]int{
			"3": {42},
		}),
		testUtils.MockWithRoadsByPlayer(map[string][]int{
			"1": {1, 2, 3, 4},
			"2": {27, 30},
			"3": {54, 69},
			"4": {32, 33},
		}),
		testUtils.MockWithDevelopmentsByPlayer(map[string]map[string]int{
			"1": {
				"Knight": 1,
			},
			"2": {
				"Victory Point": 2,
			},
		}),
		testUtils.MockWithUsedDevelopmentCardsByPlayer(map[string]map[string]int{
			"1": {
				"Knight": 2,
			},
		}),
		testUtils.MockWithPoints(),
	)

	t.Run("points - increase after Most Knights achieved", func(t *testing.T) {
		points := game.Points()
		if points["1"] != 2 {
			t.Errorf("expected player#1 to have 2 points before most knights achieved, but actually got %d", points["1"])
		}

		err := game.UseKnight("1")
		if err != nil {
			t.Errorf("expected to use knight card just fine, but actually got error %s", err.Error())
		}

		points = game.Points()
		if points["1"] != 4 {
			t.Errorf("expected player#1 to have 4 points after most knights achieved, but actually got %d", points["1"])
		}
	})
}

func TestPointsIncreaseOnStolenMostKnights(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithCurrentRoundPlayer("2"),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"2": {
				"Lumber": 10,
				"Brick":  10,
				"Sheep":  10,
				"Grain":  10,
				"Ore":    10,
			},
		}),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {1, 3},
			"2": {22, 24},
			"3": {44},
			"4": {26, 28},
		}),
		testUtils.MockWithCitiesByPlayer(map[string][]int{
			"3": {42},
		}),
		testUtils.MockWithRoadsByPlayer(map[string][]int{
			"1": {1, 2, 3, 4},
			"2": {27, 30},
			"3": {54, 69},
			"4": {32, 33},
		}),
		testUtils.MockWithDevelopmentsByPlayer(map[string]map[string]int{
			"1": {
				"Victory Point": 2,
			},
			"2": {
				"Knight": 1,
			},
		}),
		testUtils.MockWithUsedDevelopmentCardsByPlayer(map[string]map[string]int{
			"1": {
				"Knight": 3,
			},
			"2": {
				"Knight": 3,
			},
		}),
		testUtils.MockWithPoints(),
	)

	t.Run("points - increase after Most Knights achieved", func(t *testing.T) {
		points := game.Points()
		if points["1"] != 6 {
			t.Errorf("expected player#1 to have 6 points before most knights stolen, but actually got %d", points["1"])
		}
		if points["2"] != 2 {
			t.Errorf("expected player#2 to have 2 points before stealing most knights, but actually got %d", points["2"])
		}

		err := game.UseKnight("2")
		if err != nil {
			t.Errorf("expected to use knight card just fine, but actually got error %s", err.Error())
		}

		points = game.Points()
		if points["1"] != 4 {
			t.Errorf("expected player#1 to have 4 points after most knights stolen, but actually got %d", points["1"])
		}
		if points["2"] != 4 {
			t.Errorf("expected player#2 to have 4 points after stealing most knights, but actually got %d", points["2"])
		}

	})
}

func TestGameOverOnTargetPointAchievedByBuildingSettlement(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 10,
				"Brick":  10,
				"Sheep":  10,
				"Grain":  10,
				"Ore":    10,
			},
		}),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {},
			"2": {22, 24},
			"3": {44},
			"4": {26, 28},
		}),
		testUtils.MockWithCitiesByPlayer(map[string][]int{
			"1": {1, 3, 5},
			"3": {42},
		}),
		testUtils.MockWithRoadsByPlayer(map[string][]int{
			"1": {1, 2, 3, 4, 11},
			"2": {27, 30},
			"3": {54, 69},
			"4": {32, 33},
		}),
		testUtils.MockWithDevelopmentsByPlayer(map[string]map[string]int{
			"1": {
				"Victory Point": 3,
			},
		}),
		testUtils.MockWithPoints(),
	)

	t.Run("points - reaching target point after building settlement", func(t *testing.T) {
		points := game.Points()
		if points["1"] != 9 {
			t.Errorf("expected player#1 to have 9 points before building settlement, but actually got %d", points["1"])
		}

		err := game.BuildSettlement("1", 10)
		if err != nil {
			t.Errorf("expected player#1 to build settlement just fine, but actually got error %s", err.Error())
		}

		points = game.Points()
		if points["1"] != 10 {
			t.Errorf("expected player#1 to have 10 points after building settlement, but actually got %d", points["1"])
		}

		roundType := game.RoundType()
		if roundType != testUtils.GameOver {
			t.Errorf("expected game to be over after player#1 built settlement, but actually it isn't")
		}
	})
}

func TestGameOverOnTargetPointAchievedByBuildingCity(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 10,
				"Brick":  10,
				"Sheep":  10,
				"Grain":  10,
				"Ore":    10,
			},
		}),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {10},
			"2": {22, 24},
			"3": {44},
			"4": {26, 28},
		}),
		testUtils.MockWithCitiesByPlayer(map[string][]int{
			"1": {1, 3, 5},
			"3": {42},
		}),
		testUtils.MockWithRoadsByPlayer(map[string][]int{
			"1": {1, 2, 3, 4, 11},
			"2": {27, 30},
			"3": {54, 69},
			"4": {32, 33},
		}),
		testUtils.MockWithDevelopmentsByPlayer(map[string]map[string]int{
			"1": {
				"Victory Point": 2,
			},
		}),
		testUtils.MockWithPoints(),
	)

	t.Run("points - reaching target point after building city", func(t *testing.T) {
		points := game.Points()
		if points["1"] != 9 {
			t.Errorf("expected player#1 to have 9 points before building city, but actually got %d", points["1"])
		}

		err := game.BuildCity("1", 10)
		if err != nil {
			t.Errorf("expected player#1 to build city just fine, but actually got error %s", err.Error())
		}

		points = game.Points()
		if points["1"] != 10 {
			t.Errorf("expected player#1 to have 10 points after building city, but actually got %d", points["1"])
		}

		roundType := game.RoundType()
		if roundType != testUtils.GameOver {
			t.Errorf("expected game to be over after player#1 built city, but actually it isn't")
		}
	})
}

func TestGameOverOnTargetPointAchievedByBuildingLongestRoad(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 10,
				"Brick":  10,
				"Sheep":  10,
				"Grain":  10,
				"Ore":    10,
			},
		}),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"2": {22, 24},
			"3": {44},
			"4": {26, 28},
		}),
		testUtils.MockWithCitiesByPlayer(map[string][]int{
			"1": {1, 3, 5, 10},
			"3": {42},
		}),
		testUtils.MockWithRoadsByPlayer(map[string][]int{
			"1": {1, 2, 3, 4},
			"2": {27, 30},
			"3": {54, 69},
			"4": {32, 33},
		}),
		testUtils.MockWithPoints(),
	)

	t.Run("points - reaching target point after achieving longest road", func(t *testing.T) {
		points := game.Points()
		if points["1"] != 8 {
			t.Errorf("expected player#1 to have 8 points before achieving longest road, but actually got %d", points["1"])
		}

		err := game.BuildRoad("1", 5)
		if err != nil {
			t.Errorf("expected player#1 to build road just fine, but actually got error %s", err.Error())
		}

		points = game.Points()
		if points["1"] != 10 {
			t.Errorf("expected player#1 to have 10 points after achieving longest road, but actually got %d", points["1"])
		}

		roundType := game.RoundType()
		if roundType != testUtils.GameOver {
			t.Errorf("expected game to be over after player#1 achieved longest road, but actually it isn't")
		}
	})
}

func TestGameOverOnTargetPointAchievedByAcquiringMostKnightsUse(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 10,
				"Brick":  10,
				"Sheep":  10,
				"Grain":  10,
				"Ore":    10,
			},
		}),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"2": {22, 24},
			"3": {44},
			"4": {26, 28},
		}),
		testUtils.MockWithCitiesByPlayer(map[string][]int{
			"1": {1, 3, 5, 10},
			"3": {42},
		}),
		testUtils.MockWithRoadsByPlayer(map[string][]int{
			"1": {1, 2, 3, 4},
			"2": {27, 30},
			"3": {54, 69},
			"4": {32, 33},
		}),
		testUtils.MockWithDevelopmentsByPlayer(map[string]map[string]int{
			"1": {
				"Knight": 1,
			},
		}),
		testUtils.MockWithUsedDevelopmentCardsByPlayer(map[string]map[string]int{
			"1": {
				"Knight": 2,
			},
		}),
		testUtils.MockWithPoints(),
	)

	t.Run("points - reaching target point after achieving most knights", func(t *testing.T) {
		points := game.Points()
		if points["1"] != 8 {
			t.Errorf("expected player#1 to have 8 points before achieving most knights, but actually got %d", points["1"])
		}

		err := game.UseKnight("1")
		if err != nil {
			t.Errorf("expected player#1 to build use knight just fine, but actually got error %s", err.Error())
		}

		points = game.Points()
		if points["1"] != 10 {
			t.Errorf("expected player#1 to have 10 points after achieving most knights, but actually got %d", points["1"])
		}

		roundType := game.RoundType()
		if roundType != testUtils.GameOver {
			t.Errorf("expected game to be over after player#1 achieved most knights, but actually it isn't")
		}
	})
}

func TestGameOverOnTargetPointAchievedByBuyingVictoryPoint(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 10,
				"Brick":  10,
				"Sheep":  10,
				"Grain":  10,
				"Ore":    10,
			},
		}),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {17},
			"2": {22, 24},
			"3": {44},
			"4": {26, 28},
		}),
		testUtils.MockWithCitiesByPlayer(map[string][]int{
			"1": {1, 3, 5, 10},
			"3": {42},
		}),
		testUtils.MockWithRoadsByPlayer(map[string][]int{
			"1": {1, 2, 3, 4},
			"2": {27, 30},
			"3": {54, 69},
			"4": {32, 33},
		}),
		testUtils.MockWithNextDevelopmentCard("Victory Point"),
		testUtils.MockWithPoints(),
	)

	t.Run("points - reaching target point after buying victory point", func(t *testing.T) {
		points := game.Points()
		if points["1"] != 9 {
			t.Errorf("expected player#1 to have 9 points before buying victory point, but actually got %d", points["1"])
		}

		err := game.BuyDevelopmentCard("1")
		if err != nil {
			t.Errorf("expected player#1 to buy development card just fine, but actually got error %s", err.Error())
		}

		points = game.Points()
		if points["1"] != 10 {
			t.Errorf("expected player#1 to have 10 points after buying victory point, but actually got %d", points["1"])
		}

		roundType := game.RoundType()
		if roundType != testUtils.GameOver {
			t.Errorf("expected game to be over after player#1 achieved most knights, but actually it isn't")
		}
	})
}
