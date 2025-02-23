package tests

import (
	"fmt"
	"testing"

	testUtils "github.com/victoroliveirab/settlers/core/state"
)

func TestBuildRoadSetup1PhaseSuccess(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {42},
		}),
		testUtils.MockWithRoundType(testUtils.SetupRoad1),
	)

	t.Run("road build success (setup phase)", func(t *testing.T) {
		err := game.BuildRoad("1", 65)
		if err != nil {
			t.Errorf("expected to be able to build road in edge#65 during setup phase, but found error %s", err.Error())
		}
	})
}

func TestBuildRoadSetup2PhaseSuccess(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {42, 4},
		}),
		testUtils.MockWithRoundType(testUtils.SetupRoad1),
	)

	t.Run("road build success (setup phase)", func(t *testing.T) {
		err := game.BuildRoad("1", 11)
		if err != nil {
			t.Errorf("expected to be able to build road in edge#4 during setup phase, but found error %s", err.Error())
		}
	})
}

func TestBuildRoadRegularPhaseTouchingSettlementSuccess(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {42},
		}),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 4,
				"Brick":  3,
				"Sheep":  2,
				"Grain":  1,
				"Ore":    0,
			},
		}),
	)

	t.Run("road build success (regular phase)", func(t *testing.T) {
		err := game.BuildRoad("1", 65)
		if err != nil {
			t.Errorf("expected to be able to build road in edge#65 during setup phase, but found error %s", err.Error())
		}

		roads := game.AllRoads()
		newRoad := roads[65]
		if newRoad == emptyBuilding {
			t.Errorf("expected new road to show up in roads map, but it didn't")
		}

		if newRoad.Owner != "1" {
			t.Errorf("expected new road to belong to player#1, but it actually belongs to %s", newRoad.Owner)
		}

		player1ResourcesAfterBuild := game.ResourceHandByPlayer("1")

		if player1ResourcesAfterBuild["Lumber"] != 3 {
			t.Errorf("expected to have 3 Lumber after build road, but found %d", player1ResourcesAfterBuild["Lumber"])
		}

		if player1ResourcesAfterBuild["Brick"] != 2 {
			t.Errorf("expected to have 2 Brick after build road, but found %d", player1ResourcesAfterBuild["Brick"])
		}

		if player1ResourcesAfterBuild["Sheep"] != 2 {
			t.Errorf("expected to have 2 Sheep after build road, but found %d", player1ResourcesAfterBuild["Sheep"])
		}

		if player1ResourcesAfterBuild["Grain"] != 1 {
			t.Errorf("expected to have 1 Grain after build road, but found %d", player1ResourcesAfterBuild["Grain"])
		}
	})
}

func TestBuildRoadRegularPhaseTouchingRoadSuccess(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithRoadsByPlayer(map[string][]int{
			"1": {54},
		}),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 4,
				"Brick":  3,
				"Sheep":  2,
				"Grain":  1,
				"Ore":    0,
			},
		}),
	)

	t.Run("road build success (regular phase)", func(t *testing.T) {
		err := game.BuildRoad("1", 65)
		if err != nil {
			t.Errorf("expected to be able to build road in edge#65 during setup phase, but found error %s", err.Error())
		}

		roads := game.AllRoads()
		newRoad := roads[65]
		if newRoad == emptyBuilding {
			t.Errorf("expected new road to show up in roads map, but it didn't")
		}

		if newRoad.Owner != "1" {
			t.Errorf("expected new road to belong to player#1, but it actually belongs to %s", newRoad.Owner)
		}

		player1ResourcesAfterBuild := game.ResourceHandByPlayer("1")

		if player1ResourcesAfterBuild["Lumber"] != 3 {
			t.Errorf("expected to have 3 Lumber after build road, but found %d", player1ResourcesAfterBuild["Lumber"])
		}

		if player1ResourcesAfterBuild["Brick"] != 2 {
			t.Errorf("expected to have 2 Brick after build road, but found %d", player1ResourcesAfterBuild["Brick"])
		}

		if player1ResourcesAfterBuild["Sheep"] != 2 {
			t.Errorf("expected to have 2 Sheep after build road, but found %d", player1ResourcesAfterBuild["Sheep"])
		}

		if player1ResourcesAfterBuild["Grain"] != 1 {
			t.Errorf("expected to have 1 Grain after build road, but found %d", player1ResourcesAfterBuild["Grain"])
		}
	})
}

func TestBuildRoadErrorSetupPhaseNotAttachedToLastSettlement(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {42, 4},
		}),
		testUtils.MockWithRoadsByPlayer(map[string][]int{
			"1": {65},
		}),
		testUtils.MockWithRoundType(testUtils.SetupRoad1),
	)

	t.Run("road build success (setup phase)", func(t *testing.T) {
		err := game.BuildRoad("1", 54)
		if err == nil {
			t.Errorf("expected to not be able to build road in edge#54 during setup phase 2, but it built just fine")
		}
	})
}

