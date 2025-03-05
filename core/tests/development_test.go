package tests

import (
	"fmt"
	"testing"

	testUtils "github.com/victoroliveirab/settlers/core"
)

func TestBuyDevelopmentCardNotPlayerRound(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithCurrentRoundPlayer("2"),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 2,
				"Brick":  2,
				"Sheep":  2,
				"Grain":  2,
				"Ore":    2,
			},
		}),
	)
	t.Run("buy development card error - it's not the player's round", func(t *testing.T) {
		err := game.BuyDevelopmentCard("1")
		if err == nil {
			t.Errorf("expected to not be able to buy development card off round, but it bought just fine")
		}
	})
}

func TestBuyDevelopmentCardNotEnoughResources(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 2,
				"Brick":  2,
				"Sheep":  2,
				"Grain":  2,
				"Ore":    0,
			},
		}),
	)

	t.Run("buy development card error - player doesn't have enough resources", func(t *testing.T) {
		err := game.BuyDevelopmentCard("1")
		if err == nil {
			t.Errorf("expected to not be able to buy development card without enough resources, but it bought just fine")
		}
	})
}

func TestPlayKnightDevelopmentCardRobOpponentWithCards(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {1},
			"2": {42},
		}),
		testUtils.MockWithDevelopmentsByPlayer(map[string]map[string]int{
			"1": {
				"Knight": 1,
			},
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
	)

	t.Run("player has knight card, will try to rob player with cards", func(t *testing.T) {
		err := game.UseKnight("1")
		if err != nil {
			t.Errorf("expected to use knight card just fine, but actually got error %s", err.Error())
		}
		if game.RoundType() != testUtils.MoveRobberDueKnight {
			t.Errorf("expected round type to be %s after knight use, but got %s", testUtils.RoundTypeTranslation[testUtils.MoveRobberDueKnight], testUtils.RoundTypeTranslation[game.RoundType()])
		}

		err = game.MoveRobber("1", 17)
		if err != nil {
			t.Errorf("expected to move robber to tile#17 just fine, but actually got error %s", err.Error())
		}
		if game.RoundType() != testUtils.PickRobbed {
			t.Errorf("expected round type to be %s after knight use, but got %s", testUtils.RoundTypeTranslation[testUtils.PickRobbed], testUtils.RoundTypeTranslation[game.RoundType()])
		}

		err = game.RobPlayer("1", "2")
		if err != nil {
			t.Errorf("expected to let player#1 rob player#2 just fine, but actually got error %s", err.Error())
		}
		if game.RoundType() != testUtils.Regular {
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

func TestPlayKnightDevelopmentCardRobOpponentWithNoCards(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {1},
			"2": {42},
		}),
		testUtils.MockWithDevelopmentsByPlayer(map[string]map[string]int{
			"1": {
				"Knight": 1,
			},
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
	)

	t.Run("player has knight card, will try to rob player with no cards", func(t *testing.T) {
		err := game.UseKnight("1")
		if err != nil {
			t.Errorf("expected to use knight just fine, but actyally got error %s", err.Error())
		}
		if game.RoundType() != testUtils.MoveRobberDueKnight {
			t.Errorf("expected round type to be %s after knight use, but got %s", testUtils.RoundTypeTranslation[testUtils.MoveRobberDueKnight], testUtils.RoundTypeTranslation[game.RoundType()])
		}

		err = game.MoveRobber("1", 17)
		if err != nil {
			t.Errorf("expected to move robber to tile#17 just fine, but actually got error %s", err.Error())
		}
		if game.RoundType() != testUtils.PickRobbed {
			t.Errorf("expected round type to be %s after knight use, but got %s", testUtils.RoundTypeTranslation[testUtils.PickRobbed], testUtils.RoundTypeTranslation[game.RoundType()])
		}

		err = game.RobPlayer("1", "2")
		if err == nil {
			t.Errorf("expected to have error since player#2 has no cards, but actually no error was found")
		}
		if game.RoundType() != testUtils.Regular {
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

func TestPlayKnightDevelopmentCardRobItself(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {1},
			"2": {42},
		}),
		testUtils.MockWithDevelopmentsByPlayer(map[string]map[string]int{
			"1": {
				"Knight": 1,
			},
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
	)

	t.Run("player has knight card, will try to rob player with no cards", func(t *testing.T) {
		err := game.UseKnight("1")
		if err != nil {
			t.Errorf("expected to use knight just fine, but actyally got error %s", err.Error())
		}
		if game.RoundType() != testUtils.MoveRobberDueKnight {
			t.Errorf("expected round type to be %s after knight use, but got %s", testUtils.RoundTypeTranslation[testUtils.MoveRobberDueKnight], testUtils.RoundTypeTranslation[game.RoundType()])
		}

		err = game.MoveRobber("1", 17)
		if err != nil {
			t.Errorf("expected to move robber to tile#1 just fine, but actually got error %s", err.Error())
		}
		if game.RoundType() != testUtils.PickRobbed {
			t.Errorf("expected round type to be %s after knight use, but got %s", testUtils.RoundTypeTranslation[testUtils.PickRobbed], testUtils.RoundTypeTranslation[game.RoundType()])
		}

		err = game.RobPlayer("1", "1")
		if err == nil {
			t.Errorf("expected to have error since cannot rob from yourself, but actually no error was found")
		}
		if game.RoundType() != testUtils.Regular {
			t.Errorf("expected round type to be %d, but it's actually %d", testUtils.Regular, game.RoundType())
		}

		player1NumberOfResources := game.NumberOfCardsInHandByPlayer("1")
		if player1NumberOfResources != 5 {
			t.Errorf("expected player#1 to have 5 cards after robbing, but actually has %d", player1NumberOfResources)
		}
	})
}

func TestPlayKnightDevelopmentCardWithoutHavingOne(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {1},
			"2": {42},
		}),
		testUtils.MockWithDevelopmentsByPlayer(map[string]map[string]int{
			"1": {
				"Knight": 0,
			},
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
	)

	t.Run("player has no knight card, will try to play one", func(t *testing.T) {
		err := game.UseKnight("1")
		if err == nil {
			t.Errorf("expected to not be able to use knight, but actually no error was found")
		}
	})
}

func TestPlayKnightDevelopmentCardByRound(t *testing.T) {
	createGame := func(roundType int) *testUtils.GameState {
		game := testUtils.CreateTestGame(
			testUtils.MockWithRoundType(roundType),
			testUtils.MockWithDevelopmentsByPlayer(map[string]map[string]int{
				"1": {
					"Knight": 1,
				},
			}),
		)
		return game
	}

	willHaveErrorByRoundType := map[int]bool{
		testUtils.SetupSettlement1:          true,
		testUtils.SetupRoad1:                true,
		testUtils.SetupSettlement2:          true,
		testUtils.SetupRoad2:                true,
		testUtils.FirstRound:                false,
		testUtils.Regular:                   false,
		testUtils.MoveRobberDue7:            true,
		testUtils.MoveRobberDueKnight:       true,
		testUtils.PickRobbed:                true,
		testUtils.BetweenTurns:              false,
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
			err := game.UseKnight("1")
			hasErr := err != nil
			if hasErr != willHaveError {
				t.Errorf("expected error to be %v, but actually was %v", willHaveError, hasErr)
			}
		})
	}
}

func TestPlayMonopolyDevelopmentCardOpponentsHaveResource(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {1},
		}),
		testUtils.MockWithDevelopmentsByPlayer(map[string]map[string]int{
			"1": {
				"Monopoly": 1,
			},
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
				"Lumber": 3,
				"Ore":    3,
			},
			"3": {
				"Lumber": 2,
				"Ore":    2,
			},
			"4": {
				"Lumber": 3,
				"Ore":    0,
			},
		}),
	)

	t.Run("player has monopoly card, will try to rob available resource", func(t *testing.T) {
		err := game.UseMonopoly("1")
		if err != nil {
			t.Errorf("expected to use monopoly card just fine, but actually got error %s", err.Error())
		}
		if game.RoundType() != testUtils.MonopolyPickResource {
			t.Errorf("expected round type to be %s after monopoly use, but got %s", testUtils.RoundTypeTranslation[testUtils.MonopolyPickResource], testUtils.RoundTypeTranslation[game.RoundType()])
		}

		err = game.PickMonopolyResource("1", "Ore")
		if err != nil {
			t.Errorf("expected to monopolyze resource just fine, but actually got error %s", err.Error())
		}
		if game.RoundType() != testUtils.Regular {
			t.Errorf("expected round type to be %s after picking monopoly resource, but it's actually %s", testUtils.RoundTypeTranslation[testUtils.Regular], testUtils.RoundTypeTranslation[game.RoundType()])
		}

		player1NumberOfResources := game.NumberOfCardsInHandByPlayer("1")
		player2NumberOfResources := game.NumberOfCardsInHandByPlayer("2")
		player3NumberOfResources := game.NumberOfCardsInHandByPlayer("3")
		player4NumberOfResources := game.NumberOfCardsInHandByPlayer("4")
		if player1NumberOfResources != 10 {
			t.Errorf("expected player#1 to have 10 cards after monopoly use, but actually has %d", player1NumberOfResources)
		}
		if player2NumberOfResources != 3 {
			t.Errorf("expected player#2 to have 3 cards after monopoly use, but actually has %d", player2NumberOfResources)
		}
		if player3NumberOfResources != 2 {
			t.Errorf("expected player#3 to have 2 cards after monopoly use, but actually has %d", player3NumberOfResources)
		}
		if player4NumberOfResources != 3 {
			t.Errorf("expected player#4 to have 3 cards after monopoly use, but actually has %d", player4NumberOfResources)
		}
	})
}

