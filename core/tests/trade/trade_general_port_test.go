package trade_test

import (
	"fmt"
	"testing"

	testUtils "github.com/victoroliveirab/settlers/core"
)

func TestTradeWithGeneralPortSuccess(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 4,
				"Brick":  4,
				"Sheep":  3,
				"Grain":  1,
				"Ore":    1,
			},
		}),
		testUtils.MockWithPortsByPlayer(map[string][]string{
			"1": {"General", "Lumber"},
		}),
	)

	t.Run("trade with general port - success", func(t *testing.T) {
		err := game.MakeGeneralPortTrade("1", map[string]int{
			"Brick": 3,
			"Sheep": 3,
		}, map[string]int{
			"Grain": 1,
			"Ore":   1,
		})
		if err != nil {
			t.Errorf("expected to make port trades just fine, but actually got error %s", err.Error())
		}

		player1Resources := game.ResourceHandByPlayer("1")
		if player1Resources["Lumber"] != 4 {
			t.Errorf("expected player#1 to have 4 Lumber after port trades, actually got %d", player1Resources["Lumber"])
		}
		if player1Resources["Brick"] != 1 {
			t.Errorf("expected player#1 to have 1 Brick after port trades, actually got %d", player1Resources["Brick"])
		}
		if player1Resources["Sheep"] != 0 {
			t.Errorf("expected player#1 to have 0 Sheep after port trades, actually got %d", player1Resources["Sheep"])
		}
		if player1Resources["Grain"] != 2 {
			t.Errorf("expected player#1 to have 2 Grain after port trades, actually got %d", player1Resources["Grain"])
		}
		if player1Resources["Ore"] != 2 {
			t.Errorf("expected player#1 to have 2 Ore after port trades, actually got %d", player1Resources["Ore"])
		}
	})
}

func TestTradeWithGeneralPortNotAvailablePort(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 4,
				"Brick":  3,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
		}),
		testUtils.MockWithPortsByPlayer(map[string][]string{
			"1": {"Brick"},
		}),
	)

	t.Run("trade with resource port fail - doesn't own required port", func(t *testing.T) {
		err := game.MakeGeneralPortTrade("1", map[string]int{
			"Lumber": 3,
		}, map[string]int{
			"Sheep": 1,
		})
		if err == nil {
			t.Errorf("expected to not be able to make port trade, but it actually got no error")
		}

		player1Resources := game.ResourceHandByPlayer("1")
		if player1Resources["Lumber"] != 4 {
			t.Errorf("expected player#1 to have 4 Lumber after port trades, actually got %d", player1Resources["Lumber"])
		}
		if player1Resources["Brick"] != 3 {
			t.Errorf("expected player#1 to have 3 Brick after port trades, actually got %d", player1Resources["Brick"])
		}
		if player1Resources["Sheep"] != 1 {
			t.Errorf("expected player#1 to have 1 Sheep after port trades, actually got %d", player1Resources["Sheep"])
		}
		if player1Resources["Grain"] != 1 {
			t.Errorf("expected player#1 to have 1 Grain after port trades, actually got %d", player1Resources["Grain"])
		}
		if player1Resources["Ore"] != 1 {
			t.Errorf("expected player#1 to have 1 Ore after port trades, actually got %d", player1Resources["Ore"])
		}
	})
}

func TestTradeWithGeneralPortResourcePortAvailable(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 4,
				"Brick":  3,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
		}),
		testUtils.MockWithPortsByPlayer(map[string][]string{
			"1": {"General", "Brick"},
		}),
	)

	t.Run("trade with resource port fail - owns specific port", func(t *testing.T) {
		err := game.MakeGeneralPortTrade("1", map[string]int{
			"Lumber": 3,
			"Brick":  3,
		}, map[string]int{
			"Sheep": 1,
			"Ore":   1,
		})
		if err == nil {
			t.Errorf("expected to not be able to make port trade, but it actually got no error")
		}

		player1Resources := game.ResourceHandByPlayer("1")
		if player1Resources["Lumber"] != 4 {
			t.Errorf("expected player#1 to have 4 Lumber after port trades, actually got %d", player1Resources["Lumber"])
		}
		if player1Resources["Brick"] != 3 {
			t.Errorf("expected player#1 to have 3 Brick after port trades, actually got %d", player1Resources["Brick"])
		}
		if player1Resources["Sheep"] != 1 {
			t.Errorf("expected player#1 to have 1 Sheep after port trades, actually got %d", player1Resources["Sheep"])
		}
		if player1Resources["Grain"] != 1 {
			t.Errorf("expected player#1 to have 1 Grain after port trades, actually got %d", player1Resources["Grain"])
		}
		if player1Resources["Ore"] != 1 {
			t.Errorf("expected player#1 to have 1 Ore after port trades, actually got %d", player1Resources["Ore"])
		}
	})
}