func TestBuildRoadErrorAlreadyExists(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithRoadsByPlayer(map[string][]int{
			"1": {65},
		}),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 4,
				"Brick":  3,
				"Sheep":  2,
				"Grain":  1,
				"Ore":    0,
			},
		}),
	)
	t.Run("road build error - another player has road in edge", func(t *testing.T) {
		err := game.BuildRoad("1", 65)
		if err == nil {
			t.Errorf("expected to not be able to build road in edge#65, but it built just fine")
		}
	})
}

func TestBuildRoadErrorNotPlayerRound(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithCurrentRoundPlayer("2"),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {42},
		}),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 4,
				"Brick":  3,
				"Sheep":  2,
				"Grain":  1,
				"Ore":    0,
			},
		}),
	)
	t.Run("road build error - it's not the player's round", func(t *testing.T) {
		err := game.BuildRoad("1", 65)
		if err == nil {
			t.Errorf("expected to not be able to build road off round, but it built just fine")
		}
	})
}

func TestBuildRoadErrorNotEnoughResources(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {42},
		}),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 4,
				"Brick":  0,
				"Sheep":  2,
				"Grain":  1,
				"Ore":    0,
			},
		}),
	)
	t.Run("road build error - player doesn't have enough resources", func(t *testing.T) {
		err := game.BuildRoad("1", 65)
		if err == nil {
			t.Errorf("expected to not be able to build road without enough resources, but it built just fine")
		}
	})
}

func TestBuildRoadErrorNotAppropriateRound(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.MoveRobberDue7),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {42},
		}),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 4,
				"Brick":  3,
				"Sheep":  2,
				"Grain":  1,
				"Ore":    0,
			},
		}),
	)
	t.Run("road build error - player tries to build in an inappropriate phase", func(t *testing.T) {
		err := game.BuildRoad("1", 65)
		if err == nil {
			t.Errorf("expected to not be able to build road without being in setup or regular phase, but it built just fine")
		}
	})
}

func TestBuildRoadErrorNotRoadOrSettlementAttached(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.MoveRobberDue7),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {42},
		}),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 4,
				"Brick":  3,
				"Sheep":  2,
				"Grain":  1,
				"Ore":    0,
			},
		}),
	)
	t.Run("road build error - player tries to build road without connection", func(t *testing.T) {
		err := game.BuildRoad("1", 1)
		if err == nil {
			t.Errorf("expected to not be able to build road without being connected to settlement/city or road, but it built just fine")
		}
	})
}

func TestLongestRoad(t *testing.T) {
	var tests = []struct {
		description    string
		edges          []int
		settlementMap  map[string][]int
		expectedResult int
	}{
		{
			description:    "simple line",
			edges:          []int{1, 2, 7, 8, 12},
			expectedResult: 5,
		},
		{
			description:    "simple line with branch",
			edges:          []int{1, 2, 7, 8, 9, 12},
			expectedResult: 5,
		},
		{
			description:    "simple most outer ring",
			edges:          []int{4, 5, 6, 10, 11, 14, 15, 16, 20, 21, 30, 31, 32, 35, 36, 47, 48, 49, 50, 53, 61, 62, 63, 64, 66, 67, 68, 70, 71, 72},
			expectedResult: 30,
		},
		{
			description:    "most outer ring with internal roads",
			edges:          []int{1, 4, 5, 6, 10, 11, 14, 15, 16, 20, 21, 30, 31, 32, 35, 36, 47, 48, 49, 50, 53, 61, 62, 63, 64, 66, 67, 68, 70, 71, 72},
			expectedResult: 31,
		},
		{
			description:    "hangman",
			edges:          []int{23, 25, 26, 38, 39, 40, 41, 42, 43, 44},
			expectedResult: 8,
		},
		{
			description:    "tile loop",
			edges:          []int{1, 2, 3, 4, 5, 6},
			expectedResult: 6,
		},
		{
			description:    "double tile loop",
			edges:          []int{1, 2, 3, 4, 5, 6, 7, 19, 22, 23, 24},
			expectedResult: 11,
		},
		{
			description:    "triple tile loop",
			edges:          []int{1, 2, 3, 4, 5, 6, 7, 19, 22, 23, 24, 25, 39, 40, 41, 42},
			expectedResult: 15,
		},
		{
			description:    "double tile loop with connecting edge",
			edges:          []int{1, 2, 3, 4, 5, 6, 7, 8, 12, 24, 25, 26, 27},
			expectedResult: 13,
		},
		{
			description:    "interrupted line",
			edges:          []int{1, 2, 7, 8},
			expectedResult: 2,
			settlementMap: map[string][]int{
				"1": {1, 8},
				"2": {3},
			},
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("longest road - %s", tt.description)
		t.Run(testname, func(t *testing.T) {
			game := testUtils.CreateTestGame(
				testUtils.MockWithSettlementsByPlayer(tt.settlementMap),
				testUtils.MockWithRoadsByPlayer(map[string][]int{
					"1": tt.edges,
				}),
			)
			length := game.LongestRoadLengthByPlayer("1")
			if length != tt.expectedResult {
				t.Errorf("expected longest road to have length %d, but got %d", tt.expectedResult, length)
			}
		})
	}
}