func TestPlayMonopolyDevelopmentCardOpponentsDontHaveResource(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {1},
		}),
		testUtils.MockWithDevelopmentsByPlayer(map[string]map[string]int{
			"1": {
				"Monopoly": 1,
			},
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
				"Lumber": 3,
				"Ore":    0,
			},
			"3": {
				"Lumber": 2,
				"Ore":    0,
			},
			"4": {
				"Lumber": 3,
				"Ore":    0,
			},
		}),
	)

	t.Run("player has monopoly card, will try to rob unavailable resource", func(t *testing.T) {
		err := game.UseMonopoly("1")
		if err != nil {
			t.Errorf("expected to rob all Ore just fine, but actually got error %s", err.Error())
		}
		if game.RoundType() != testUtils.MonopolyPickResource {
			t.Errorf("expected round type to be %s after monopoly use, but got %s", testUtils.RoundTypeTranslation[testUtils.MonopolyPickResource], testUtils.RoundTypeTranslation[game.RoundType()])
		}
		err = game.PickMonopolyResource("1", "Ore")
		if err != nil {
			t.Errorf("expected to monopolyze resource just fine, but actually got error %s", err.Error())
		}
		if game.RoundType() != testUtils.Regular {
			t.Errorf("expected round type to be %s after picking monopoly resource, but it's actually %s", testUtils.RoundTypeTranslation[testUtils.Regular], testUtils.RoundTypeTranslation[game.RoundType()])
		}

		player1NumberOfResources := game.NumberOfCardsInHandByPlayer("1")
		player2NumberOfResources := game.NumberOfCardsInHandByPlayer("2")
		player3NumberOfResources := game.NumberOfCardsInHandByPlayer("3")
		player4NumberOfResources := game.NumberOfCardsInHandByPlayer("4")
		if player1NumberOfResources != 5 {
			t.Errorf("expected player#1 to have 5 cards after monopoly use, but actually has %d", player1NumberOfResources)
		}
		if player2NumberOfResources != 3 {
			t.Errorf("expected player#2 to have 3 cards after monopoly use, but actually has %d", player2NumberOfResources)
		}
		if player3NumberOfResources != 2 {
			t.Errorf("expected player#3 to have 2 cards after monopoly use, but actually has %d", player3NumberOfResources)
		}
		if player4NumberOfResources != 3 {
			t.Errorf("expected player#4 to have 3 cards after monopoly use, but actually has %d", player4NumberOfResources)
		}
	})
}

