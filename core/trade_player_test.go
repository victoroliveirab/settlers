package core

import (
	"fmt"
	"testing"

	"github.com/victoroliveirab/settlers/core/packages/trade"
)

func TestCreateTradeOfferWithAvailableResources(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(Regular),
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
		if trade.Requester != "1" {
			t.Errorf("expected trade to belong to player#1, but actually belongs to player#%s", trade.Requester)
		}
		if trade.Status != "Open" {
			t.Errorf("expected trade offer status to be \"Open\", but actually got %s", trade.Status)
		}
		for id, opponent := range trade.Responses {
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
	game := CreateTestGame(
		MockWithRoundType(Regular),
		MockWithResourcesByPlayer(map[string]map[string]int{
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
	createGame := func(roundType int) *GameState {
		game := CreateTestGame(
			MockWithRoundType(roundType),
			MockWithResourcesByPlayer(map[string]map[string]int{
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
		SetupSettlement1:          true,
		SetupRoad1:                true,
		SetupSettlement2:          true,
		SetupRoad2:                true,
		FirstRound:                true,
		Regular:                   false,
		MoveRobberDue7:            true,
		MoveRobberDueKnight:       true,
		PickRobbed:                true,
		BetweenTurns:              true,
		BuildRoad1Development:     true,
		BuildRoad2Development:     true,
		MonopolyPickResource:      true,
		YearOfPlentyPickResources: true,
		DiscardPhase:              true,
	}

	for roundType, willHaveError := range willHaveErrorByRoundType {
		testname := fmt.Sprintf("round type: %s, will have error: %v", RoundTypeTranslation[roundType], willHaveError)
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
	game := CreateTestGame(
		MockWithRoundType(Regular),
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

func TestCreateTradeOfferWithBlockedPlayers(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(Regular),
		MockWithResourcesByPlayer(map[string]map[string]int{
			"1": {
				"Lumber": 1,
				"Brick":  1,
				"Sheep":  1,
				"Grain":  1,
				"Ore":    1,
			},
		}),
	)

	t.Run("create trade offer - has blocked player", func(t *testing.T) {
		_, err := game.MakeTradeOffer("1", map[string]int{
			"Lumber": 1,
		}, map[string]int{
			"Brick": 1,
		}, []string{"2"})
		if err != nil {
			t.Errorf("expected to make trade offer just fine, but actually got error %s", err.Error())
		}

		trade := game.ActiveTradeOffers()[0]
		if trade.Responses["2"].Blocked == false {
			t.Errorf("expected to have player#2 blocked on trade, but actually isn't blocked")
		}
	})
}

func TestCreateCounterTradeOfferWithAvailableResources(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(Regular),
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
			"Lumber": 2,
			"Brick":  1,
		}, map[string]int{
			"Ore": 1,
		})
		if err != nil {
			t.Errorf("expected to make counter trade offer just fine, but actually got error %s", err.Error())
		}

		allTrades := game.ActiveTradeOffers()
		if len(allTrades) != 2 {
			t.Errorf("expected to have 2 active trade offers, but actually got %d", len(allTrades))
		}

		trade1 := allTrades[0]
		trade2 := allTrades[1]

		if trade2.ParentID != trade1.ID {
			t.Errorf("expected counter trade offer to have parentID %d, but actually got %d", trade1.ID, trade2.ParentID)
		}
		if trade2.Responses["2"].Status != "Accepted" {
			t.Errorf("expected counter trade offer to be already accepted by creator, but actually got status %s", trade2.Responses["2"].Status)
		}
	})
}

func TestCreateCounterTradeOfferOfCounterTradeOffer(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(Regular),
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
	)

	t.Run("create a counter trade offer - parent trade offer is a counter trade offer", func(t *testing.T) {
		tradeID, err := game.MakeTradeOffer("1", map[string]int{
			"Lumber": 1,
			"Brick":  1,
		}, map[string]int{
			"Ore": 1,
		}, []string{})
		if err != nil {
			t.Errorf("expected to make trade offer just fine, but actually got error %s", err.Error())
		}

		counterTradeID, err := game.MakeCounterTradeOffer("2", tradeID, map[string]int{
			"Lumber": 2,
			"Brick":  1,
		}, map[string]int{
			"Ore": 1,
		})
		if err != nil {
			t.Errorf("expected to make counter offer just fine, but actually got error %s", err.Error())
		}

		_, err = game.MakeCounterTradeOffer("1", counterTradeID, map[string]int{
			"Lumber": 1,
			"Brick":  1,
			"Sheep":  1,
		}, map[string]int{
			"Ore": 1,
		})
		if err != nil {
			t.Errorf("expected to make counter offer of counter offer just fine, but actually got error %s", err.Error())
		}

		numberOfActiveTrades := len(game.ActiveTradeOffers())
		if numberOfActiveTrades != 3 {
			t.Errorf("expected to have 3 active trades, but actually got %d", numberOfActiveTrades)
		}
	})
}

func TestCreateCounterTradeOfferWithNoAvailableResources(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(Regular),
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
		if trade.Responses["2"].Status != "Open" {
			t.Errorf("expected to have player#2 status still \"Open\" after failed counter offer creation, but actually has %s", trade.Responses["2"].Status)
		}
	})
}

func TestCreateCounterTradeOfferNonExistentTradeOffer(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(Regular),
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
	)

	t.Run("create a counter trade offer - trade offer doesn't exist", func(t *testing.T) {
		_, err := game.MakeCounterTradeOffer("2", 0, map[string]int{
			"Ore": 2,
		}, map[string]int{
			"Lumber": 2,
			"Brick":  1,
		})
		if err == nil {
			t.Errorf("expected not to be able to make counter trade offer, but actually no error was found")
		}
	})
}

// func TestCreateCounterTradeOfferParentTradeOfferAlreadyFinalized(t *testing.T) {
// 	game := CreateTestGame(
// 		MockWithRoundType(Regular),
// 		MockWithResourcesByPlayer(map[string]map[string]int{
// 			"1": {
// 				"Lumber": 1,
// 				"Brick":  1,
// 				"Sheep":  1,
// 				"Grain":  1,
// 				"Ore":    1,
// 			},
// 			"2": {
// 				"Lumber": 1,
// 				"Brick":  1,
// 				"Sheep":  1,
// 				"Grain":  1,
// 				"Ore":    1,
// 			},
// 		}),
// 	)
//
// 	t.Run("create a counter trade offer - trade offer already finalized", func(t *testing.T) {
//
// 	})
// }

func TestCreateCounterTradeOfferAsBlockedPlayer(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(Regular),
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
	)

	t.Run("create a counter trade offer - player is blocked from parent trade offer", func(t *testing.T) {
		tradeID, err := game.MakeTradeOffer("1", map[string]int{
			"Lumber": 1,
			"Brick":  1,
		}, map[string]int{
			"Ore": 1,
		}, []string{"2"})
		if err != nil {
			t.Errorf("expected to make trade offer just fine, but actually got error %s", err.Error())
		}

		_, err = game.MakeCounterTradeOffer("2", tradeID, map[string]int{
			"Ore": 1,
		}, map[string]int{
			"Lumber": 1,
		})
		if err == nil {
			t.Errorf("expected not to be able to make counter trade offer as blocked player, but actually no error was found")
		}
	})
}

