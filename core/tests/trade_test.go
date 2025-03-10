package tests

import (
	"fmt"
	"testing"

	testUtils "github.com/victoroliveirab/settlers/core"
)

func TestTradeWithBankWithAvailableResources(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
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
		err := game.MakeBankTrade("1", "Lumber", "Ore")
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
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
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
		err := game.MakeBankTrade("1", "Lumber", "Ore")
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
	createGame := func(roundType int) *testUtils.GameState {
		game := testUtils.CreateTestGame(
			testUtils.MockWithRoundType(roundType),
			testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
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
			err := game.MakeBankTrade("1", "Lumber", "Ore")
			hasErr := err != nil
			if hasErr != willHaveError {
				t.Errorf("expected error to be %v, but actually was %v", willHaveError, hasErr)
			}
		})
	}
}

func TestTradeWithBankNotPlayerRound(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
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
		err := game.MakeBankTrade("2", "Lumber", "Ore")
		if err == nil {
			t.Errorf("expected to not be able to trade with bank during other player's round, but traded just fine")
		}
	})
}

func TestTradeWithPortGeneralWithAvailableResources(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
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

	t.Run("trade with general port - player has available resources", func(t *testing.T) {
		var vertexID int
		for vertexID = range game.AllSettlements() {
			break
		}
		err := game.MakePortTrade("1", vertexID, "Lumber", "Ore")
		if err != nil {
			t.Errorf("expected to trade with port just fine, but actually got error %s", err.Error())
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

func TestTradeWithPortGeneralWithNoAvailableResources(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 2,
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

	t.Run("trade with port - player doesn't have available resources", func(t *testing.T) {
		var vertexID int
		for vertexID = range game.AllSettlements() {
			break
		}
		err := game.MakePortTrade("1", vertexID, "Lumber", "Ore")
		if err == nil {
			t.Errorf("expected to not be able to trade with port, but actually traded just fine")
		}

		player1Resources := game.ResourceHandByPlayer("1")
		if player1Resources["Lumber"] != 2 {
			t.Errorf("expected player#1 to have 2 Lumber, actually got %d", player1Resources["Lumber"])
		}
		if player1Resources["Ore"] != 1 {
			t.Errorf("expected player#1 to have 1 Ore, actually got %d", player1Resources["Ore"])
		}
	})
}

func TestTradeWithPortSpecificWithAvailableResources(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 2,
				"Brick":  1,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
		}),
		testUtils.MockWithPortsByPlayer(map[string][]string{
			"1": {"Lumber"},
		}),
	)

	t.Run("trade with specific port - player has available resources", func(t *testing.T) {
		var vertexID int
		for vertexID = range game.AllSettlements() {
			break
		}
		err := game.MakePortTrade("1", vertexID, "Lumber", "Ore")
		if err != nil {
			t.Errorf("expected to trade with port just fine, but actually got error %s", err.Error())
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

func TestTradeWithPortSpecificWithNoAvailableResources(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 1,
				"Brick":  1,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
		}),
		testUtils.MockWithPortsByPlayer(map[string][]string{
			"1": {"Lumber"},
		}),
	)

	t.Run("trade with specific port - player doesn't have available resources", func(t *testing.T) {
		var vertexID int
		for vertexID = range game.AllSettlements() {
			break
		}
		err := game.MakePortTrade("1", vertexID, "Lumber", "Ore")
		if err == nil {
			t.Errorf("expected to not be able to trade with port, but actually traded just fine")
		}

		player1Resources := game.ResourceHandByPlayer("1")
		if player1Resources["Lumber"] != 1 {
			t.Errorf("expected player#1 to have 1 Lumber, actually got %d", player1Resources["Lumber"])
		}
		if player1Resources["Ore"] != 1 {
			t.Errorf("expected player#1 to have 1 Ore, actually got %d", player1Resources["Ore"])
		}
	})
}

func TestTradeWithPortSpecificPortIsOtherResource(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 2,
				"Brick":  1,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
		}),
		testUtils.MockWithPortsByPlayer(map[string][]string{
			"1": {"Brick"},
		}),
	)

	t.Run("trade with specific port - player has different available resources", func(t *testing.T) {
		var vertexID int
		for vertexID = range game.AllSettlements() {
			break
		}
		err := game.MakePortTrade("1", vertexID, "Lumber", "Ore")
		if err == nil {
			t.Errorf("expected to not be able to trade with port, but actually traded just fine")
		}

		player1Resources := game.ResourceHandByPlayer("1")
		if player1Resources["Lumber"] != 2 {
			t.Errorf("expected player#1 to have 2 Lumber, actually got %d", player1Resources["Lumber"])
		}
		if player1Resources["Ore"] != 1 {
			t.Errorf("expected player#1 to have 1 Ore, actually got %d", player1Resources["Ore"])
		}
	})
}

