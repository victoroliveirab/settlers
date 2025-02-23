package tests

import (
	testUtils "github.com/victoroliveirab/settlers/core/state"
	"testing"
)

func TestBuildSettlementSetupPhaseSuccess(t *testing.T) {
	game := testUtils.CreateTestGame()

	t.Run("settlement build success (setup phase)", func(t *testing.T) {
		err := game.BuildSettlement("1", 42)
		if err != nil {
			t.Errorf("expected to be able to build settlement in vertex#42 during setup phase, but found error %s", err.Error())
		}
	})
}

func TestBuildSettlementRegularPhaseSuccess(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithRoadsByPlayer(map[string][]int{
			"1": {65},
		},
		),
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

	t.Run("settlement build success (regular phase)", func(t *testing.T) {
		err := game.BuildSettlement("1", 42)
		if err != nil {
			t.Errorf("expected to be able to build settlement in vertex#42 during regular phase, but found error %s", err.Error())
		}

		settlements := game.AllSettlements()
		newSettlement := settlements[42]
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
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {32},
		}),
		testUtils.MockWithRoadsByPlayer(map[string][]int{
			"1": {65},
		},
		),
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
	t.Run("settlement build error - player has settlement in vertex", func(t *testing.T) {
		err := game.BuildSettlement("1", 32)
		if err == nil {
			t.Errorf("expected to not be able to build settlement in vertex#32, but it built just fine")
		}
	})
}

func TestBuildSettlementErrorAlreadyExistsOtherPlayer(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"2": {32},
		}),
		testUtils.MockWithRoadsByPlayer(map[string][]int{
			"1": {65},
		},
		),
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
	t.Run("settlement build error - another player has settlement in vertex", func(t *testing.T) {
		err := game.BuildSettlement("1", 32)
		if err == nil {
			t.Errorf("expected to not be able to build settlement in vertex#32, but it built just fine")
		}
	})
}

func TestBuildSettlementErrorNotPlayerRound(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithCurrentRoundPlayer("2"),
		testUtils.MockWithRoadsByPlayer(map[string][]int{
			"1": {65},
		},
		),
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
	t.Run("settlement build error - it's not the player's round", func(t *testing.T) {
		err := game.BuildSettlement("1", 42)
		if err == nil {
			t.Errorf("expected to not be able to build settlement off round, but it built just fine")
		}
	})
}

func TestBuildSettlementErrorNotEnoughResources(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithRoadsByPlayer(map[string][]int{
			"1": {65},
		},
		),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
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
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.MoveRobberDue7),
		testUtils.MockWithRoadsByPlayer(map[string][]int{
			"1": {65},
		},
		),
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
	t.Run("settlement build error - player tries to build in an inappropriate phase", func(t *testing.T) {
		err := game.BuildSettlement("1", 42)
		if err == nil {
			t.Errorf("expected to not be able to build settlement without being in setup or regular phase, but it built just fine")
		}
	})
}

func TestBuildSettlementErrorAlreadyExistsCity(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.MoveRobberDue7),
		testUtils.MockWithRoadsByPlayer(map[string][]int{
			"1": {65},
		},
		),
		testUtils.MockWithCitiesByPlayer(map[string][]int{
			"2": {42},
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
	t.Run("settlement build error - player tries to build in vertex occupied by city", func(t *testing.T) {
		err := game.BuildSettlement("1", 42)
		if err == nil {
			t.Errorf("expected to not be able to build settlement on top of city, but it built just fine")
		}
	})
}

func TestBuildSettlementErrorBuildingWithoutRoad(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
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
	t.Run("settlement build error - player tries to build without near road", func(t *testing.T) {
		err := game.BuildSettlement("1", 42)
		if err == nil {
			t.Errorf("expected to not be able to build settlement without near road, but it build just fine")
		}
	})
}

func TestBuildSettlementErrorAlreadyBuildingSameEdge(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithRoadsByPlayer(map[string][]int{
			"1": {65},
		},
		),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"2": {31},
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
	t.Run("settlement build error - player tries to build in an occupied edge", func(t *testing.T) {
		err := game.BuildSettlement("1", 32)
		if err == nil {
			t.Errorf("expected to not be able to build settlement in the same edge, but it build just fine")
		}
	})
}