func TestPlayMonopolyDevelopmentCardWithoutHavingOne(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {1},
		}),
		testUtils.MockWithDevelopmentsByPlayer(map[string]map[string]int{
			"1": {
				"Monopoly": 0,
			},
		}),
	)

	t.Run("player has no monopoly card, will try to play one", func(t *testing.T) {
		err := game.UseMonopoly("1")
		if err == nil {
			t.Errorf("expected to not be able to use monopoly, but actually no error was found")
		}
	})
}

func TestPlayMonopolyDevelopmentCardByRound(t *testing.T) {
	createGame := func(roundType int) *testUtils.GameState {
		game := testUtils.CreateTestGame(
			testUtils.MockWithRoundType(roundType),
			testUtils.MockWithDevelopmentsByPlayer(map[string]map[string]int{
				"1": {
					"Monopoly": 1,
				},
			}),
		)
		return game
	}

	willHaveErrorByRoundType := map[int]bool{
		testUtils.SetupSettlement1:          true,
		testUtils.SetupRoad1:                true,
		testUtils.SetupSettlement2:          true,
		testUtils.SetupRoad2:                true,
		testUtils.FirstRound:                false,
		testUtils.Regular:                   false,
		testUtils.MoveRobberDue7:            true,
		testUtils.MoveRobberDueKnight:       true,
		testUtils.PickRobbed:                true,
		testUtils.BetweenTurns:              false,
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
			err := game.UseMonopoly("1")
			hasErr := err != nil
			if hasErr != willHaveError {
				t.Errorf("expected error to be %v, but actually was %v", willHaveError, hasErr)
			}
		})
	}
}

