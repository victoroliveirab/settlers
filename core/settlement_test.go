package core

import (
	"testing"

	"github.com/victoroliveirab/settlers/core/packages/board"
	"github.com/victoroliveirab/settlers/core/packages/round"
	"github.com/victoroliveirab/settlers/utils"
)

func TestBuildSettlementSetupPhaseSuccess(t *testing.T) {
	game := CreateTestGame()

	t.Run("settlement build success (setup phase)", func(t *testing.T) {
		err := game.BuildSettlement("1", 42)
		if err != nil {
			t.Errorf("expected to be able to build settlement in vertex#42 during setup phase, but found error %s", err.Error())
		}
	})
}

func TestBuildSettlementRegularPhaseSuccess(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(round.Regular),
		MockWithRoadsByPlayer(map[string][]int{
			"1": {65},
		},
		),
		MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 4,
				"Brick":  3,
				"Sheep":  2,
				"Grain":  1,
				"Ore":    0,
			},
		}),
	)

	t.Run("settlement build success (regular phase)", func(t *testing.T) {
		err := game.BuildSettlement("1", 42)
		if err != nil {
			t.Errorf("expected to be able to build settlement in vertex#42 during regular phase, but found error %s", err.Error())
		}

		settlements := game.GetAllSettlements()
		t.Log(settlements)
		newSettlement := settlements[42]
		var emptyBuilding = board.Building{Owner: "", ID: 0}
		if newSettlement == emptyBuilding {
			t.Errorf("expected new settlement to show up in settlements map, but it didn't")
		}

		if newSettlement.Owner != "1" {
			t.Errorf("expected new settlement to belong to player#1, but it actually belongs to %s", newSettlement.Owner)
		}

		player1ResourcesAfterBuild := game.ResourceHandByPlayer("1")

		if player1ResourcesAfterBuild["Lumber"] != 3 {
			t.Errorf("expected to have 3 Lumber after build settlement, but found %d", player1ResourcesAfterBuild["Lumber"])
		}

		if player1ResourcesAfterBuild["Brick"] != 2 {
			t.Errorf("expected to have 2 Brick after build settlement, but found %d", player1ResourcesAfterBuild["Brick"])
		}

		if player1ResourcesAfterBuild["Sheep"] != 1 {
			t.Errorf("expected to have 1 Sheep after build settlement, but found %d", player1ResourcesAfterBuild["Sheep"])
		}

		if player1ResourcesAfterBuild["Grain"] != 0 {
			t.Errorf("expected to have 0 Grain after build settlement, but found %d", player1ResourcesAfterBuild["Grain"])
		}
	})
}

func TestBuildSettlementErrorAlreadyExistsByPlayer(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(round.Regular),
		MockWithSettlementsByPlayer(map[string][]int{
			"1": {32},
		}),
		MockWithRoadsByPlayer(map[string][]int{
			"1": {65},
		},
		),
		MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 4,
				"Brick":  3,
				"Sheep":  2,
				"Grain":  1,
				"Ore":    0,
			},
		}),
	)
	t.Run("settlement build error - player has settlement in vertex", func(t *testing.T) {
		err := game.BuildSettlement("1", 32)
		if err == nil {
			t.Errorf("expected to not be able to build settlement in vertex#32, but it built just fine")
		}
	})
}

func TestBuildSettlementErrorAlreadyExistsOtherPlayer(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(round.Regular),
		MockWithSettlementsByPlayer(map[string][]int{
			"2": {32},
		}),
		MockWithRoadsByPlayer(map[string][]int{
			"1": {65},
		},
		),
		MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 4,
				"Brick":  3,
				"Sheep":  2,
				"Grain":  1,
				"Ore":    0,
			},
		}),
	)
	t.Run("settlement build error - another player has settlement in vertex", func(t *testing.T) {
		err := game.BuildSettlement("1", 32)
		if err == nil {
			t.Errorf("expected to not be able to build settlement in vertex#32, but it built just fine")
		}
	})
}

func TestBuildSettlementErrorNotPlayerRound(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(round.Regular),
		MockWithCurrentRoundPlayer("2"),
		MockWithRoadsByPlayer(map[string][]int{
			"1": {65},
		},
		),
		MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 4,
				"Brick":  3,
				"Sheep":  2,
				"Grain":  1,
				"Ore":    0,
			},
		}),
	)
	t.Run("settlement build error - it's not the player's round", func(t *testing.T) {
		err := game.BuildSettlement("1", 42)
		if err == nil {
			t.Errorf("expected to not be able to build settlement off round, but it built just fine")
		}
	})
}

func TestBuildSettlementErrorNotEnoughResources(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(round.Regular),
		MockWithRoadsByPlayer(map[string][]int{
			"1": {65},
		},
		),
		MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 0,
				"Brick":  1,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    0,
			},
		}),
	)

	t.Run("settlement build error - player doesn't have enough resources", func(t *testing.T) {
		err := game.BuildSettlement("1", 42)
		if err == nil {
			t.Errorf("expected to not be able to build settlement without enough resources, but it built just fine")
		}
	})
}