func TestTradeWithGeneralPortNotAvailableResources(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 2,
				"Brick":  3,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
		}),
		testUtils.MockWithPortsByPlayer(map[string][]string{
			"1": {"General"},
		}),
	)

	t.Run("trade with resource port fail - doesn't have required resources", func(t *testing.T) {
		err := game.MakeResourcePortTrade("1", map[string]int{
			"Lumber": 3,
			"Brick":  3,
		}, map[string]int{
			"Sheep": 1,
			"Grain": 1,
		})
		if err == nil {
			t.Errorf("expected to not be able to make port trade, but it actually got no error")
		}

		player1Resources := game.ResourceHandByPlayer("1")
		if player1Resources["Lumber"] != 2 {
			t.Errorf("expected player#1 to have 2 Lumber after port trades, actually got %d", player1Resources["Lumber"])
		}
		if player1Resources["Brick"] != 3 {
			t.Errorf("expected player#1 to have 3 Brick after port trades, actually got %d", player1Resources["Brick"])
		}
		if player1Resources["Sheep"] != 1 {
			t.Errorf("expected player#1 to have 1 Sheep after port trades, actually got %d", player1Resources["Sheep"])
		}
		if player1Resources["Grain"] != 1 {
			t.Errorf("expected player#1 to have 1 Grain after port trades, actually got %d", player1Resources["Grain"])
		}
		if player1Resources["Ore"] != 1 {
			t.Errorf("expected player#1 to have 1 Ore after port trades, actually got %d", player1Resources["Ore"])
		}
	})
}

func TestTradeWithGeneralPortByRound(t *testing.T) {
	createGame := func(roundType int) *testUtils.GameState {
		game := testUtils.CreateTestGame(
			testUtils.MockWithRoundType(roundType),
			testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
				"1": {
					"Lumber": 3,
					"Brick":  1,
					"Sheep":  1,
					"Grain":  1,
					"Ore":    1,
				},
			}),
			testUtils.MockWithPortsByPlayer(map[string][]string{
				"1": {"General"},
			}),
		)
		return game
	}

	willHaveErrorByRoundType := map[int]bool{
		testUtils.SetupSettlement1:          true,
		testUtils.SetupRoad1:                true,
		testUtils.SetupSettlement2:          true,
		testUtils.SetupRoad2:                true,
		testUtils.FirstRound:                true,
		testUtils.Regular:                   false,
		testUtils.MoveRobberDue7:            true,
		testUtils.MoveRobberDueKnight:       true,
		testUtils.PickRobbed:                true,
		testUtils.BetweenTurns:              true,
		testUtils.BuildRoad1Development:     true,
		testUtils.BuildRoad2Development:     true,
		testUtils.MonopolyPickResource:      true,
		testUtils.YearOfPlentyPickResources: true,
		testUtils.DiscardPhase:              true,
	}

	for roundType, willHaveError := range willHaveErrorByRoundType {
		testname := fmt.Sprintf("round type: %s, will have error: %v", testUtils.RoundTypeTranslation[roundType], willHaveError)
		t.Run(testname, func(t *testing.T) {
			game := createGame(roundType)
			err := game.MakeGeneralPortTrade("1", map[string]int{"Lumber": 3}, map[string]int{"Ore": 1})
			hasErr := err != nil
			if hasErr != willHaveError {
				t.Errorf("expected error to be %v, but actually was %v", willHaveError, hasErr)
			}
		})
	}
}

func TestTradeWithGeneralPortNotPlayerRound(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"2": {
				"Lumber": 3,
				"Brick":  1,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
		}),
		testUtils.MockWithPortsByPlayer(map[string][]string{
			"2": {"General"},
		}),
	)

	t.Run("trade with resource port - not player's round", func(t *testing.T) {
		err := game.MakeGeneralPortTrade("2", map[string]int{
			"Lumber": 3,
		}, map[string]int{
			"Ore": 1,
		})
		if err == nil {
			t.Errorf("expected to not be able to trade with port during other player's round, but traded just fine")
		}
	})
}