func TestTradeWithPortByRound(t *testing.T) {
	createGame := func(roundType int) *testUtils.GameState {
		game := testUtils.CreateTestGame(
			testUtils.MockWithRoundType(roundType),
			testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
				"1": {
					"Lumber": 2,
					"Brick":  1,
					"Sheep":  1,
					"Grain":  1,
					"Ore":    1,
				},
			}),
			testUtils.MockWithPortsByPlayer(map[string][]string{
				"1": {"Lumber"},
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
			var vertexID int
			for vertexID = range game.AllSettlements() {
				break
			}
			err := game.MakePortTrade("1", vertexID, "Lumber", "Ore")
			hasErr := err != nil
			if hasErr != willHaveError {
				t.Errorf("expected error to be %v, but actually was %v", willHaveError, hasErr)
			}
		})
	}
}

func TestTradeWithPortNotPlayerRound(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"2": {
				"Lumber": 2,
				"Brick":  1,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
		}),
		testUtils.MockWithPortsByPlayer(map[string][]string{
			"2": {"Lumber"},
		}),
	)

	t.Run("trade with port - not player's round", func(t *testing.T) {
		var vertexID int
		for vertexID = range game.AllSettlements() {
			break
		}
		err := game.MakePortTrade("2", vertexID, "Lumber", "Ore")
		if err == nil {
			t.Errorf("expected to not be able to trade with port during other player's round, but traded just fine")
		}
	})
}

func TestCreateTradeOfferWithAvailableResources(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
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

	t.Run("create a trade offer - player has available resources", func(t *testing.T) {
		_, err := game.MakeTradeOffer("1", map[string]int{
			"Lumber": 1,
			"Brick":  1,
		}, map[string]int{
			"Ore": 1,
		}, []string{})
		if err != nil {
			t.Errorf("expected to make trade offer just fine, but actually got error %s", err.Error())
		}

		trade := game.ActiveTradeOffers()[0]
		if trade.PlayerID != "1" {
			t.Errorf("expected trade to belong to player#1, but actually belongs to player#%s", trade.PlayerID)
		}
		if trade.Status != "Open" {
			t.Errorf("expected trade offer status to be \"Open\", but actually got %s", trade.Status)
		}
		for id, opponent := range trade.Opponents {
			if opponent.Status != "Open" {
				t.Errorf("expected opponent#%s trade status to be \"Open\", but actually got %s", id, opponent.Status)
			}
		}
		if trade.Finalized {
			t.Errorf("expected trade offer not to be finalized, but actually got that it is")
		}
		if trade.ParentID != -1 {
			t.Errorf("expected parentID to be -1, but actually got %d", trade.ParentID)
		}
	})
}

func TestCreateTradeOfferWithNoAvailableResources(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 1,
				"Brick":  0,
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

	t.Run("create a trade offer - player hasn't all available resources", func(t *testing.T) {
		_, err := game.MakeTradeOffer("1", map[string]int{
			"Lumber": 1,
			"Brick":  1,
		}, map[string]int{
			"Ore": 1,
		}, []string{})
		if err == nil {
			t.Errorf("expected not to be able to make trade offer, but actually no error was found")
		}

		allTrades := game.ActiveTradeOffers()
		if len(allTrades) != 0 {
			t.Errorf("expected to not have active trades, actually got %d", len(allTrades))
		}
	})
}