func TestBuildSettlementErrorNotAppropriateRound(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(round.MoveRobberDue7),
		MockWithRoadsByPlayer(map[string][]int{
			"1": {65},
		},
		),
		MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 4,
				"Brick":  3,
				"Sheep":  2,
				"Grain":  1,
				"Ore":    0,
			},
		}),
	)
	t.Run("settlement build error - player tries to build in an inappropriate phase", func(t *testing.T) {
		err := game.BuildSettlement("1", 42)
		if err == nil {
			t.Errorf("expected to not be able to build settlement without being in setup or regular phase, but it built just fine")
		}
	})
}

func TestBuildSettlementErrorAlreadyExistsCity(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(round.MoveRobberDue7),
		MockWithRoadsByPlayer(map[string][]int{
			"1": {65},
		},
		),
		MockWithCitiesByPlayer(map[string][]int{
			"2": {42},
		}),
		MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 4,
				"Brick":  3,
				"Sheep":  2,
				"Grain":  1,
				"Ore":    0,
			},
		}),
	)
	t.Run("settlement build error - player tries to build in vertex occupied by city", func(t *testing.T) {
		err := game.BuildSettlement("1", 42)
		if err == nil {
			t.Errorf("expected to not be able to build settlement on top of city, but it built just fine")
		}
	})
}

func TestBuildSettlementErrorBuildingWithoutRoad(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(round.Regular),
		MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 4,
				"Brick":  3,
				"Sheep":  2,
				"Grain":  1,
				"Ore":    0,
			},
		}),
	)
	t.Run("settlement build error - player tries to build without near road", func(t *testing.T) {
		err := game.BuildSettlement("1", 42)
		if err == nil {
			t.Errorf("expected to not be able to build settlement without near road, but it build just fine")
		}
	})
}

func TestBuildSettlementErrorAlreadyBuildingSameEdge(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(round.Regular),
		MockWithRoadsByPlayer(map[string][]int{
			"1": {65},
		},
		),
		MockWithSettlementsByPlayer(map[string][]int{
			"2": {31},
		}),
		MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 4,
				"Brick":  3,
				"Sheep":  2,
				"Grain":  1,
				"Ore":    0,
			},
		}),
	)
	t.Run("settlement build error - player tries to build in an occupied edge", func(t *testing.T) {
		err := game.BuildSettlement("1", 32)
		if err == nil {
			t.Errorf("expected to not be able to build settlement in the same edge, but it build just fine")
		}
	})
}

func TestAvailableVerticesPlayerRoundRegularPhase(t *testing.T) {
	createGame := func(settlementMap, cityMap, roadMap map[string][]int) *GameState {
		game := CreateTestGame(
			MockWithRoundType(round.Regular),
			MockWithSettlementsByPlayer(settlementMap),
			MockWithCitiesByPlayer(cityMap),
			MockWithRoadsByPlayer(roadMap),
		)
		return game
	}

	var tests = []struct {
		description    string
		cityMap        map[string][]int
		roadMap        map[string][]int
		settlementMap  map[string][]int
		expectedResult []int
	}{
		{
			description: "no available vertice to build - first regular round",
			cityMap:     map[string][]int{},
			roadMap: map[string][]int{
				"1": {1, 55},
				"2": {7, 41},
			},
			settlementMap: map[string][]int{
				"1": {1, 42},
				"2": {7, 32},
			},
			expectedResult: []int{},
		},
		{
			description: "no available vertice to build - only occupied edges",
			cityMap:     map[string][]int{},
			roadMap: map[string][]int{
				"1": {1, 2, 6, 55},
				"2": {7, 41},
			},
			settlementMap: map[string][]int{
				"1": {1, 42},
				"2": {7, 32},
			},
			expectedResult: []int{},
		},
		{
			description: "available vertice to build - single edge",
			cityMap: map[string][]int{
				"1": {1},
			},
			roadMap: map[string][]int{
				"1": {1, 2, 6, 55},
				"2": {41, 43},
			},
			settlementMap: map[string][]int{
				"1": {42},
				"2": {32, 34},
			},
			expectedResult: []int{3},
		},
		{
			description: "available vertices to build - general",
			cityMap:     map[string][]int{},
			roadMap: map[string][]int{
				"1": {1, 2, 3, 4, 5, 6, 55, 56, 57},
				"2": {41, 43},
			},
			settlementMap: map[string][]int{
				"1": {1, 42},
				"2": {32, 34},
			},
			expectedResult: []int{3, 4, 5, 44},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			game := createGame(tt.settlementMap, tt.cityMap, tt.roadMap)
			vertices, err := game.AvailableVertices("1")
			if err != nil {
				t.Errorf("expected to be able to check available vertices, actually got error %s", err.Error())
			}
			expectedResultSet := utils.SetFromSlice(tt.expectedResult)
			actualResultSet := utils.SetFromSlice(vertices)
			if !expectedResultSet.Equal(actualResultSet) {
				t.Errorf("expected available vertices to be %v, actually got %v", tt.expectedResult, vertices)
			}
		})
	}

}