func TestPlayRoadBuildingDevelopmentCardAvailablePathAvailableRoads(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {1},
		}),
		testUtils.MockWithDevelopmentsByPlayer(map[string]map[string]int{
			"1": {
				"Road Building": 1,
			},
		}),
		testUtils.MockWithRoadsByPlayer(map[string][]int{
			"1": {1, 2, 7},
		}),
	)

	t.Run("player has road building card, will try to build roads in available/connected edges", func(t *testing.T) {
		err := game.UseRoadBuilding("1")
		if err != nil {
			t.Errorf("expected to use road building card just fine, but actually got error %s", err.Error())
		}
		if game.RoundType() != testUtils.BuildRoad1Development {
			t.Errorf("expected round type to be %s after road building use, but got %s", testUtils.RoundTypeTranslation[testUtils.BuildRoad1Development], testUtils.RoundTypeTranslation[game.RoundType()])
		}

		err = game.PickRoadBuildingSpot("1", 8)
		if err != nil {
			t.Errorf("expected to build road on edge#8 just fine, but actually got error %s", err.Error())
		}
		if game.RoundType() != testUtils.BuildRoad2Development {
			t.Errorf("expected round type to be %s after 1st road building spot picked, but got %s", testUtils.RoundTypeTranslation[testUtils.BuildRoad2Development], testUtils.RoundTypeTranslation[game.RoundType()])
		}

		err = game.PickRoadBuildingSpot("1", 3)
		if err != nil {
			t.Errorf("expected to build road on edge#3 just fine, but actually got error %s", err.Error())
		}
		if game.RoundType() != testUtils.Regular {
			t.Errorf("expected round type to be %s after 2nd road building spot picked, but got %s", testUtils.RoundTypeTranslation[testUtils.Regular], testUtils.RoundTypeTranslation[game.RoundType()])
		}
	})
}

