package tests

import (
	"fmt"
	"testing"

	testUtils "github.com/victoroliveirab/settlers/core/state"
	"github.com/victoroliveirab/settlers/utils"
)

// TODO: will have to add more tests when FriendlyRobber is active
func TestRobbablePlayers(t *testing.T) {
	createGame := func(settlementMap, cityMap map[string][]int) *testUtils.GameState {
		game := testUtils.CreateTestGame(
			testUtils.MockWithRoundType(testUtils.PickRobbed),
			testUtils.MockWithBlockedTile(17),
			testUtils.MockWithSettlementsByPlayer(settlementMap),
			testUtils.MockWithCitiesByPlayer(cityMap),
		)
		return game
	}

	var tests = []struct {
		description    string
		cityMap        map[string][]int
		settlementMap  map[string][]int
		expectedResult []string
	}{
		{
			description: "single owner, but not current player",
			cityMap: map[string][]int{
				"1": {1},
				"2": {42},
			},
			settlementMap:  map[string][]int{},
			expectedResult: []string{"2"},
		},
		{
			description: "single owner, current player",
			cityMap: map[string][]int{
				"1": {42},
				"2": {1},
			},
			settlementMap:  map[string][]int{},
			expectedResult: []string{},
		},
		{
			description: "multiple owners of tile, including current player",
			cityMap: map[string][]int{
				"1": {49},
				"2": {42},
			},
			settlementMap: map[string][]int{
				"3": {40},
			},
			expectedResult: []string{"2", "3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			game := createGame(tt.settlementMap, tt.cityMap)

			robbablePlayers, err := game.RobbablePlayers("1")
			if err != nil {
				t.Errorf("expected to run game.RobbablePlayers(\"1\") just fine, but actually got error %s", err.Error())
			}
			robbablePlayersSet := utils.SetFromSlice(robbablePlayers)
			expectedResultSet := utils.SetFromSlice(tt.expectedResult)
			if !robbablePlayersSet.Equal(expectedResultSet) {
				t.Errorf("expected robbable players to be %v, but it's actually %v", tt.expectedResult, robbablePlayers)
			}
		})
	}
}

func TestDiscardAmountByPlayer(t *testing.T) {
	createGame := func(resourceMap map[string]map[string]int) *testUtils.GameState {
		game := testUtils.CreateTestGame(
			testUtils.MockWithRoundType(testUtils.DiscardPhase),
			testUtils.MockWithResourcesByPlayer(resourceMap),
		)
		return game
	}

	var tests = []struct {
		description    string
		resourceMap    map[string]map[string]int
		expectedResult map[string]int
	}{
		{
			description: "no player has more than 7 cards",
			resourceMap: map[string]map[string]int{
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
			},
			expectedResult: map[string]int{
				"1": 0,
				"2": 0,
				"3": 0,
				"4": 0,
			},
		},
		{
			description: "multiple player have more than 7 cards",
			resourceMap: map[string]map[string]int{
				"1": {
					"Lumber": 1,
					"Brick":  2,
					"Sheep":  3,
					"Grain":  4,
					"Ore":    5, // total: 15
				},
				"2": {
					"Lumber": 1,
					"Brick":  2,
					"Sheep":  3,
					"Grain":  4,
					"Ore":    4, // total: 14
				},
				"3": {
					"Lumber": 1,
					"Brick":  2,
					"Sheep":  3,
					"Grain":  4,
					"Ore":    3, // total: 13
				},
			},
			expectedResult: map[string]int{
				"1": 7,
				"2": 7,
				"3": 6,
				"4": 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			game := createGame(tt.resourceMap)
			discardAmountMap := make(map[string]int)
			for i := 1; i <= 4; i++ {
				id := fmt.Sprintf("%d", i)
				discardAmountMap[id] = game.DiscardAmountByPlayer(id)
			}

			for id, expectedResult := range tt.expectedResult {
				if expectedResult != discardAmountMap[id] {
					t.Errorf("expected player %s to have to discard %d cards, but actually is %d", id, expectedResult, discardAmountMap[id])
				}
			}
		})
	}

}