func TestCreateTradeOfferByRound(t *testing.T) {
	createGame := func(roundType int) *testUtils.GameState {
		game := testUtils.CreateTestGame(
			testUtils.MockWithRoundType(roundType),
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
			_, err := game.MakeTradeOffer("1",
				map[string]int{
					"Lumber": 1,
				},
				map[string]int{
					"Ore": 1,
				},
				[]string{})
			hasErr := err != nil
			if hasErr != willHaveError {
				t.Errorf("expected error to be %v, but actually was %v", willHaveError, hasErr)
			}
		})
	}
}

func TestCreateTradeOfferNotPlayerRound(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
		testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
			"2": {
				"Lumber": 1,
				"Brick":  1,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
		}),
	)

	t.Run("create trade offer - not player's round", func(t *testing.T) {
		_, err := game.MakeTradeOffer("2", map[string]int{
			"Lumber": 1,
		}, map[string]int{
			"Brick": 1,
		}, []string{})
		if err == nil {
			t.Errorf("expected to not be able to create trade offer during other player's round, but traded just fine")
		}
	})
}

func TestCreateCounterTradeOfferWithAvailableResources(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
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

	t.Run("create a counter trade offer - player has available resources", func(t *testing.T) {
		tradeID, err := game.MakeTradeOffer("1", map[string]int{
			"Lumber": 1,
			"Brick":  1,
		}, map[string]int{
			"Ore": 1,
		}, []string{})
		if err != nil {
			t.Errorf("expected to make trade offer just fine, but actually got error %s", err.Error())
		}

		_, err = game.MakeCounterTradeOffer("2", tradeID, map[string]int{
			"Ore": 1,
		}, map[string]int{
			"Lumber": 2,
			"Brick":  1,
		})
		if err != nil {
			t.Errorf("expected to make counter trade offer just fine, but actually got error %s", err.Error())
		}

		allTrades := game.ActiveTradeOffers()
		if len(allTrades) != 2 {
			t.Errorf("expected to have 2 active trade offers, but actually got %d", len(allTrades))
		}

		var trade1, trade2 testUtils.Trade
		if allTrades[0].ID == tradeID {
			trade1 = allTrades[0]
			trade2 = allTrades[1]
		} else {
			trade1 = allTrades[1]
			trade2 = allTrades[0]
		}

		if len(trade1.Counters) != 1 {
			t.Log(trade1)
			t.Log(trade2)
			t.Log(game.ActiveTradeOffers())
			t.Errorf("expected original trade offer to have 1 counter offer, but actually has length %d", len(trade1.Counters))
		}
		if trade1.Counters[0] != trade2.ID {
			t.Errorf("expected original trade offer to have counter offer with id %d, but actually got %d", trade2.ID, trade1.Counters[0])
		}
		if trade2.ParentID != trade1.ID {
			t.Errorf("expected counter trade offer to have parentID %d, but actually got %d", trade1.ID, trade2.ParentID)
		}
		if trade2.Opponents["2"].Status != "Accepted" {
			t.Errorf("expected counter trade offer to be already accepted by creator, but actually got status %s", trade2.Opponents["2"].Status)
		}
	})
}

func TestCreateCounterTradeOfferWithNoAvailableResources(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
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

	t.Run("create a counter trade offer - player hasn't enough available resources", func(t *testing.T) {
		tradeID, err := game.MakeTradeOffer("1", map[string]int{
			"Lumber": 1,
			"Brick":  1,
		}, map[string]int{
			"Ore": 1,
		}, []string{})
		if err != nil {
			t.Errorf("expected to make trade offer just fine, but actually got error %s", err.Error())
		}

		_, err = game.MakeCounterTradeOffer("2", tradeID, map[string]int{
			"Ore": 2,
		}, map[string]int{
			"Lumber": 2,
			"Brick":  1,
		})
		if err == nil {
			t.Errorf("expected not to be able to make counter trade offer, but actually no error was found")
		}

		allTrades := game.ActiveTradeOffers()
		if len(allTrades) != 1 {
			t.Errorf("expected to have 1 active trade offer, but actually got %d", len(allTrades))
		}

		trade := allTrades[0]
		if len(trade.Counters) != 0 {
			t.Errorf("expected to have no counter offers, but actually has %d", len(trade.Counters))
		}
		if trade.Opponents["2"].Status != "Open" {
			t.Errorf("expected to have player#2 status still \"Open\" after failed counter offer creation, but actually has %s", trade.Opponents["2"].Status)
		}
	})
}

