package core

import (
	"fmt"
	"time"

	"github.com/victoroliveirab/settlers/utils"
)

func equalResourceMaps(a, b map[string]int) bool {
	if len(a) != len(b) {
		return false
	}
	for key, valA := range a {
		valB, ok := b[key]
		if !ok || valB != valA {
			return false
		}
	}
	return true
}

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

	responses := make(map[string]*TradePlayerEntry)

	for _, player := range state.players {
		if player.ID == playerID {
			continue
		}
		responses[player.ID] = &TradePlayerEntry{
			Status:  NoResponse,
			Blocked: utils.SliceContains(blockedPlayers, player.ID),
		}
	}

	state.playersTrades[tradeID] = &Trade{
		ID:        tradeID,
		Requester: playerID,
		Creator:   playerID,
		Responses: responses,
		Offer:     givenResources,
		Request:   requestedResources,
		Status:    TradeOpen,
		ParentID:  -1,
		Finalized: false,
		Timestamp: time.Now().UnixMilli(),
	}

	return tradeID, nil
}

func (state *GameState) MakeCounterTradeOffer(playerID string, tradeID int, givenResources, requestedResources map[string]int) (int, error) {
	parentTrade, exists := state.playersTrades[tradeID]
	if !exists {
		err := fmt.Errorf("Invalid tradeID %d", tradeID)
		return -1, err
	}

	if parentTrade.Status != TradeOpen {
		err := fmt.Errorf("Cannot create counter offer: Trade#%d is not in Open status", tradeID)
		return -1, err
	}

	if parentTrade.Creator == playerID {
		err := fmt.Errorf("Cannot create counter offer to own offer")
		return -1, err
	}

	if parentTrade.Responses[playerID] != nil && parentTrade.Responses[playerID].Blocked {
		err := fmt.Errorf("Cannot create counter offer: blocked from Trade#%d", tradeID)
		return -1, err
	}

	for resource, quantity := range givenResources {
		totalFromOfferedResource := state.playerResourceHandMap[playerID][resource]
		if totalFromOfferedResource < quantity {
			err := fmt.Errorf("Cannot make such counter offer: wants to give %d %s, but only have %d", quantity, resource, totalFromOfferedResource)
			return -1, err
		}
	}

	if equalResourceMaps(givenResources, parentTrade.Offer) && equalResourceMaps(requestedResources, parentTrade.Request) {
		return -1, fmt.Errorf("counter offer must be different from the original offer")
	}

	counterTradeID := state.playerTradeId
	state.playerTradeId++

	responses := make(map[string]*TradePlayerEntry)

	for responsePlayer, responsePlayerParams := range parentTrade.Responses {
		responses[responsePlayer] = &TradePlayerEntry{
			Status:  NoResponse,
			Blocked: responsePlayerParams.Blocked,
		}
		if responsePlayer == playerID {
			responses[responsePlayer].Status = Accepted
		}
	}

	counterTrade := &Trade{
		ID:        counterTradeID,
		Requester: parentTrade.Requester,
		Creator:   playerID,
		Responses: responses,
		Offer:     requestedResources,
		Request:   givenResources,
		Status:    TradeOpen,
		ParentID:  tradeID,
		Finalized: false,
		Timestamp: time.Now().UnixMilli(),
	}

	state.playersTrades[counterTradeID] = counterTrade
	state.tradeParentToChild[tradeID] = append(state.tradeParentToChild[tradeID], counterTradeID)
	return counterTradeID, nil
}