func TestPlayRoadBuildingDevelopmentCardUnavailablePathAvailableRoads(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {1},
		}),
		testUtils.MockWithDevelopmentsByPlayer(map[string]map[string]int{
			"1": {
				"Road Building": 1,
			},
		}),
		testUtils.MockWithRoadsByPlayer(map[string][]int{
			"1": {1, 2, 7},
			"2": {8},
		}),
	)

	t.Run("player has road building card, will try to build roads in unavailable/unconnected edges", func(t *testing.T) {
		err := game.UseRoadBuilding("1")
		if err != nil {
			t.Errorf("expected to use road building card just fine, but actually got error %s", err.Error())
		}
		if game.RoundType() != testUtils.BuildRoad1Development {
			t.Errorf("expected round type to be %s after road building use, but got %s", testUtils.RoundTypeTranslation[testUtils.BuildRoad1Development], testUtils.RoundTypeTranslation[game.RoundType()])
		}

		err = game.PickRoadBuildingSpot("1", 2)
		if err == nil {
			t.Errorf("expected to not be able to road building pick edge#2, but it picked just fine")
		}
		if game.RoundType() != testUtils.BuildRoad1Development {
			t.Errorf("expected round type to be %s after road building pick failure, but got %s", testUtils.RoundTypeTranslation[testUtils.BuildRoad1Development], testUtils.RoundTypeTranslation[game.RoundType()])
		}

		err = game.PickRoadBuildingSpot("1", 8)
		if err == nil {
			t.Errorf("expected to not be able to road building pick edge#8, but it picked just fine")
		}
		if game.RoundType() != testUtils.BuildRoad1Development {
			t.Errorf("expected round type to be %s after road building pick failure, but got %s", testUtils.RoundTypeTranslation[testUtils.BuildRoad1Development], testUtils.RoundTypeTranslation[game.RoundType()])
		}

		err = game.PickRoadBuildingSpot("1", 72)
		if err == nil {
			t.Errorf("expected to not be able to road building pick edge#72, but it picked just fine")
		}
		if game.RoundType() != testUtils.BuildRoad1Development {
			t.Errorf("expected round type to be %s after road building pick failure, but got %s", testUtils.RoundTypeTranslation[testUtils.BuildRoad1Development], testUtils.RoundTypeTranslation[game.RoundType()])
		}

		err = game.PickRoadBuildingSpot("1", 3)
		if err != nil {
			t.Errorf("expected to build road on edge#3 just fine, but actually got error %s", err.Error())
		}
		if game.RoundType() != testUtils.BuildRoad2Development {
			t.Errorf("expected round type to be %s after 1st road building spot picked, but got %s", testUtils.RoundTypeTranslation[testUtils.BuildRoad2Development], testUtils.RoundTypeTranslation[game.RoundType()])
		}

		err = game.PickRoadBuildingSpot("1", 3)
		if err == nil {
			t.Errorf("expected to not be able to road building pick edge#3, but it picked just fine")
		}
		if game.RoundType() != testUtils.BuildRoad2Development {
			t.Errorf("expected round type to be %s after road building pick failure, but got %s", testUtils.RoundTypeTranslation[testUtils.BuildRoad2Development], testUtils.RoundTypeTranslation[game.RoundType()])
		}

		err = game.PickRoadBuildingSpot("1", 8)
		if err == nil {
			t.Errorf("expected to not be able to road building pick edge#8, but it picked just fine")
		}
		if game.RoundType() != testUtils.BuildRoad2Development {
			t.Errorf("expected round type to be %s after road building pick failure, but got %s", testUtils.RoundTypeTranslation[testUtils.BuildRoad2Development], testUtils.RoundTypeTranslation[game.RoundType()])
		}

		err = game.PickRoadBuildingSpot("1", 72)
		if err == nil {
			t.Errorf("expected to not be able to road building pick edge#72, but it picked just fine")
		}
		if game.RoundType() != testUtils.BuildRoad2Development {
			t.Errorf("expected round type to be %s after road building pick failure, but got %s", testUtils.RoundTypeTranslation[testUtils.BuildRoad2Development], testUtils.RoundTypeTranslation[game.RoundType()])
		}

		err = game.PickRoadBuildingSpot("1", 4)
		if err != nil {
			t.Errorf("expected to build road on edge#4 just fine, but actually got error %s", err.Error())
		}
		if game.RoundType() != testUtils.Regular {
			t.Errorf("expected round type to be %s after 2nd road building spot picked, but got %s", testUtils.RoundTypeTranslation[testUtils.Regular], testUtils.RoundTypeTranslation[game.RoundType()])
		}
	})
}

func TestPlayRoadBuildingDevelopmentCardAvailablePathUnavailableRoads(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {1},
		}),
		testUtils.MockWithDevelopmentsByPlayer(map[string]map[string]int{
			"1": {
				"Road Building": 1,
			},
		}),
		testUtils.MockWithRoadsByPlayer(map[string][]int{
			"1": {1, 2, 7, 24, 23, 22, 18, 17, 35, 36, 32, 53, 50, 66, 63, 64, 67, 68, 70, 71},
			"2": {8},
		}),
	)

	t.Run("player has road building card, will try to build roads in available edges but no available roads", func(t *testing.T) {
		err := game.UseRoadBuilding("1")
		if err == nil {
			t.Errorf("expected to not be allowed to use road building card, but used just fine")
		}
	})
}

func TestPlayRoadBuildingDevelopmentCardAvailablePathOnly1AvailableRoad(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {1},
		}),
		testUtils.MockWithDevelopmentsByPlayer(map[string]map[string]int{
			"1": {
				"Road Building": 1,
			},
		}),
		testUtils.MockWithRoadsByPlayer(map[string][]int{
			"1": {1, 2, 7, 24, 23, 22, 18, 17, 35, 36, 32, 53, 50, 66, 63, 64, 67, 68, 70},
			"2": {8},
		}),
	)

	t.Run("player has road building card, will try to build roads in available edges but only 1 available road", func(t *testing.T) {
		err := game.UseRoadBuilding("1")
		if err != nil {
			t.Errorf("expected to use road building card just fine, but actually got error %s", err.Error())
		}
		if game.RoundType() != testUtils.BuildRoad1Development {
			t.Errorf("expected round type to be %s after road building use, but got %s", testUtils.RoundTypeTranslation[testUtils.BuildRoad1Development], testUtils.RoundTypeTranslation[game.RoundType()])
		}

		err = game.PickRoadBuildingSpot("1", 3)
		if err != nil {
			t.Errorf("expected to build road on edge#3 just fine, but actually got error %s", err.Error())
		}
		if game.RoundType() != testUtils.Regular {
			t.Errorf("expected round type to be %s after 1st road building spot picked that reached max roads, but got %s", testUtils.RoundTypeTranslation[testUtils.Regular], testUtils.RoundTypeTranslation[game.RoundType()])
		}
	})
}