func TestCreateCounterTradeOfferEqualToParentTrade(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(Regular),
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
	)

	t.Run("create a counter trade offer - offer is equal to original trade", func(t *testing.T) {
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
			"Lumber": 1,
			"Brick":  1,
		}, map[string]int{
			"Ore": 1,
		})
		if err == nil {
			t.Errorf("expected not to be able to make counter trade offer equal to parent trade, but actually no error was found")
		}
	})
}

func TestCreateCounterTradeOfferByRound(t *testing.T) {
	createGame := func(roundType int) *GameState {
		game := CreateTestGame(
			MockWithRoundType(roundType),
			MockWithResourcesByPlayer(map[string]map[string]int{
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
		SetupSettlement1:          true,
		SetupRoad1:                true,
		SetupSettlement2:          true,
		SetupRoad2:                true,
		FirstRound:                true,
		Regular:                   false,
		MoveRobberDue7:            true,
		MoveRobberDueKnight:       true,
		PickRobbed:                true,
		BetweenTurns:              true,
		BuildRoad1Development:     true,
		BuildRoad2Development:     true,
		MonopolyPickResource:      true,
		YearOfPlentyPickResources: true,
		DiscardPhase:              true,
	}

	for roundType, willHaveError := range willHaveErrorByRoundType {
		testname := fmt.Sprintf("round type: %s, will have error: %v", RoundTypeTranslation[roundType], willHaveError)
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
					"Lumber": 2,
				}, map[string]int{
					"Ore": 1,
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

func TestCreateCounterTradeOfferOwnTrade(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(Regular),
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
	game := CreateTestGame(
		MockWithRoundType(Regular),
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
		if trade.Responses["2"].Status != "Accepted" {
			t.Errorf("expected player#2 state to be \"Accepted\", but actually got %s", trade.Responses["2"].Status)
		}
		if trade.Responses["3"].Status != "Open" {
			t.Errorf("expected player#3 state to be \"Open\", but actually got %s", trade.Responses["3"].Status)
		}
		if trade.Responses["4"].Status != "Open" {
			t.Errorf("expected player#4 state to be \"Open\", but actually got %s", trade.Responses["4"].Status)
		}
	})
}

func TestAcceptTradeOfferIsCounterTradeOffer(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(Regular),
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
	)

	t.Run("accept a counter trade offer", func(t *testing.T) {
		tradeID, err := game.MakeTradeOffer("1", map[string]int{
			"Lumber": 1,
			"Brick":  1,
		}, map[string]int{
			"Ore": 1,
		}, []string{})
		if err != nil {
			t.Errorf("expected to make trade offer just fine, but actually got error %s", err.Error())
		}

		counterTradeID, err := game.MakeCounterTradeOffer("2", tradeID, map[string]int{
			"Lumber": 1,
			"Brick":  1,
			"Sheep":  1,
		}, map[string]int{
			"Ore": 1,
		})
		if err != nil {
			t.Errorf("expected to make counter trade offer just fine, but actually got error %s", err.Error())
		}

		err = game.AcceptTradeOffer("3", counterTradeID)
		if err != nil {
			t.Errorf("expected to accept counter offer just fine, but actually got error %s", err.Error())
		}
	})
}

func TestAcceptTradeOfferWithNoAvailableResources(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(Regular),
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
	})
}

func TestAcceptTradeOfferAsBlockedPlayer(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(Regular),
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
	)

	t.Run("accept a trade offer - player is blocked", func(t *testing.T) {
		tradeID, err := game.MakeTradeOffer("1", map[string]int{
			"Lumber": 1,
			"Brick":  1,
		}, map[string]int{
			"Ore": 1,
		}, []string{"2"})
		if err != nil {
			t.Errorf("expected to make trade offer just fine, but actually got error %s", err.Error())
		}

		err = game.AcceptTradeOffer("2", tradeID)
		if err == nil {
			t.Errorf("expected not to be able to accept trade offer as a blocked player, but no error was found")
		}
	})
}