func TestCreateCounterTradeOfferByRound(t *testing.T) {
	createGame := func(roundType int) *testUtils.GameState {
		game := testUtils.CreateTestGame(
			testUtils.MockWithRoundType(roundType),
			testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
				"1": {
					"Lumber": 1,
					"Brick":  1,
					"Sheep":  1,
					"Grain":  1,
					"Ore":    1,
				},
				"2": {
					"Ore": 1,
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
			tradeID, _ := game.MakeTradeOffer("1",
				map[string]int{
					"Lumber": 1,
				},
				map[string]int{
					"Ore": 1,
				},
				[]string{})
			_, err := game.MakeCounterTradeOffer("2", tradeID,
				map[string]int{
					"Ore": 1,
				},
				map[string]int{
					"Lumber": 2,
				},
			)
			hasErr := err != nil
			if hasErr != willHaveError {
				t.Errorf("expected error to be %v, but actually was %v", willHaveError, hasErr)
				t.Error(err.Error())
			}
		})
	}
}

func TestCreateCounterTradeOfferOwnRound(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
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

	t.Run("create a counter trade offer - counter offer to own offer", func(t *testing.T) {
		tradeID, err := game.MakeTradeOffer("1", map[string]int{
			"Lumber": 1,
			"Brick":  1,
		}, map[string]int{
			"Ore": 1,
		}, []string{})
		if err != nil {
			t.Errorf("expected to make trade offer just fine, but actually got error %s", err.Error())
		}

		_, err = game.MakeCounterTradeOffer("1", tradeID, map[string]int{
			"Ore": 1,
		}, map[string]int{
			"Lumber": 2,
			"Brick":  1,
		})
		if err == nil {
			t.Errorf("expected to not be able to create counter trade offer to own trade offer, but it actually made just fine")
		}
	})
}

func TestAcceptTradeOfferWithAvailableResources(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
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

	t.Run("accept a trade offer - players have the resources available", func(t *testing.T) {
		tradeID, err := game.MakeTradeOffer("1", map[string]int{
			"Lumber": 1,
			"Brick":  1,
		}, map[string]int{
			"Ore": 1,
		}, []string{})
		if err != nil {
			t.Errorf("expected to make trade offer just fine, but actually got error %s", err.Error())
		}

		err = game.AcceptTradeOffer("2", tradeID)
		if err != nil {
			t.Errorf("expected to accept offer just fine, but actually got error %s", err.Error())
		}

		trade := game.ActiveTradeOffers()[0]
		if trade.Opponents["2"].Status != "Accepted" {
			t.Errorf("expected player#2 state to be \"Accepted\", but actually got %s", trade.Opponents["2"].Status)
		}
		if trade.Opponents["3"].Status != "Open" {
			t.Errorf("expected player#3 state to be \"Open\", but actually got %s", trade.Opponents["3"].Status)
		}
		if trade.Opponents["4"].Status != "Open" {
			t.Errorf("expected player#4 state to be \"Open\", but actually got %s", trade.Opponents["4"].Status)
		}
	})
}

