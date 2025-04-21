package core

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core/packages/trade"
	"github.com/victoroliveirab/settlers/utils"
)

func (state *GameState) MakeBankTrade(playerID string, givenResources, requestedResources map[string]int) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot trade with bank during other player's turn")
		return err
	}

	if state.roundType != Regular {
		err := fmt.Errorf("Cannot trade with bank during %s", RoundTypeTranslation[state.roundType])
		return err
	}

	// REFACTOR: perhaps should be moved to the trade manager
	ownedPorts := state.PortsByPlayer(playerID)
	if utils.SliceContains(ownedPorts, "General") {
		err := fmt.Errorf("Cannot trade with bank: owns General port")
		return err
	}

	playerState := state.playersStates[playerID]
	return state.trade.MakeBankTrade(
		playerState,
		state.bankTradeAmount,
		givenResources,
		requestedResources,
	)
}

func (state *GameState) MakeGeneralPortTrade(playerID string, givenResources, requestedResources map[string]int) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot trade with port during other player's turn")
		return err
	}

	if state.roundType != Regular {
		err := fmt.Errorf("Cannot trade with port during %s", RoundTypeTranslation[state.roundType])
		return err
	}

	playerState := state.playersStates[playerID]
	return state.trade.MakeGeneralPortTrade(
		playerState,
		state.generalPortCost,
		givenResources,
		requestedResources,
	)
}

func (state *GameState) MakeResourcePortTrade(playerID string, givenResources, requestedResources map[string]int) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot trade with port during other player's turn")
		return err
	}

	if state.roundType != Regular {
		err := fmt.Errorf("Cannot trade with port during %s", RoundTypeTranslation[state.roundType])
		return err
	}

	playerState := state.playersStates[playerID]
	return state.trade.MakeResourcePortTrade(
		playerState,
		state.resourcePortCost,
		givenResources,
		requestedResources,
	)
}

func (state *GameState) MakeTradeOffer(playerID string, givenResources, requestedResources map[string]int, blockedPlayers []string) (int, error) {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot create trade offer during other player's turn")
		return -1, err
	}

	if state.roundType != Regular {
		err := fmt.Errorf("Cannot create trade offer during %s", RoundTypeTranslation[state.roundType])
		return -1, err
	}

	playerState := state.playersStates[playerID]
	return state.trade.MakeTradeOffer(
		playerState,
		givenResources,
		requestedResources,
		state.players,
		blockedPlayers,
	)
}

func (state *GameState) MakeCounterTradeOffer(playerID string, tradeID int, givenResources, requestedResources map[string]int) (int, error) {
	playerState := state.playersStates[playerID]
	return state.trade.MakeCounterTradeOffer(
		playerState,
		tradeID,
		givenResources,
		requestedResources,
		state.players,
	)
}

func (state *GameState) AcceptTradeOffer(playerID string, tradeID int) error {
	playerState := state.playersStates[playerID]
	return state.trade.AcceptTradeOffer(playerState, tradeID)
}

func (state *GameState) FinalizeTrade(playerID, accepterID string, tradeID int) error {
	// REFACTOR: Probably unnecessary? -> will be cought below
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot finalize a trade during other player's round")
		return err
	}
	ownerState := state.playersStates[playerID]
	accepterState := state.playersStates[accepterID]
	return state.trade.FinalizeTrade(ownerState, accepterState, tradeID)
}

func (state *GameState) RejectTradeOffer(playerID string, tradeID int) error {
	playerState := state.playersStates[playerID]
	return state.trade.RejectTradeOffer(playerState, tradeID)
}

func (state *GameState) CancelTradeOffer(playerID string, tradeID int) error {
	playerState := state.playersStates[playerID]
	return state.trade.CancelTradeOffer(playerState, tradeID)
}

func (state *GameState) GetTradeByID(tradeID int) *trade.Trade {
	return state.trade.GetTrade(tradeID)
}
