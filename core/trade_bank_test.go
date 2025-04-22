package core

import (
	"fmt"
	"testing"

	"github.com/victoroliveirab/settlers/core/packages/round"
)

func TestTradeWithBankWithAvailableResources(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(round.Regular),
		MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 4,
				"Brick":  1,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
		}),
	)

	t.Run("trade with bank - player has available resources", func(t *testing.T) {
		err := game.MakeBankTrade("1", map[string]int{"Lumber": 4}, map[string]int{"Ore": 1})
		if err != nil {
			t.Errorf("expected to trade with bank just fine, but actually got error %s", err.Error())
		}

		player1Resources := game.ResourceHandByPlayer("1")
		if player1Resources["Lumber"] != 0 {
			t.Errorf("expected player#1 to have 0 Lumber, actually got %d", player1Resources["Lumber"])
		}
		if player1Resources["Ore"] != 2 {
			t.Errorf("expected player#1 to have 2 Ore, actually got %d", player1Resources["Ore"])
		}
	})
}

func TestTradeWithBankWithNoAvailableResources(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(round.Regular),
		MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 3,
				"Brick":  1,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
		}),
	)

	t.Run("trade with bank - player doesn't have available resources", func(t *testing.T) {
		err := game.MakeBankTrade("1", map[string]int{"Lumber": 4}, map[string]int{"Ore": 1})
		if err == nil {
			t.Errorf("expected to not be able to trade with bank, but actually traded just fine")
		}

		player1Resources := game.ResourceHandByPlayer("1")
		if player1Resources["Lumber"] != 3 {
			t.Errorf("expected player#1 to have 3 Lumber, actually got %d", player1Resources["Lumber"])
		}
		if player1Resources["Ore"] != 1 {
			t.Errorf("expected player#1 to have 1 Ore, actually got %d", player1Resources["Ore"])
		}
	})
}

func TestTradeWithBankByRound(t *testing.T) {
	createGame := func(roundType round.Type) *GameState {
		game := CreateTestGame(
			MockWithRoundType(roundType),
			MockWithResourcesByPlayer(map[string]map[string]int{
				"1": {
					"Lumber": 4,
					"Brick":  1,
					"Sheep":  1,
					"Grain":  1,
					"Ore":    1,
				},
			}),
		)
		return game
	}

	willHaveErrorByRoundType := map[round.Type]bool{
		round.SetupSettlement1:          true,
		round.SetupRoad1:                true,
		round.SetupSettlement2:          true,
		round.SetupRoad2:                true,
		round.FirstRound:                true,
		round.Regular:                   false,
		round.MoveRobberDue7:            true,
		round.MoveRobberDueKnight:       true,
		round.PickRobbed:                true,
		round.BetweenTurns:              true,
		round.BuildRoad1Development:     true,
		round.BuildRoad2Development:     true,
		round.MonopolyPickResource:      true,
		round.YearOfPlentyPickResources: true,
		round.DiscardPhase:              true,
	}

	for roundType, willHaveError := range willHaveErrorByRoundType {
		testname := fmt.Sprintf("round type: %d, will have error: %v", roundType, willHaveError)
		t.Run(testname, func(t *testing.T) {
			game := createGame(roundType)
			err := game.MakeBankTrade("1", map[string]int{"Lumber": 4}, map[string]int{"Ore": 1})
			hasErr := err != nil
			if hasErr != willHaveError {
				t.Errorf("expected error to be %v, but actually was %v", willHaveError, hasErr)
			}
		})
	}
}

func TestTradeWithBankNotPlayerRound(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(round.Regular),
		MockWithResourcesByPlayer(map[string]map[string]int{
			"2": {
				"Lumber": 1,
				"Brick":  1,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
		}),
	)

	t.Run("trade with bank - not player's round", func(t *testing.T) {
		err := game.MakeBankTrade("1", map[string]int{"Lumber": 4}, map[string]int{"Ore": 1})
		if err == nil {
			t.Errorf("expected to not be able to trade with bank during other player's round, but traded just fine")
		}
	})
}