func TestAcceptTradeOfferWithNoAvailableResources(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
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

	t.Run("accept a trade offer - player has not the resources available", func(t *testing.T) {
		tradeID, err := game.MakeTradeOffer("1", map[string]int{
			"Lumber": 1,
			"Brick":  1,
		}, map[string]int{
			"Ore": 2,
		}, []string{})
		if err != nil {
			t.Errorf("expected to make trade offer just fine, but actually got error %s", err.Error())
		}

		err = game.AcceptTradeOffer("2", tradeID)
		if err == nil {
			t.Errorf("expected not to be able to accept trade offer, but actually no error was found")
		}

		trade := game.ActiveTradeOffers()[0]
		if trade.Opponents["2"].Status != "Declined" {
			t.Errorf("expected player#2 status to have changed to \"Declined\" after trying to accept without enough resources, but actually got %s", trade.Opponents["2"].Status)
		}
	})
}

func TestAcceptTradeOfferByRound(t *testing.T) {
	createGame := func(roundType int) *testUtils.GameState {
		game := testUtils.CreateTestGame(
			testUtils.MockWithRoundType(roundType),
			testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
				"1": {
					"Lumber": 1,
					"Brick":  1,
					"Sheep":  1,
					"Grain":  1,
					"Ore":    1,
				},
				"2": {
					"Ore": 1,
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
			tradeID, _ := game.MakeTradeOffer("1",
				map[string]int{
					"Lumber": 1,
				},
				map[string]int{
					"Ore": 1,
				},
				[]string{})
			err := game.AcceptTradeOffer("2", tradeID)
			hasErr := err != nil
			if hasErr != willHaveError {
				t.Errorf("expected error to be %v, but actually was %v", willHaveError, hasErr)
			}
		})
	}
}

func TestAcceptTradeOfferOwnOffer(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
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

	t.Run("accept a trade offer - accept own offer", func(t *testing.T) {
		tradeID, err := game.MakeTradeOffer("1", map[string]int{
			"Lumber": 1,
			"Brick":  1,
		}, map[string]int{
			"Ore": 1,
		}, []string{})
		if err != nil {
			t.Errorf("expected to make trade offer just fine, but actually got error %s", err.Error())
		}

		err = game.AcceptTradeOffer("1", tradeID)
		if err == nil {
			t.Errorf("expected to not be able to accept own trade offer, but it actually accepted just fine")
		}
	})
}

func TestFinalizeAcceptedTradeOfferWithAvailableResources(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
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

	t.Run("finalize a regular trade offer - players have the resources available", func(t *testing.T) {
		tradeID, err := game.MakeTradeOffer("1", map[string]int{
			"Lumber": 1,
			"Brick":  1,
		}, map[string]int{
			"Ore": 1,
		}, []string{})
		if err != nil {
			t.Errorf("expected to make trade offer just fine, but actually got error %s", err.Error())
		}

		err = game.AcceptTradeOffer("2", tradeID)
		if err != nil {
			t.Errorf("expected to accept offer just fine, but actually got error %s", err.Error())
		}

		allTrades := game.ActiveTradeOffers()
		if len(allTrades) == 0 {
			t.Errorf("expected trade offer not to be finalized, but actually got that it is")
		}

		err = game.FinalizeTrade("1", "2", tradeID)
		if err != nil {
			t.Errorf("expected to finalize offer just fine, but actually got error %s", err.Error())
		}

		player1Resources := game.ResourceHandByPlayer("1")
		player2Resources := game.ResourceHandByPlayer("2")
		if player1Resources["Lumber"] != 0 {
			t.Errorf("expected player#1 to have 0 Lumber, actually got %d", player1Resources["Lumber"])
		}
		if player1Resources["Brick"] != 0 {
			t.Errorf("expected player#1 to have 0 Brick, actually got %d", player1Resources["Brick"])
		}
		if player1Resources["Ore"] != 2 {
			t.Errorf("expected player#1 to have 2 Ore, actually got %d", player1Resources["Ore"])
		}
		if player2Resources["Lumber"] != 2 {
			t.Errorf("expected player#2 to have 2 Lumber, actually got %d", player2Resources["Lumber"])
		}
		if player2Resources["Brick"] != 2 {
			t.Errorf("expected player#2 to have 2 Brick, actually got %d", player2Resources["Brick"])
		}
		if player2Resources["Ore"] != 0 {
			t.Errorf("expected player#2 to have 0 Ore, actually got %d", player2Resources["Ore"])
		}

		activeTrades := game.ActiveTradeOffers()
		if len(activeTrades) != 0 {
			t.Errorf("expected trade to be finalized after exchanging resources, but actually got that it isn't")
		}
	})
}