func (state *GameState) AcceptTradeOffer(playerID string, tradeID int) error {
	trade, exists := state.playersTrades[tradeID]

	if !exists {
		err := fmt.Errorf("Invalid tradeID %d", tradeID)
		return err
	}

	if trade.Status != TradeOpen {
		err := fmt.Errorf("Cannot create counter offer: Trade#%d is not in Open status", tradeID)
		return err
	}

	if trade.Responses[playerID] == nil || trade.Responses[playerID].Blocked {
		err := fmt.Errorf("Cannot accept offer: not part of trade#%d opponents", tradeID)
		return err
	}

	for resource, quantity := range trade.Request {
		if state.playerResourceHandMap[playerID][resource] < quantity {
			err := fmt.Errorf("Cannot accept offer %d: not enough %s", tradeID, resource)
			// state.RejectTradeOffer(playerID, tradeID)
			return err
		}
	}

	trade.Responses[playerID].Status = Accepted
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

	if trade.Status != TradeOpen {
		err := fmt.Errorf("Cannot finalize offer: Trade#%d is not in Open status", tradeID)
		return err
	}

	if trade.Requester != playerID {
		err := fmt.Errorf("Cannot finalize a trade created by other player")
		return err
	}

	if playerID == accepterID {
		err := fmt.Errorf("Cannot finalize a trade between the same players")
		return err
	}

	if trade.Responses[accepterID] == nil || trade.Responses[accepterID].Status != Accepted {
		err := fmt.Errorf("Cannot finalize a trade with a player that didn't accept")
		return err
	}

	var err error
	// Check if original offerer still has the available resources - could have accepted a different offer in the mean time
	for resource, quantity := range trade.Offer {
		totalFromOfferedResource := state.playerResourceHandMap[playerID][resource]
		if totalFromOfferedResource < quantity {
			err = fmt.Errorf("Offer %d cannot be accepted at the moment: player %s wants to give %d %s, but only has %d", tradeID, trade.Requester, quantity, resource, totalFromOfferedResource)
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
		trade.Responses[accepterID].Status = "Declined"
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
	trade.Status = TradeFinalized

	if trade.ParentID >= 0 {
		// If finalizing a counter offer, close the parent and all siblings
		parentID := trade.ParentID
		if parentTrade, ok := state.playersTrades[parentID]; ok {
			parentTrade.Status = TradeClosed
		}
		for _, siblingID := range state.tradeParentToChild[parentID] {
			if siblingTrade, ok := state.playersTrades[siblingID]; ok {
				siblingTrade.Status = TradeClosed
			}
		}
	} else {
		// If finalizing the parent trade, close all its children
		for _, childID := range state.tradeParentToChild[tradeID] {
			if childTrade, ok := state.playersTrades[childID]; ok {
				childTrade.Status = TradeClosed
			}
		}
	}

	return nil
}

func (state *GameState) RejectTradeOffer(playerID string, tradeID int) error {
	trade, exists := state.playersTrades[tradeID]
	if !exists {
		err := fmt.Errorf("Invalid tradeID %d", tradeID)
		return err
	}

	if trade.Status != TradeOpen {
		err := fmt.Errorf("Cannot reject offer: Trade#%d is not in Open status", tradeID)
		return err
	}

	if trade.Responses[playerID] != nil && trade.Responses[playerID].Blocked {
		err := fmt.Errorf("Cannot reject offer: blocked from Trade#%d", tradeID)
		return err
	}

	if trade.Requester == playerID {
		_, ok := state.playersTrades[trade.ParentID]
		if !ok {
			err := fmt.Errorf("Cannot reject own offer")
			return err
		}
		trade.Status = TradeClosed
		return nil
	}

	trade.Responses[playerID].Status = Declined
	return nil
}

func (state *GameState) CancelTradeOffer(playerID string, tradeID int) error {
	trade, exists := state.playersTrades[tradeID]
	if !exists {
		err := fmt.Errorf("Invalid tradeID %d", tradeID)
		return err
	}

	if playerID != state.playersTrades[tradeID].Creator {
		err := fmt.Errorf("Cannot cancel offer: not owned trade")
		return err
	}

	state.cancelOffer(trade)
	return nil
}

func (state *GameState) cancelOffer(offer *Trade) {
	offer.Finalized = true
	offer.Status = TradeClosed
}
