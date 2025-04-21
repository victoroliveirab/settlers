package core

import (
	"testing"
)

func TestHandleDice7TilePlayerHasCards(t *testing.T) {
	rand := StubRand(7)
	game := CreateTestGame(
		MockWithRoundType(BetweenTurns),
		MockWithSettlementsByPlayer(map[string][]int{
			"1": {1},
			"2": {42},
		}),
		MockWithResourcesByPlayer(map[string]map[string]int{
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
		MockWithRand(rand),
	)

	t.Run("dice 7, will try to rob player with cards", func(t *testing.T) {
		game.RollDice("1")
		if game.RoundType() != MoveRobberDue7 {
			t.Errorf("expected round type to be %d, but it's actually %d", MoveRobberDue7, game.RoundType())
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
		if roundType != Regular {
			t.Errorf("expected round type to be %d, but it's actually %d", Regular, game.RoundType())
		}

		player1NumberOfResources := game.NumberOfCardsInHandByPlayer("1")
		player2NumberOfResources := game.NumberOfCardsInHandByPlayer("2")
		if player1NumberOfResources != 6 {
			t.Errorf("expected player#1 to have 6 cards after robbing, but actually has %d", player1NumberOfResources)
		}
		if player2NumberOfResources != 4 {
			t.Errorf("expected player#2 to have 4 cards after being robbed, but actually has %d", player2NumberOfResources)
		}
	})
}

func TestHandleDice7TilePlayerHasNoCards(t *testing.T) {
	rand := StubRand(7)
	game := CreateTestGame(
		MockWithRoundType(BetweenTurns),
		MockWithSettlementsByPlayer(map[string][]int{
			"1": {1},
			"2": {42},
		}),
		MockWithResourcesByPlayer(map[string]map[string]int{
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
		MockWithRand(rand),
	)

	t.Run("dice 7, will try to rob player with no cards", func(t *testing.T) {
		game.RollDice("1")
		if game.RoundType() != MoveRobberDue7 {
			t.Errorf("expected round type to be %d, but it's actually %d", MoveRobberDue7, game.RoundType())
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
		if roundType != Regular {
			t.Errorf("expected round type to be %d, but it's actually %d", Regular, game.RoundType())
		}

		player1NumberOfResources := game.NumberOfCardsInHandByPlayer("1")
		player2NumberOfResources := game.NumberOfCardsInHandByPlayer("2")
		if player1NumberOfResources != 5 {
			t.Errorf("expected player#1 to have 5 cards after robbing, but actually has %d", player1NumberOfResources)
		}
		if player2NumberOfResources != 0 {
			t.Errorf("expected player#2 to have 0 cards after being robbed, but actually has %d", player2NumberOfResources)
		}
	})
}

func TestHandleDice7TilePlayerRobsItself(t *testing.T) {
	rand := StubRand(7)
	game := CreateTestGame(
		MockWithRoundType(BetweenTurns),
		MockWithSettlementsByPlayer(map[string][]int{
			"1": {1},
			"2": {42},
		}),
		MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 1,
				"Brick":  1,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
		}),
		MockWithRand(rand),
	)

	t.Run("dice 7, will try to rob itself", func(t *testing.T) {
		game.RollDice("1")
		if game.RoundType() != MoveRobberDue7 {
			t.Errorf("expected round type to be %d, but it's actually %d", MoveRobberDue7, game.RoundType())
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
		if roundType != Regular {
			t.Errorf("expected round type to be %d, but it's actually %d", Regular, game.RoundType())
		}

		player1NumberOfResources := game.NumberOfCardsInHandByPlayer("1")
		if player1NumberOfResources != 5 {
			t.Errorf("expected player#1 to have 5 cards after robbing phase, but actually has %d", player1NumberOfResources)
		}
	})
}

func TestHandleDice7MoveRobberToNotOwnedTile(t *testing.T) {
	rand := StubRand(7)
	game := CreateTestGame(
		MockWithRoundType(BetweenTurns),
		MockWithSettlementsByPlayer(map[string][]int{
			"1": {1},
			"2": {42},
		}),
		MockWithResourcesByPlayer(map[string]map[string]int{
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
		MockWithRand(rand),
	)

	t.Run("dice 7, will move robber to tile not owned by anyone", func(t *testing.T) {
		game.RollDice("1")
		if game.RoundType() != MoveRobberDue7 {
			t.Errorf("expected round type to be %d, but it's actually %d", MoveRobberDue7, game.RoundType())
		}

		err := game.MoveRobber("1", 10)
		if err != nil {
			t.Errorf("expected to move robber to tile#10 just fine, but actually got error %s", err.Error())
		}

		roundType := game.RoundType()
		if roundType != Regular {
			t.Errorf("expected round type to be %d, but it's actually %d", Regular, game.RoundType())
		}

		player1NumberOfResources := game.NumberOfCardsInHandByPlayer("1")
		player2NumberOfResources := game.NumberOfCardsInHandByPlayer("2")
		if player1NumberOfResources != 5 {
			t.Errorf("expected player#1 to have 5 cards after robber moved, but actually has %d", player1NumberOfResources)
		}
		if player2NumberOfResources != 5 {
			t.Errorf("expected player#2 to have 5 cards after robber moved, but actually has %d", player2NumberOfResources)
		}
	})
}

func TestHandleDice7MoveRobberToTileOnlyOwnedByPlayer(t *testing.T) {
	rand := StubRand(7)
	game := CreateTestGame(
		MockWithRoundType(BetweenTurns),
		MockWithSettlementsByPlayer(map[string][]int{
			"1": {1},
			"2": {42},
		}),
		MockWithResourcesByPlayer(map[string]map[string]int{
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
		MockWithRand(rand),
	)

	t.Run("dice 7, will move robber to tile only owned by itself", func(t *testing.T) {
		game.RollDice("1")
		if game.RoundType() != MoveRobberDue7 {
			t.Errorf("expected round type to be %d, but it's actually %d", MoveRobberDue7, game.RoundType())
		}

		err := game.MoveRobber("1", 1)
		if err != nil {
			t.Errorf("expected to move robber to tile#1 just fine, but actually got error %s", err.Error())
		}

		roundType := game.RoundType()
		if roundType != Regular {
			t.Errorf("expected round type to be %d, but it's actually %d", Regular, game.RoundType())
		}
	})
}

func TestHandleDice7MoveRobberToOwnedTileButTriesToRobUnaffectedPlayer(t *testing.T) {
	rand := StubRand(7)
	game := CreateTestGame(
		MockWithRoundType(BetweenTurns),
		MockWithSettlementsByPlayer(map[string][]int{
			"1": {1},
			"2": {42},
			"3": {3},
		}),
		MockWithResourcesByPlayer(map[string]map[string]int{
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
			"3": {
				"Lumber": 1,
				"Brick":  1,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
		}),
		MockWithRand(rand),
	)

	t.Run("dice 7, will try to rob player that doesn't own new blocked tile", func(t *testing.T) {
		game.RollDice("1")
		if game.RoundType() != MoveRobberDue7 {
			t.Errorf("expected round type to be %s, but it's actually %s", RoundTypeTranslation[MoveRobberDue7], RoundTypeTranslation[game.RoundType()])
		}

		err := game.MoveRobber("1", 17)
		if err != nil {
			t.Errorf("expected to move robber to tile#17 just fine, but actually got error %s", err.Error())
		}

		err = game.RobPlayer("1", "3")
		if err == nil {
			t.Errorf("expected to not let player#1 rob player#3, but actually robbed just fine")
		}

		roundType := game.RoundType()
		if roundType != PickRobbed {
			t.Errorf("expected round type to be %s, but it's actually %s", RoundTypeTranslation[PickRobbed], RoundTypeTranslation[game.RoundType()])
		}

		player1NumberOfResources := game.NumberOfCardsInHandByPlayer("1")
		player3NumberOfResources := game.NumberOfCardsInHandByPlayer("2")
		if player1NumberOfResources != 5 {
			t.Errorf("expected player#1 to have 5 cards after trying to rob wrong player, but actually has %d", player1NumberOfResources)
		}
		if player3NumberOfResources != 5 {
			t.Errorf("expected player#3 to have 5 cards after try to being robbed failed, but actually has %d", player3NumberOfResources)
		}
	})
}