func TestFinalizeAcceptedTradeOfferWithNoAvailableResources(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
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

	t.Run("finalize a regular trade offer - players has not the resources available", func(t *testing.T) {
		tradeID, err := game.MakeTradeOffer("1", map[string]int{
			"Lumber": 1,
			"Brick":  1,
		}, map[string]int{
			"Ore": 1,
		}, []string{})
		if err != nil {
			t.Errorf("expected to make trade offer just fine, but actually got error %s", err.Error())
		}
		secondTradeID, err := game.MakeTradeOffer("1", map[string]int{
			"Grain": 1,
		}, map[string]int{
			"Ore": 1,
		}, []string{})

		err = game.AcceptTradeOffer("2", tradeID)
		if err != nil {
			t.Errorf("expected to accept offer just fine, but actually got error %s", err.Error())
		}
		err = game.AcceptTradeOffer("2", secondTradeID)
		if err != nil {
			t.Errorf("expected to accept offer just fine, but actually got error %s", err.Error())
		}

		err = game.FinalizeTrade("1", "2", tradeID)
		if err != nil {
			t.Errorf("expected to finalize offer just fine, but actually got error %s", err.Error())
		}
		err = game.FinalizeTrade("1", "2", secondTradeID)
		if err == nil {
			t.Errorf("expected not to finalize offer, but actually got no error")
		}

		player1Resources := game.ResourceHandByPlayer("1")
		player2Resources := game.ResourceHandByPlayer("2")
		if player1Resources["Lumber"] != 0 {
			t.Errorf("expected player#1 to have 0 Lumber, actually got %d", player1Resources["Lumber"])
		}
		if player1Resources["Brick"] != 0 {
			t.Errorf("expected player#1 to have 0 Brick, actually got %d", player1Resources["Brick"])
		}
		if player1Resources["Ore"] != 2 {
			t.Errorf("expected player#1 to have 2 Ore, actually got %d", player1Resources["Ore"])
		}
		if player2Resources["Lumber"] != 2 {
			t.Errorf("expected player#2 to have 2 Lumber, actually got %d", player2Resources["Lumber"])
		}
		if player2Resources["Brick"] != 2 {
			t.Errorf("expected player#2 to have 2 Brick, actually got %d", player2Resources["Brick"])
		}
		if player2Resources["Ore"] != 0 {
			t.Errorf("expected player#2 to have 0 Ore, actually got %d", player2Resources["Ore"])
		}

		activeTrades := game.ActiveTradeOffers()
		if len(activeTrades) != 1 {
			t.Errorf("expected to have 1 active trade offer left, but actually got length %d", len(activeTrades))
		}
	})
}

