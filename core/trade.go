package core

import (
	"fmt"
	"time"

	"github.com/victoroliveirab/settlers/utils"
)

func (state *GameState) MakeBankTrade(playerID string, givenResource, desiredResource string) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot trade with bank during other player's turn")
		return err
	}

	if state.roundType != Regular {
		err := fmt.Errorf("Cannot trade with bank during %s", RoundTypeTranslation[state.roundType])
		return err
	}

	if state.playerResourceHandMap[playerID][givenResource] < state.bankTradeAmount {
		err := fmt.Errorf("Cannot trade with bank: need %d %s, only have %d", state.bankTradeAmount, givenResource, state.playerResourceHandMap[playerID][givenResource])
		return err
	}

	state.playerResourceHandMap[playerID][givenResource] -= state.bankTradeAmount
	state.playerResourceHandMap[playerID][desiredResource]++
	return nil
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

	for resource, quantity := range givenResources {
		totalFromOfferedResource := state.playerResourceHandMap[playerID][resource]
		if totalFromOfferedResource < quantity {
			err := fmt.Errorf("Cannot make such offer: wants to give %d %s, but only have %d", quantity, resource, totalFromOfferedResource)
			return -1, err
		}
	}

	tradeID := state.playerTradeId
	state.playerTradeId++

	opponents := make(map[string]*TradePlayerEntry)

	for _, opponent := range state.players {
		if opponent.ID == playerID {
			continue
		}
		opponents[opponent.ID] = &TradePlayerEntry{
			Status:  "Open",
			Blocked: utils.SliceContains(blockedPlayers, opponent.ID),
		}
	}

	state.playersTrades[tradeID] = &Trade{
		ID:        tradeID,
		PlayerID:  playerID,
		Opponents: opponents,
		Offer:     givenResources,
		Request:   requestedResources,
		Status:    "Open",
		ParentID:  -1,
		Counters:  []int{},
		Finalized: false,
		Timestamp: time.Now().UnixMilli(),
	}

	return tradeID, nil
}

func (state *GameState) MakePortTrade(playerID string, vertexID int, givenResource, wantedResource string) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot trade with port during other player's turn")
		return err
	}

	if state.roundType != Regular {
		err := fmt.Errorf("Cannot trade with port during %s", RoundTypeTranslation[state.roundType])
		return err
	}

	var portType string
	for _, portID := range state.playerPortMap[playerID] {
		if portID == vertexID {
			portType = state.ports[portID]
		}
	}

	if portType == "" {
		err := fmt.Errorf("Vertex#%d doesn't have a port owned by player %s", vertexID, playerID)
		return err
	}

	if portType != "General" && portType != givenResource {
		err := fmt.Errorf("Vertex#%d port is of resource %s, not %s", vertexID, portType, givenResource)
		return err
	}

	var neededResources int
	if portType == "General" {
		neededResources = 3
	} else {
		if portType != givenResource {
			err := fmt.Errorf("Port type is %s, but given resource was %s", portType, givenResource)
			return err
		}
		neededResources = 2
	}

	if state.playerResourceHandMap[playerID][givenResource] < neededResources {
		err := fmt.Errorf("Player doesn't have %d %s to trade in port %s", neededResources, givenResource, portType)
		return err
	}

	state.playerResourceHandMap[playerID][givenResource] -= neededResources
	state.playerResourceHandMap[playerID][wantedResource]++
	return nil
}

