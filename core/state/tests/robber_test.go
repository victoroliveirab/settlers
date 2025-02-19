package tests

import (
	"testing"

	testUtils "github.com/victoroliveirab/settlers/core/state"
)

func TestHandleDice7TilePlayerHasCards(t *testing.T) {
	rand := testUtils.StubRand(7)
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.BetweenTurns),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {1},
			"2": {42},
		}),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 1,
				"Brick":  1,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
			"2": {
				"Lumber": 1,
				"Brick":  1,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
		}),
		testUtils.MockWithRand(rand),
	)

	t.Run("dice 7, will try to rob player with cards", func(t *testing.T) {
		game.RollDice("1")
		if game.RoundType() != testUtils.MoveRobberDue7 {
			t.Errorf("expected round type to be %d, but it's actually %d", testUtils.MoveRobberDue7, game.RoundType())
		}

		err := game.MoveRobber("1", 17)
		if err != nil {
			t.Errorf("expected to move robber to tile#17 just fine, but actually got error %s", err.Error())
		}

		err = game.RobPlayer("1", "2")
		if err != nil {
			t.Errorf("expected to let player#1 rob player#2 just fine, but actually got error %s", err.Error())
		}

		roundType := game.RoundType()
		if roundType != testUtils.Regular {
			t.Errorf("expected round type to be %d, but it's actually %d", testUtils.Regular, game.RoundType())
		}

		player1NumberOfResources := game.NumberOfCardsInHandByPlayer("1")
		player2NumberOfResources := game.NumberOfCardsInHandByPlayer("2")
		if player1NumberOfResources != 6 {
			t.Errorf("expected player#1 to have 6 cards after robbing, but actually has %d", player1NumberOfResources)
		}
		if player2NumberOfResources != 4 {
			t.Errorf("expected player#2 to have 4 cards after being robbed, but actually has %d", player1NumberOfResources)
		}
	})
}

func TestHandleDice7TilePlayerHasNoCards(t *testing.T) {
	rand := testUtils.StubRand(7)
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.BetweenTurns),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {1},
			"2": {42},
		}),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 1,
				"Brick":  1,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
			"2": {
				"Lumber": 0,
				"Brick":  0,
				"Sheep":  0,
				"Grain":  0,
				"Ore":    0,
			},
		}),
		testUtils.MockWithRand(rand),
	)

	t.Run("dice 7, will try to rob player with no cards", func(t *testing.T) {
		game.RollDice("1")
		if game.RoundType() != testUtils.MoveRobberDue7 {
			t.Errorf("expected round type to be %d, but it's actually %d", testUtils.MoveRobberDue7, game.RoundType())
		}

		err := game.MoveRobber("1", 17)
		if err != nil {
			t.Errorf("expected to move robber to tile#17 just fine, but actually got error %s", err.Error())
		}

		err = game.RobPlayer("1", "2")
		if err == nil {
			t.Errorf("expected to have error since player#2 has no cards, but actually no error was found")
		}

		roundType := game.RoundType()
		if roundType != testUtils.Regular {
			t.Errorf("expected round type to be %d, but it's actually %d", testUtils.Regular, game.RoundType())
		}

		player1NumberOfResources := game.NumberOfCardsInHandByPlayer("1")
		player2NumberOfResources := game.NumberOfCardsInHandByPlayer("2")
		if player1NumberOfResources != 5 {
			t.Errorf("expected player#1 to have 5 cards after robbing, but actually has %d", player1NumberOfResources)
		}
		if player2NumberOfResources != 0 {
			t.Errorf("expected player#2 to have 0 cards after being robbed, but actually has %d", player1NumberOfResources)
		}
	})
}

func TestHandleDice7TilePlayerRobsItself(t *testing.T) {
	rand := testUtils.StubRand(7)
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.BetweenTurns),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {1},
			"2": {42},
		}),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 1,
				"Brick":  1,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
		}),
		testUtils.MockWithRand(rand),
	)

	t.Run("dice 7, will try to rob itself", func(t *testing.T) {
		game.RollDice("1")
		if game.RoundType() != testUtils.MoveRobberDue7 {
			t.Errorf("expected round type to be %d, but it's actually %d", testUtils.MoveRobberDue7, game.RoundType())
		}

		err := game.MoveRobber("1", 1)
		if err != nil {
			t.Errorf("expected to move robber to tile#1 just fine, but actually got error %s", err.Error())
		}

		err = game.RobPlayer("1", "1")
		if err == nil {
			t.Errorf("expected to have error since cannot rob from yourself, but actually no error was found")
		}

		roundType := game.RoundType()
		if roundType != testUtils.Regular {
			t.Errorf("expected round type to be %d, but it's actually %d", testUtils.Regular, game.RoundType())
		}

		player1NumberOfResources := game.NumberOfCardsInHandByPlayer("1")
		if player1NumberOfResources != 5 {
			t.Errorf("expected player#1 to have 5 cards after robbing phase, but actually has %d", player1NumberOfResources)
		}
	})
}