func TestPlayRoadBuildingDevelopmentCardWithoutHavingOne(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {1},
		}),
		testUtils.MockWithDevelopmentsByPlayer(map[string]map[string]int{
			"1": {
				"Road Building": 0,
			},
		}),
	)

	t.Run("player has no road building card, will try to play one", func(t *testing.T) {
		err := game.UseRoadBuilding("1")
		if err == nil {
			t.Errorf("expected to not be able to use road building, but actually no error was found")
		}
	})
}

func TestPlayRoadBuildingDevelopmentCardByRound(t *testing.T) {
	createGame := func(roundType int) *testUtils.GameState {
		game := testUtils.CreateTestGame(
			testUtils.MockWithRoundType(roundType),
			testUtils.MockWithDevelopmentsByPlayer(map[string]map[string]int{
				"1": {
					"Road Building": 1,
				},
			}),
		)
		return game
	}

	willHaveErrorByRoundType := map[int]bool{
		testUtils.SetupSettlement1:          true,
		testUtils.SetupRoad1:                true,
		testUtils.SetupSettlement2:          true,
		testUtils.SetupRoad2:                true,
		testUtils.FirstRound:                false,
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
			err := game.UseRoadBuilding("1")
			hasErr := err != nil
			if hasErr != willHaveError {
				t.Errorf("expected error to be %v, but actually was %v", willHaveError, hasErr)
			}
		})
	}
}

func TestPlayYearOfPlentyDevelopmentCardAvailableResources(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {1},
		}),
		testUtils.MockWithDevelopmentsByPlayer(map[string]map[string]int{
			"1": {
				"Year of Plenty": 1,
			},
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
	)

	t.Run("player has year of plenty card, will try to acquire available resources", func(t *testing.T) {
		err := game.UseYearOfPlenty("1")
		if err != nil {
			t.Errorf("expected to use year of plenty card just fine, but actually got error %s", err.Error())
		}
		if game.RoundType() != testUtils.YearOfPlentyPickResources {
			t.Errorf("expected round type to be %s after year of plenty use, but got %s", testUtils.RoundTypeTranslation[testUtils.YearOfPlentyPickResources], testUtils.RoundTypeTranslation[game.RoundType()])
		}

		err = game.PickYearOfPlentyResources("1", "Grain", "Ore")
		if err != nil {
			t.Errorf("expected to receive the two selected resources just fine, but actually got error %s", err.Error())
		}
		if game.RoundType() != testUtils.Regular {
			t.Errorf("expected round type to be %s after picking monopoly resource, but it's actually %s", testUtils.RoundTypeTranslation[testUtils.Regular], testUtils.RoundTypeTranslation[game.RoundType()])
		}
	})
}

func TestPlayYearOfPlentyDevelopmentCardWithoutHavingOne(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithSettlementsByPlayer(map[string][]int{
			"1": {1},
		}),
		testUtils.MockWithDevelopmentsByPlayer(map[string]map[string]int{
			"1": {
				"Year of Plenty": 0,
			},
		}),
	)

	t.Run("player has no monopoly card, will try to play one", func(t *testing.T) {
		err := game.UseYearOfPlenty("1")
		if err == nil {
			t.Errorf("expected to not be able to use year of plenty card, but actually no error was found")
		}
	})
}

func TestPlayYearOfPlentyDevelopmentCardByRound(t *testing.T) {
	createGame := func(roundType int) *testUtils.GameState {
		game := testUtils.CreateTestGame(
			testUtils.MockWithRoundType(roundType),
			testUtils.MockWithDevelopmentsByPlayer(map[string]map[string]int{
				"1": {
					"Year of Plenty": 1,
				},
			}),
		)
		return game
	}

	willHaveErrorByRoundType := map[int]bool{
		testUtils.SetupSettlement1:          true,
		testUtils.SetupRoad1:                true,
		testUtils.SetupSettlement2:          true,
		testUtils.SetupRoad2:                true,
		testUtils.FirstRound:                false,
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
			err := game.UseYearOfPlenty("1")
			hasErr := err != nil
			if hasErr != willHaveError {
				t.Errorf("expected error to be %v, but actually was %v", willHaveError, hasErr)
			}
		})
	}
}