// FIXME: check blocked players to not allow them to make counter offers
// TODO: don't allow counter trade offer to be equal to trade offer
func (state *GameState) MakeCounterTradeOffer(playerID string, tradeID int, counterOfferedResources, counterRequestedResources map[string]int) (int, error) {
	trade, exists := state.playersTrades[tradeID]
	if !exists {
		err := fmt.Errorf("Invalid tradeID %d", tradeID)
		return -1, err
	}

	if trade.PlayerID == playerID {
		err := fmt.Errorf("Cannot create counter offer to own offer")
		return -1, err
	}

	for resource, quantity := range counterOfferedResources {
		totalFromOfferedResource := state.playerResourceHandMap[playerID][resource]
		if totalFromOfferedResource < quantity {
			err := fmt.Errorf("Cannot make such counter offer: wants to give %d %s, but only have %d", quantity, resource, totalFromOfferedResource)
			return -1, err
		}
	}

	counterTradeID := state.playerTradeId
	state.playerTradeId++

	opponents := make(map[string]*TradePlayerEntry)

	for _, opponent := range state.players {
		if opponent.ID == trade.PlayerID {
			continue
		}
		var status string
		if opponent.ID == playerID {
			status = "Accepted"
		} else {
			status = "Open"
		}

		opponents[opponent.ID] = &TradePlayerEntry{
			Status:  status,
			Blocked: false, // TODO: have a easy way to block originally blocked players
		}
	}

	counterTrade := &Trade{
		ID:        counterTradeID,
		PlayerID:  playerID,
		Opponents: opponents,
		Offer:     counterOfferedResources,
		Request:   counterRequestedResources,
		Status:    "Open",
		ParentID:  tradeID,
		Counters:  nil,
		Finalized: false,
		Timestamp: time.Now().UnixMilli(),
	}

	state.playersTrades[tradeID].Counters = append(state.playersTrades[tradeID].Counters, counterTradeID)
	state.playersTrades[counterTradeID] = counterTrade
	return counterTradeID, nil
}

func (state *GameState) AcceptTradeOffer(playerID string, tradeID int) error {
	trade, exists := state.playersTrades[tradeID]
	if !exists {
		err := fmt.Errorf("Invalid tradeID %d", tradeID)
		return err
	}

	if playerID == state.playersTrades[tradeID].PlayerID {
		err := fmt.Errorf("Cannot accept own offer")
		return err
	}

	for resource, quantity := range trade.Request {
		if state.playerResourceHandMap[playerID][resource] < quantity {
			err := fmt.Errorf("Cannot accept offer %d: not enough %s", tradeID, resource)
			state.RejectTradeOffer(playerID, tradeID)
			return err
		}
	}

	state.playersTrades[tradeID].Opponents[playerID].Status = "Accepted"
	return nil
}

func (state *GameState) FinalizeTrade(playerID, accepterID string, tradeID int) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot finalize a trade during other player's round")
		return err
	}

	trade, exists := state.playersTrades[tradeID]
	if !exists {
		err := fmt.Errorf("Invalid tradeID %d", tradeID)
		return err
	}

	if trade.PlayerID != playerID {
		err := fmt.Errorf("Cannot finalize a trade created by other player")
		return err
	}

	if playerID == accepterID {
		err := fmt.Errorf("Cannot finalize a trade between the same players")
		return err
	}

	if trade.Opponents[accepterID].Status != "Accepted" {
		err := fmt.Errorf("Cannot finalize a trade with a player that didn't accept")
		return err
	}

	var err error
	// Check if original offerer still has the available resources - could have accepted a different offer in the mean time
	for resource, quantity := range trade.Offer {
		totalFromOfferedResource := state.playerResourceHandMap[playerID][resource]
		if totalFromOfferedResource < quantity {
			err = fmt.Errorf("Offer %d cannot be accepted at the moment: player %s wants to give %d %s, but only has %d", tradeID, trade.PlayerID, quantity, resource, totalFromOfferedResource)
			break
		}
	}

	if err != nil {
		state.cancelOffer(trade)
		return err
	}

	// Check if accepter still has the available resources - could have accepted a different offer in the mean time
	for resource, quantity := range trade.Request {
		totalFromRequestedResource := state.playerResourceHandMap[accepterID][resource]
		if totalFromRequestedResource < quantity {
			err = fmt.Errorf("Offer %d cannot be accepted at the moment by player %s: they don't have %d %s", tradeID, accepterID, quantity, resource)
			break
		}
	}

	if err != nil {
		trade.Opponents[accepterID].Status = "Declined"
		return err
	}

	for resource, quantity := range trade.Offer {
		if quantity > 0 {
			state.playerResourceHandMap[playerID][resource] -= quantity
			state.playerResourceHandMap[accepterID][resource] += quantity
		}
	}
	for resource, quantity := range trade.Request {
		if quantity > 0 {
			state.playerResourceHandMap[playerID][resource] += quantity
			state.playerResourceHandMap[accepterID][resource] -= quantity
		}
	}

	trade.Finalized = true

	activeTrades := state.ActiveTradeOffers()
	for _, activeTrade := range activeTrades {
		if activeTrade.ParentID == trade.ID {
			state.cancelOffer(state.playersTrades[activeTrade.ID])
		}
	}

	return nil
}