func TestAcceptCounterTradeOfferAsBlockedPlayerOfOriginalTrade(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(Regular),
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
	)

	t.Run("accept a counter trade offer - player is blocked from parent trade", func(t *testing.T) {
		tradeID, err := game.MakeTradeOffer("1", map[string]int{
			"Lumber": 1,
			"Brick":  1,
		}, map[string]int{
			"Ore": 1,
		}, []string{"2"})
		if err != nil {
			t.Errorf("expected to make trade offer just fine, but actually got error %s", err.Error())
		}

		counterTradeID, err := game.MakeCounterTradeOffer("3", tradeID, map[string]int{
			"Lumber": 1,
			"Brick":  1,
			"Sheep":  1,
		}, map[string]int{
			"Ore": 1,
		})
		if err != nil {
			t.Errorf("expected to make counter trade offer just fine, but actually got error %s", err.Error())
		}

		err = game.AcceptTradeOffer("2", counterTradeID)
		if err == nil {
			t.Errorf("expected not to be able to accept counter trade offer as a blocked player from parent trade, but no error was found")
		}
	})
}

func TestAcceptTradeOfferByRound(t *testing.T) {
	createGame := func(roundType int) *GameState {
		game := CreateTestGame(
			MockWithRoundType(roundType),
			MockWithResourcesByPlayer(map[string]map[string]int{
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
		SetupSettlement1:          true,
		SetupRoad1:                true,
		SetupSettlement2:          true,
		SetupRoad2:                true,
		FirstRound:                true,
		Regular:                   false,
		MoveRobberDue7:            true,
		MoveRobberDueKnight:       true,
		PickRobbed:                true,
		BetweenTurns:              true,
		BuildRoad1Development:     true,
		BuildRoad2Development:     true,
		MonopolyPickResource:      true,
		YearOfPlentyPickResources: true,
		DiscardPhase:              true,
	}

	for roundType, willHaveError := range willHaveErrorByRoundType {
		testname := fmt.Sprintf("round type: %s, will have error: %v", RoundTypeTranslation[roundType], willHaveError)
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
	game := CreateTestGame(
		MockWithRoundType(Regular),
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

func TestAcceptTradeOfferOriginalCreatorTriesToAcceptCounterOffer(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(Regular),
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
	)

	t.Run("accept a counter trade offer", func(t *testing.T) {
		tradeID, err := game.MakeTradeOffer("1", map[string]int{
			"Lumber": 1,
			"Brick":  1,
		}, map[string]int{
			"Ore": 1,
		}, []string{})
		if err != nil {
			t.Errorf("expected to make trade offer just fine, but actually got error %s", err.Error())
		}

		counterTradeID, err := game.MakeCounterTradeOffer("2", tradeID, map[string]int{
			"Lumber": 1,
			"Brick":  1,
			"Sheep":  1,
		}, map[string]int{
			"Ore": 1,
		})
		if err != nil {
			t.Errorf("expected to make counter trade offer just fine, but actually got error %s", err.Error())
		}

		err = game.AcceptTradeOffer("1", counterTradeID)
		if err == nil {
			t.Errorf("expected to not be able to accept counter offer as creator of original trade offer, but no error was found")
		}
	})
}

func TestAcceptTradeOfferNonExistentTradeOffer(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(Regular),
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
	)

	t.Run("create a counter trade offer - trade offer doesn't exist", func(t *testing.T) {
		err := game.AcceptTradeOffer("2", 0)
		if err == nil {
			t.Errorf("expected not to be able to accept non existent trade offer, but actually no error was found")
		}
	})
}

// func TestAcceptTradeOfferAlreadyFinalized(t *testing.T) {}

func TestRejectTradeOffer(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(Regular),
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
	)

	t.Run("reject a trade offer", func(t *testing.T) {
		tradeID, err := game.MakeTradeOffer("1", map[string]int{
			"Lumber": 1,
			"Brick":  1,
		}, map[string]int{
			"Ore": 1,
		}, []string{})
		if err != nil {
			t.Errorf("expected to make trade offer just fine, but actually got error %s", err.Error())
		}

		err = game.RejectTradeOffer("2", tradeID)
		if err != nil {
			t.Errorf("expected to reject offer just fine, but actually got error %s", err.Error())
		}

		trade := game.ActiveTradeOffers()[0]
		if trade.Responses["2"].Status != "Declined" {
			t.Errorf("expected player#2 state to be \"Declined\", but actually got %s", trade.Responses["2"].Status)
		}
		if trade.Responses["3"].Status != "Open" {
			t.Errorf("expected player#3 state to be \"Open\", but actually got %s", trade.Responses["3"].Status)
		}
		if trade.Responses["4"].Status != "Open" {
			t.Errorf("expected player#4 state to be \"Open\", but actually got %s", trade.Responses["4"].Status)
		}
	})
}

func TestRejectTradeOfferOwnCreator(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(Regular),
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
	)

	t.Run("reject a trade offer", func(t *testing.T) {
		tradeID, err := game.MakeTradeOffer("1", map[string]int{
			"Lumber": 1,
			"Brick":  1,
		}, map[string]int{
			"Ore": 1,
		}, []string{})
		if err != nil {
			t.Errorf("expected to make trade offer just fine, but actually got error %s", err.Error())
		}

		err = game.RejectTradeOffer("1", tradeID)
		if err == nil {
			t.Errorf("expected not to be able to reject own offer, but actually got no error")
		}
	})
}

func TestRejectTradeOfferOriginalCreatorRejectsCounterOffer(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(Regular),
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
	)

	t.Run("reject a counter trade offer as creator of original", func(t *testing.T) {
		tradeID, err := game.MakeTradeOffer("1", map[string]int{
			"Lumber": 1,
			"Brick":  1,
		}, map[string]int{
			"Ore": 1,
		}, []string{})
		if err != nil {
			t.Errorf("expected to make trade offer just fine, but actually got error %s", err.Error())
		}

		counterTradeID, err := game.MakeCounterTradeOffer("2", tradeID, map[string]int{
			"Lumber": 1,
			"Brick":  1,
			"Sheep":  1,
		}, map[string]int{
			"Ore": 1,
		})
		if err != nil {
			t.Errorf("expected to make counter trade offer just fine, but actually got error %s", err.Error())
		}

		err = game.AcceptTradeOffer("3", counterTradeID)
		if err != nil {
			t.Errorf("expected to accept counter offer just fine, but actually got error %s", err.Error())
		}

		err = game.RejectTradeOffer("1", counterTradeID)
		if err != nil {
			t.Errorf("expected to reject counter offer as creator of parent just fine, but actually got error %s", err.Error())
		}

		trade1 := game.Trades()[0]
		trade2 := game.Trades()[1]

		if trade1.Status != trade.TradeOpen {
			t.Errorf("expected parent trade to still be opened, but it actually has status %s", trade1.Status)
		}
		if trade2.Status != trade.TradeClosed {
			t.Errorf("expected counter trade to be closed, but it actually has status %s", trade2.Status)
		}
	})
}

func TestRejectTradeOfferAsBlockedPlayer(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(Regular),
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
	)

	t.Run("reject a counter trade offer as creator of original", func(t *testing.T) {
		tradeID, err := game.MakeTradeOffer("1", map[string]int{
			"Lumber": 1,
			"Brick":  1,
		}, map[string]int{
			"Ore": 1,
		}, []string{"2"})
		if err != nil {
			t.Errorf("expected to make trade offer just fine, but actually got error %s", err.Error())
		}

		err = game.RejectTradeOffer("2", tradeID)
		if err == nil {
			t.Errorf("expected not to be able to reject trade offer as blocked player, but actually no error was found")
		}
	})
}

