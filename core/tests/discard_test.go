package tests

import (
	"testing"

	testUtils "github.com/victoroliveirab/settlers/core"
)

func TestEnterDiscardPhaseAfter7AndAPlayerHasTooMuchCards(t *testing.T) {
	rand := testUtils.StubRand(7)
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.BetweenTurns),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 3,
				"Brick":  3,
				"Sheep":  3,
				"Grain":  3,
				"Ore":    3,
			},
		}),
		testUtils.MockWithRand(rand),
	)

	t.Run("dice 7, must enter discard phase", func(t *testing.T) {
		game.RollDice("1")
		if game.RoundType() != testUtils.DiscardPhase {
			t.Errorf("expected round type to be %s, but it's actually %s", testUtils.RoundTypeTranslation[testUtils.DiscardPhase], testUtils.RoundTypeTranslation[game.RoundType()])
		}
	})
}

// TODO: make the test that goes through each type of round
// func TestDiscardPlayerCardsNotAppropriateRound(t *testing.T) {
// }

func TestDiscardPlayerCardsDoesntNeedToDiscardError(t *testing.T) {
	rand := testUtils.StubRand(7)
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.BetweenTurns),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 1,
				"Brick":  1,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
			"2": {
				"Lumber": 3,
				"Brick":  3,
				"Sheep":  3,
				"Grain":  3,
				"Ore":    3,
			},
		}),
		testUtils.MockWithRand(rand),
	)

	t.Run("discard attempt, doesn't need to discard", func(t *testing.T) {
		game.RollDice("1")
		err := game.DiscardPlayerCards("1", map[string]int{
			"Lumber": 1,
		})
		if err == nil {
			t.Errorf("expected to have error due to discarding while not needing it, but actually discarded just fine")
		}
	})
}

func TestDiscardPlayerAlreadyDiscardedError(t *testing.T) {
	rand := testUtils.StubRand(7)
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.BetweenTurns),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 3,
				"Brick":  3,
				"Sheep":  3,
				"Grain":  3,
				"Ore":    3,
			},
		}),
		testUtils.MockWithRand(rand),
	)

	t.Run("discard attempt, trying to discard twice the same turn", func(t *testing.T) {
		game.RollDice("1")
		err := game.DiscardPlayerCards("1", map[string]int{
			"Lumber": 3,
			"Brick":  3,
			"Sheep":  1,
		})
		if err != nil {
			t.Errorf("expected to discard resources correctly, but actually got error %s", err.Error())
		}

		err = game.DiscardPlayerCards("1", map[string]int{
			"Sheep": 2,
			"Grain": 2,
		})
		if err == nil {
			t.Errorf("expected to have error due to discarding multiple times the same round, but actually discarded just fine")
		}
	})

}

func TestDiscardPlayerTryToDiscardMoreResourceThanPossessedError(t *testing.T) {
	rand := testUtils.StubRand(7)
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.BetweenTurns),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 3,
				"Brick":  3,
				"Sheep":  3,
				"Grain":  3,
				"Ore":    3,
			},
		}),
		testUtils.MockWithRand(rand),
	)

	t.Run("discard attempt, tries to discard more resources than possessed", func(t *testing.T) {
		game.RollDice("1")
		err := game.DiscardPlayerCards("1", map[string]int{
			"Lumber": 3,
			"Brick":  4,
		})
		if err == nil {
			t.Errorf("expected to have error due to discarding more bricks than possessed, but actually discarded just fine")
		}
	})
}

func TestDiscardPlayerTryToDiscardLessResourcesThanRequiredError(t *testing.T) {
	rand := testUtils.StubRand(7)
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.BetweenTurns),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 3,
				"Brick":  3,
				"Sheep":  3,
				"Grain":  3,
				"Ore":    3,
			},
		}),
		testUtils.MockWithRand(rand),
	)

	t.Run("discard attempt, tries to discard less resources than required", func(t *testing.T) {
		game.RollDice("1")
		err := game.DiscardPlayerCards("1", map[string]int{
			"Lumber": 3,
			"Brick":  3,
		})
		if err == nil {
			t.Errorf("expected to have error due to discarding less resources than required, but actually discarded just fine")
		}
	})
}

func TestDiscardPlayerTryToDiscardMoreResourcesThanRequiredError(t *testing.T) {
	rand := testUtils.StubRand(7)
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.BetweenTurns),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 3,
				"Brick":  3,
				"Sheep":  3,
				"Grain":  3,
				"Ore":    3,
			},
		}),
		testUtils.MockWithRand(rand),
	)

	t.Run("discard attempt, tries to discard more resources than required", func(t *testing.T) {
		game.RollDice("1")
		err := game.DiscardPlayerCards("1", map[string]int{
			"Lumber": 3,
			"Brick":  3,
			"Sheep":  2,
		})
		if err == nil {
			t.Errorf("expected to have error due to discarding more resources than required, but actually discarded just fine")
		}
	})
}