func TestFinalizeAcceptedTradeOfferByRound(t *testing.T) {
	createGame := func(roundType int) *testUtils.GameState {
		game := testUtils.CreateTestGame(
			testUtils.MockWithRoundType(roundType),
			testUtils.MockWithResourcesByPlayer(map[string]map[string]int{
				"1": {
					"Lumber": 1,
					"Brick":  1,
					"Sheep":  1,
					"Grain":  1,
					"Ore":    1,
				},
				"2": {
					"Ore": 1,
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
			tradeID, _ := game.MakeTradeOffer("1",
				map[string]int{
					"Lumber": 1,
				},
				map[string]int{
					"Ore": 1,
				},
				[]string{})
			game.AcceptTradeOffer("2", tradeID)
			err := game.FinalizeTrade("1", "2", tradeID)
			hasErr := err != nil
			if hasErr != willHaveError {
				t.Errorf("expected error to be %v, but actually was %v", willHaveError, hasErr)
			}
		})
	}
}

func TestAcceptCounterOfferWithAvailableResources(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
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

	t.Run("finalize a counter trade offer - players have the resources available", func(t *testing.T) {
		tradeID, err := game.MakeTradeOffer("1", map[string]int{
			"Lumber": 1,
			"Brick":  1,
		}, map[string]int{
			"Ore": 1,
		}, []string{})
		if err != nil {
			t.Errorf("expected to make trade offer just fine, but actually got error %s", err.Error())
		}

		counterOfferID, err := game.MakeCounterTradeOffer("2", tradeID, map[string]int{
			"Ore": 1,
		}, map[string]int{
			"Lumber": 1,
			"Brick":  1,
			"Sheep":  1,
		})
		if err != nil {
			t.Errorf("expected to make counter trade offer just fine, but actually got error %s", err.Error())
		}

		err = game.FinalizeCounterTrade("1", "2", counterOfferID)
		if err != nil {
			t.Errorf("expected to accept offer just fine, but actually got error %s", err.Error())
		}

		activeTrades := game.ActiveTradeOffers()
		if len(activeTrades) != 0 {
			t.Errorf("expected to not have any active trades, but actually got length %d", len(activeTrades))
		}

		player1Resources := game.ResourceHandByPlayer("1")
		player2Resources := game.ResourceHandByPlayer("2")
		if player1Resources["Lumber"] != 0 {
			t.Errorf("expected player#1 to have 0 Lumber, actually got %d", player1Resources["Lumber"])
		}
		if player1Resources["Brick"] != 0 {
			t.Errorf("expected player#1 to have 0 Brick, actually got %d", player1Resources["Brick"])
		}
		if player1Resources["Sheep"] != 0 {
			t.Errorf("expected player#1 to have 0 Sheep, actually got %d", player1Resources["Sheep"])
		}
		if player1Resources["Ore"] != 2 {
			t.Errorf("expected player#1 to have 2 Ore, actually got %d", player1Resources["Ore"])
		}
		if player2Resources["Lumber"] != 2 {
			t.Errorf("expected player#2 to have 2 Lumber, actually got %d", player2Resources["Lumber"])
		}
		if player2Resources["Brick"] != 2 {
			t.Errorf("expected player#2 to have 2 Brick, actually got %d", player2Resources["Brick"])
		}
		if player2Resources["Sheep"] != 2 {
			t.Errorf("expected player#2 to have 2 Sheep, actually got %d", player2Resources["Sheep"])
		}
		if player2Resources["Ore"] != 0 {
			t.Errorf("expected player#2 to have 0 Ore, actually got %d", player2Resources["Ore"])
		}
	})
}

func TestAcceptTradeOfferWithActiveCounterOfferAndAvailableResources(t *testing.T) {
	game := testUtils.CreateTestGame(
		testUtils.MockWithRoundType(testUtils.Regular),
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
			"3": {
				"Lumber": 1,
				"Brick":  1,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
		}),
	)

	t.Run("finalize a trade offer with an active counter offer - players have the resources available", func(t *testing.T) {
		tradeID, err := game.MakeTradeOffer("1", map[string]int{
			"Lumber": 1,
			"Brick":  1,
		}, map[string]int{
			"Ore": 1,
		}, []string{})
		if err != nil {
			t.Errorf("expected to make trade offer just fine, but actually got error %s", err.Error())
		}

		_, err = game.MakeCounterTradeOffer("2", tradeID, map[string]int{
			"Ore": 1,
		}, map[string]int{
			"Lumber": 1,
			"Brick":  1,
			"Sheep":  1,
		})
		if err != nil {
			t.Errorf("expected to make counter trade offer just fine, but actually got error %s", err.Error())
		}

		err = game.AcceptTradeOffer("3", tradeID)
		if err != nil {
			t.Errorf("expected to accept offer just fine, but actually got error %s", err.Error())
		}

		err = game.FinalizeTrade("1", "3", tradeID)
		if err != nil {
			t.Errorf("expected to finalize offer just fine, but actually got error %s", err.Error())
		}

		activeTrades := game.ActiveTradeOffers()
		if len(activeTrades) != 0 {
			t.Errorf("expected to not have any active trades, but actually got length %d", len(activeTrades))
		}
	})
}