func TestFinalizeAcceptedTradeOfferWithAvailableResources(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(Regular),
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

func TestFinalizeAcceptedTradeOfferWithAccepterNoLongerHavingAvailableResources(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(Regular),
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

		currentTrade := activeTrades[0]
		if currentTrade.Responses["2"].Status != trade.Declined {
			t.Errorf("expected player 2 status to auto transition to declined, but actually got status %s", currentTrade.Responses["2"].Status)
		}
	})
}

func TestFinalizeAcceptedTradeOfferWithRequesterNoLongerHavingAvailableResources(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(Regular),
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
			"Lumber": 1,
			"Brick":  1,
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
		if len(activeTrades) != 0 {
			t.Errorf("expected to have 0 active trade offers left, but actually got length %d", len(activeTrades))
		}
	})
}

func TestFinalizeAcceptedTradeOfferByRound(t *testing.T) {
	createGame := func(roundType int) *GameState {
		game := CreateTestGame(
			MockWithRoundType(roundType),
			MockWithResourcesByPlayer(map[string]map[string]int{
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
		SetupSettlement1:          true,
		SetupRoad1:                true,
		SetupSettlement2:          true,
		SetupRoad2:                true,
		FirstRound:                true,
		Regular:                   false,
		MoveRobberDue7:            true,
		MoveRobberDueKnight:       true,
		PickRobbed:                true,
		BetweenTurns:              true,
		BuildRoad1Development:     true,
		BuildRoad2Development:     true,
		MonopolyPickResource:      true,
		YearOfPlentyPickResources: true,
		DiscardPhase:              true,
	}

	for roundType, willHaveError := range willHaveErrorByRoundType {
		testname := fmt.Sprintf("round type: %s, will have error: %v", RoundTypeTranslation[roundType], willHaveError)
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

func TestFinalizeCounterOfferWithAvailableResources(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(Regular),
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
			"Lumber": 1,
			"Brick":  1,
			"Sheep":  1,
		}, map[string]int{
			"Ore": 1,
		})
		if err != nil {
			t.Errorf("expected to make counter trade offer just fine, but actually got error %s", err.Error())
		}

		err = game.FinalizeTrade("1", "2", counterOfferID)
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

func TestFinalizeTradeOfferWithActiveCounterOfferAndAvailableResources(t *testing.T) {
	game := CreateTestGame(
		MockWithRoundType(Regular),
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