func TestTradeWithGeneralPortIncorrectNumberOfGivenResources(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 4,
				"Brick":  3,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
		}),
		testUtils.MockWithPortsByPlayer(map[string][]string{
			"1": {"General"},
		}),
	)

	t.Run("trade with resource port fail - incorrect number of given resources", func(t *testing.T) {
		err := game.MakeResourcePortTrade("1", map[string]int{
			"Lumber": 4,
			"Brick":  3,
		}, map[string]int{
			"Sheep": 2,
			"Grain": 1,
		})
		if err == nil {
			t.Errorf("expected to not be able to make port trade, but it actually got no error")
		}

		player1Resources := game.ResourceHandByPlayer("1")
		if player1Resources["Lumber"] != 4 {
			t.Errorf("expected player#1 to have 4 Lumber after port trades, actually got %d", player1Resources["Lumber"])
		}
		if player1Resources["Brick"] != 3 {
			t.Errorf("expected player#1 to have 3 Brick after port trades, actually got %d", player1Resources["Brick"])
		}
		if player1Resources["Sheep"] != 1 {
			t.Errorf("expected player#1 to have 1 Sheep after port trades, actually got %d", player1Resources["Sheep"])
		}
		if player1Resources["Grain"] != 1 {
			t.Errorf("expected player#1 to have 1 Grain after port trades, actually got %d", player1Resources["Grain"])
		}
		if player1Resources["Ore"] != 1 {
			t.Errorf("expected player#1 to have 1 Ore after port trades, actually got %d", player1Resources["Ore"])
		}
	})
}

func TestTradeWithGeneralPortIncorrectNumberOfRequestedResources(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 4,
				"Brick":  3,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
		}),
		testUtils.MockWithPortsByPlayer(map[string][]string{
			"1": {"General"},
		}),
	)

	t.Run("trade with resource port fail - incorrect number of requested resources", func(t *testing.T) {
		err := game.MakeResourcePortTrade("1", map[string]int{
			"Lumber": 3,
			"Brick":  3,
		}, map[string]int{
			"Sheep": 1,
		})
		if err == nil {
			t.Errorf("expected to not be able to make port trade, but it actually got no error")
		}

		player1Resources := game.ResourceHandByPlayer("1")
		if player1Resources["Lumber"] != 4 {
			t.Errorf("expected player#1 to have 4 Lumber after port trades, actually got %d", player1Resources["Lumber"])
		}
		if player1Resources["Brick"] != 3 {
			t.Errorf("expected player#1 to have 3 Brick after port trades, actually got %d", player1Resources["Brick"])
		}
		if player1Resources["Sheep"] != 1 {
			t.Errorf("expected player#1 to have 1 Sheep after port trades, actually got %d", player1Resources["Sheep"])
		}
		if player1Resources["Grain"] != 1 {
			t.Errorf("expected player#1 to have 1 Grain after port trades, actually got %d", player1Resources["Grain"])
		}
		if player1Resources["Ore"] != 1 {
			t.Errorf("expected player#1 to have 1 Ore after port trades, actually got %d", player1Resources["Ore"])
		}
	})
}

func TestTradeWithGeneralPortResourcePresentBothInGivenAndRequested(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 4,
				"Brick":  3,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
		}),
		testUtils.MockWithPortsByPlayer(map[string][]string{
			"1": {"General"},
		}),
	)

	t.Run("trade with resource port fail - incorrect number of given resources", func(t *testing.T) {
		err := game.MakeResourcePortTrade("1", map[string]int{
			"Lumber": 3,
			"Brick":  3,
		}, map[string]int{
			"Lumber": 1,
			"Sheep":  1,
		})
		if err == nil {
			t.Errorf("expected to not be able to make port trade, but it actually got no error")
		}

		player1Resources := game.ResourceHandByPlayer("1")
		if player1Resources["Lumber"] != 4 {
			t.Errorf("expected player#1 to have 4 Lumber after port trades, actually got %d", player1Resources["Lumber"])
		}
		if player1Resources["Brick"] != 3 {
			t.Errorf("expected player#1 to have 3 Brick after port trades, actually got %d", player1Resources["Brick"])
		}
		if player1Resources["Sheep"] != 1 {
			t.Errorf("expected player#1 to have 1 Sheep after port trades, actually got %d", player1Resources["Sheep"])
		}
		if player1Resources["Grain"] != 1 {
			t.Errorf("expected player#1 to have 1 Grain after port trades, actually got %d", player1Resources["Grain"])
		}
		if player1Resources["Ore"] != 1 {
			t.Errorf("expected player#1 to have 1 Ore after port trades, actually got %d", player1Resources["Ore"])
		}
	})
}