func TestDiscardPlayerTryToDiscardMoreR(t *testing.T) {
	rand := testUtils.StubRand(7)
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.BetweenTurns),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 3,
				"Brick":  3,
				"Sheep":  3,
				"Grain":  3,
				"Ore":    3,
			},
		}),
		testUtils.MockWithRand(rand),
	)

	t.Run("discard attempt, tries to discard more resources than required", func(t *testing.T) {
		game.RollDice("1")
		err := game.DiscardPlayerCards("1", map[string]int{
			"Lumber": 3,
			"Brick":  3,
			"Sheep":  2,
		})
		if err == nil {
			t.Errorf("expected to have error due to discarding more resources than required, but actually discarded just fine")
		}
	})
}

func TestDiscardPlayerSuccess(t *testing.T) {
	rand := testUtils.StubRand(7)
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.BetweenTurns),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 3,
				"Brick":  3,
				"Sheep":  3,
				"Grain":  3,
				"Ore":    3,
			},
		}),
		testUtils.MockWithRand(rand),
	)

	t.Run("discard attempt, discards correctly", func(t *testing.T) {
		game.RollDice("1")
		err := game.DiscardPlayerCards("1", map[string]int{
			"Lumber": 3,
			"Brick":  3,
			"Sheep":  1,
		})
		if err != nil {
			t.Errorf("expected to discard resources correctly, but actually got error %s", err.Error())
		}

		player1Resources := game.ResourceHandByPlayer("1")

		if player1Resources["Lumber"] != 0 {
			t.Errorf("expected player#1 to have 0 Lumber, actually got %d", player1Resources["Lumber"])
		}
		if player1Resources["Brick"] != 0 {
			t.Errorf("expected player#1 to have 0 Brick, actually got %d", player1Resources["Brick"])
		}
		if player1Resources["Sheep"] != 2 {
			t.Errorf("expected player#1 to have 2 Sheep, actually got %d", player1Resources["Sheep"])
		}
		if player1Resources["Grain"] != 3 {
			t.Errorf("expected player#1 to have 3 Grain, actually got %d", player1Resources["Grain"])
		}
		if player1Resources["Ore"] != 3 {
			t.Errorf("expected player#1 to have 3 Ore, actually got %d", player1Resources["Ore"])
		}

		if game.RoundType() != testUtils.MoveRobberDue7 {
			t.Errorf("expected round type to be %s, but it's actually %s", testUtils.RoundTypeTranslation[testUtils.MoveRobberDue7], testUtils.RoundTypeTranslation[game.RoundType()])
		}
	})
}

func TestDiscardPlayerChangeToMoveRobberRoundAfterLastRequiredPlayerDiscards(t *testing.T) {
	rand := testUtils.StubRand(7)
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.BetweenTurns),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 3,
				"Brick":  3,
				"Sheep":  3,
				"Grain":  3,
				"Ore":    3,
			},
			"2": {
				"Lumber": 3,
				"Brick":  3,
				"Sheep":  3,
				"Grain":  3,
				"Ore":    3,
			},
		}),
		testUtils.MockWithRand(rand),
	)
	t.Run("discard phase, over only after all required players discard", func(t *testing.T) {
		game.RollDice("1")
		err := game.DiscardPlayerCards("1", map[string]int{
			"Lumber": 3,
			"Brick":  3,
			"Sheep":  1,
		})
		if err != nil {
			t.Errorf("expected to discard resources correctly, but actually got error %s", err.Error())
		}
		if game.RoundType() != testUtils.DiscardPhase {
			t.Errorf("expected round type to be %s, but it's actually %s", testUtils.RoundTypeTranslation[testUtils.DiscardPhase], testUtils.RoundTypeTranslation[game.RoundType()])
		}

		err = game.DiscardPlayerCards("2", map[string]int{
			"Lumber": 3,
			"Brick":  3,
			"Sheep":  1,
		})
		if err != nil {
			t.Errorf("expected to discard resources correctly, but actually got error %s", err.Error())
		}
		if game.RoundType() != testUtils.MoveRobberDue7 {
			t.Errorf("expected round type to be %s, but it's actually %s", testUtils.RoundTypeTranslation[testUtils.MoveRobberDue7], testUtils.RoundTypeTranslation[game.RoundType()])
		}
	})
}
