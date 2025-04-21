package core

import (
	"fmt"
	"maps"
	"testing"
)

func TestRollDiceNotPlayerRound(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(BetweenTurns),
		MockWithCurrentRoundPlayer("2"),
	)
	t.Run("dice roll attempt out of player's round", func(t *testing.T) {
		err := game.RollDice("1")
		if err == nil {
			t.Errorf("expected to have error while trying to roll dice not being player's round, but rolled just fine")
		}
	})
}

func TestRollDiceAcrossAllRoundTypes(t *testing.T) {
	createGame := func(roundType int) *GameState {
		game := CreateTestGame(
			MockWithRoundType(roundType),
		)
		return game
	}

	willHaveErrorByRoundType := map[int]bool{
		SetupSettlement1: true,
		SetupRoad1:       true,
		SetupSettlement2: true,
		SetupRoad2:       true,
		FirstRound:       false,
		Regular:          true,
		MoveRobberDue7:   true,
		PickRobbed:       true,
		BetweenTurns:     false,
		DiscardPhase:     true,
	}

	for roundType, willHaveError := range willHaveErrorByRoundType {
		testname := fmt.Sprintf("round type: %s, will have error: %v", RoundTypeTranslation[roundType], willHaveError)
		t.Run(testname, func(t *testing.T) {
			game := createGame(roundType)
			err := game.RollDice("1")
			hasErr := err != nil
			if hasErr != willHaveError {
				t.Errorf("expected error to be %v, but actually was %v", willHaveError, hasErr)
			}
		})
	}
}

func TestHandleDiceNot7NotBlocked(t *testing.T) {
	createGame := func(sum int) *GameState {
		rand := StubRand(sum)
		game := CreateTestGame(
			MockWithRoundType(BetweenTurns),
			MockWithSettlementsByPlayer(map[string][]int{
				"1": {40, 11, 6},
			}),
			MockWithCitiesByPlayer(map[string][]int{
				"1": {32},
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
		return game
	}

	var tests = []struct {
		dice           int
		resourcesAfter map[string]int
	}{
		{
			dice: 2,
			resourcesAfter: map[string]int{
				"Lumber": 1,
				"Brick":  1,
				"Sheep":  3,
				"Grain":  1,
				"Ore":    1,
			},
		},
		{
			dice: 3,
			resourcesAfter: map[string]int{
				"Lumber": 1,
				"Brick":  1,
				"Sheep":  1,
				"Grain":  2,
				"Ore":    1,
			},
		},
		{
			dice: 4,
			resourcesAfter: map[string]int{
				"Lumber": 3,
				"Brick":  2,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
		},
		{
			dice: 5,
			resourcesAfter: map[string]int{
				"Lumber": 1,
				"Brick":  1,
				"Sheep":  1,
				"Grain":  2,
				"Ore":    1,
			},
		},
		{
			dice: 6,
			resourcesAfter: map[string]int{
				"Lumber": 1,
				"Brick":  1,
				"Sheep":  2,
				"Grain":  1,
				"Ore":    1,
			},
		},
		{
			dice: 8,
			resourcesAfter: map[string]int{
				"Lumber": 1,
				"Brick":  1,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
		},
		{
			dice: 9,
			resourcesAfter: map[string]int{
				"Lumber": 2,
				"Brick":  1,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
		},
		{
			dice: 10,
			resourcesAfter: map[string]int{
				"Lumber": 1,
				"Brick":  3,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
		},
		{
			dice: 11,
			resourcesAfter: map[string]int{
				"Lumber": 1,
				"Brick":  1,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    2,
			},
		},
		{
			dice: 12,
			resourcesAfter: map[string]int{
				"Lumber": 1,
				"Brick":  1,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("dice %d", tt.dice)
		t.Run(testname, func(t *testing.T) {
			game := createGame(tt.dice)
			game.RollDice("1")
			actualResources := game.ResourceHandByPlayer("1")
			if !maps.Equal(tt.resourcesAfter, actualResources) {
				t.Errorf("expected %v, got %v", tt.resourcesAfter, actualResources)
			}
		})
	}
}

func TestHandleDiceNot7Blocked(t *testing.T) {
	rand := StubRand(4)
	game := CreateTestGame(
		MockWithRoundType(BetweenTurns),
		MockWithSettlementsByPlayer(map[string][]int{
			"1": {40, 11, 6},
		}),
		MockWithCitiesByPlayer(map[string][]int{
			"1": {32},
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
		MockWithBlockedTile(10),
		MockWithRand(rand),
	)

	t.Run("dice 4, tileID 10 blocked", func(t *testing.T) {
		game.RollDice("1")
		actualResources := game.ResourceHandByPlayer("1")

		if actualResources["Lumber"] != 1 {
			t.Errorf("expected to have 1 Lumber, actually have %d", actualResources["Lumber"])
		}
		if actualResources["Brick"] != 2 {
			t.Errorf("expected to have 2 Brick, actually have %d", actualResources["Brick"])
		}
		if actualResources["Sheep"] != 1 {
			t.Errorf("expected to have 1 Sheep, actually have %d", actualResources["Sheep"])
		}
		if actualResources["Grain"] != 1 {
			t.Errorf("expected to have 1 Grain, actually have %d", actualResources["Grain"])
		}
		if actualResources["Ore"] != 1 {
			t.Errorf("expected to have 1 Ore, actually have %d", actualResources["Ore"])
		}
	})
}