// REFACTOR: Use FinalizeTrade to both operations
func (state *GameState) FinalizeCounterTrade(playerID, proposerID string, counterTradeID int) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot finalize a trade during other player's round")
		return err
	}

	counterTrade, exists := state.playersTrades[counterTradeID]
	if !exists {
		err := fmt.Errorf("Invalid tradeID %d", counterTradeID)
		return err
	}

	trade, exists := state.playersTrades[counterTrade.ParentID]
	if !exists {
		// TODO: cancel counter trade
		err := fmt.Errorf("Invalid tradeID %d", counterTrade.ParentID)
		return err
	}

	if trade.PlayerID != playerID {
		err := fmt.Errorf("Cannot finalize a trade created by other player")
		return err
	}

	if playerID == proposerID {
		err := fmt.Errorf("Cannot finalize a counter trade between the same players")
		return err
	}

	var err error
	// Check if counter offerer still has the available resources - could have accepted a different offer in the mean time
	for resource, quantity := range counterTrade.Offer {
		totalFromOfferedResource := state.playerResourceHandMap[proposerID][resource]
		if totalFromOfferedResource < quantity {
			err = fmt.Errorf("Offer %d cannot be accepted at the moment: player %s wants to give %d %s, but only has %d", counterTradeID, proposerID, quantity, resource, totalFromOfferedResource)
			break
		}
	}

	if err != nil {
		state.cancelOffer(counterTrade)
		return err
	}

	// Check if player still has the available resources - could have accepted a different offer in the mean time
	for resource, quantity := range counterTrade.Request {
		totalFromRequestedResource := state.playerResourceHandMap[playerID][resource]
		if totalFromRequestedResource < quantity {
			err = fmt.Errorf("Offer %d cannot be accepted at the moment by player %s: they don't have %d %s", counterTradeID, playerID, quantity, resource)
			break
		}
	}

	if err != nil {
		state.cancelOffer(counterTrade)
		return err
	}

	for resource, quantity := range counterTrade.Offer {
		if quantity > 0 {
			state.playerResourceHandMap[proposerID][resource] -= quantity
			state.playerResourceHandMap[playerID][resource] += quantity
		}
	}
	for resource, quantity := range counterTrade.Request {
		if quantity > 0 {
			state.playerResourceHandMap[proposerID][resource] += quantity
			state.playerResourceHandMap[playerID][resource] -= quantity
		}
	}

	counterTrade.Finalized = true
	state.cancelOffer(trade)
	return nil
}

func (state *GameState) RejectTradeOffer(playerID string, tradeID int) error {
	trade, exists := state.playersTrades[tradeID]
	if !exists {
		err := fmt.Errorf("Invalid tradeID %d", tradeID)
		return err
	}

	if trade.ParentID >= 0 {
		parentTrade := state.playersTrades[trade.ParentID]
		// Creator of trade offer is rejecting counter offer
		if parentTrade.PlayerID == playerID {
			state.playersTrades[tradeID].Finalized = true
			state.playersTrades[tradeID].Status = "Closed"
			return nil
		}
	}

	trade.Opponents[playerID].Status = "Declined"
	return nil
}

func (state *GameState) cancelOffer(offer *Trade) {
	offer.Finalized = true
	offer.Status = "Closed"
}
