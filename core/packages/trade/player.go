package trade

import (
	"fmt"
	"time"

	"github.com/victoroliveirab/settlers/core/packages/player"
	coreT "github.com/victoroliveirab/settlers/core/types"
	"github.com/victoroliveirab/settlers/utils"
)

func (tm *Instance) MakeTradeOffer(
	playerState *player.Instance,
	givenResources map[string]int,
	requestedResources map[string]int,
	players []coreT.Player,
	blockedPlayers []string,
) (int, error) {
	playerID := playerState.ID
	for resource, quantity := range givenResources {
		totalFromOfferedResource := playerState.Resources[resource]
		if totalFromOfferedResource < quantity {
			err := fmt.Errorf("Cannot make such offer: wants to give %d %s, but only have %d", quantity, resource, totalFromOfferedResource)
			return -1, err
		}
	}

	tradeID := tm.nextTradeID
	tm.nextTradeID++
	responses := make(map[string]*TradePlayerEntry)
	for _, player := range players {
		if player.ID == playerID {
			continue
		}
		responses[player.ID] = &TradePlayerEntry{
			Status:  NoResponse,
			Blocked: utils.SliceContains(blockedPlayers, player.ID),
		}
	}

	tm.trades[tradeID] = &Trade{
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

func (tm *Instance) MakeCounterTradeOffer(
	playerState *player.Instance,
	tradeID int,
	givenResources map[string]int,
	requestedResources map[string]int,
	players []coreT.Player,
) (int, error) {
	playerID := playerState.ID
	parentTrade, exists := tm.trades[tradeID]
	if !exists {
		err := fmt.Errorf("Cannot create counter offer: invalid tradeID %d", tradeID)
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

	var offeredResources map[string]int

	if playerID == parentTrade.Requester {
		offeredResources = givenResources
	} else {
		offeredResources = requestedResources
	}

	for resource, quantity := range offeredResources {
		totalFromOfferedResource := playerState.Resources[resource]
		if totalFromOfferedResource < quantity {
			err := fmt.Errorf("Cannot craate counter offer: wants to give %d %s, but only have %d", quantity, resource, totalFromOfferedResource)
			return -1, err
		}
	}

	if equalResourceMaps(givenResources, parentTrade.Offer) && equalResourceMaps(requestedResources, parentTrade.Request) {
		return -1, fmt.Errorf("Cannot create counter offer: must be different from the original offer")
	}

	counterTradeID := tm.nextTradeID
	tm.nextTradeID++
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

	tm.trades[counterTradeID] = &Trade{
		ID:        counterTradeID,
		Requester: parentTrade.Requester,
		Creator:   playerID,
		Responses: responses,
		Offer:     givenResources,
		Request:   requestedResources,
		Status:    TradeOpen,
		ParentID:  tradeID,
		Finalized: false,
		Timestamp: time.Now().UnixMilli(),
	}
	tm.parentToChildMap[tradeID] = append(tm.parentToChildMap[tradeID], counterTradeID)
	return counterTradeID, nil
}

func (tm *Instance) AcceptTradeOffer(playerState *player.Instance, tradeID int) error {
	playerID := playerState.ID
	trade, exists := tm.trades[tradeID]
	if !exists {
		err := fmt.Errorf("Cannot accept trade offer: invalid tradeID %d", tradeID)
		return err
	}

	if trade.Status != TradeOpen {
		err := fmt.Errorf("Cannot accept trade offer: trade#%d is not in Open status", tradeID)
		return err
	}

	if trade.Responses[playerID] == nil || trade.Responses[playerID].Blocked {
		err := fmt.Errorf("Cannot accept offer: not part of trade#%d opponents", tradeID)
		return err
	}

	for resource, quantity := range trade.Request {
		if playerState.Resources[resource] < quantity {
			err := fmt.Errorf("Cannot accept offer %d: not enough %s", tradeID, resource)
			return err
		}
	}

	trade.Responses[playerID].Status = Accepted
	return nil
}

func (tm *Instance) FinalizeTrade(
	ownerState *player.Instance,
	accepterState *player.Instance,
	tradeID int,
) error {
	playerID := ownerState.ID
	accepterID := accepterState.ID
	trade, exists := tm.trades[tradeID]
	if !exists {
		err := fmt.Errorf("Cannot finalize trade offer: invalid tradeID %d", tradeID)
		return err
	}

	if trade.Status != TradeOpen {
		err := fmt.Errorf("Cannot finalize trade offer: trade#%d is not in Open status", tradeID)
		return err
	}

	if trade.Requester != playerID {
		err := fmt.Errorf("Cannot finalize trade offer: not owned trade")
		return err
	}

	if playerID == accepterID {
		err := fmt.Errorf("Cannot finalize trade offer: cannot be creator and accepter")
		return err
	}

	if trade.Responses[accepterID] == nil || trade.Responses[accepterID].Status != Accepted {
		err := fmt.Errorf("Cannot finalize trade offer: player %s didn't accept trade offer", accepterID)
		return err
	}

	var err error
	// Check if original offerer still has the available resources - could have accepted a different offer in the mean time
	for resource, quantity := range trade.Offer {
		totalFromOfferedResource := ownerState.Resources[resource]
		if totalFromOfferedResource < quantity {
			err = fmt.Errorf("Offer %d cannot be accepted at the moment: player %s wants to give %d %s, but only has %d", tradeID, trade.Requester, quantity, resource, totalFromOfferedResource)
			break
		}
	}
	if err != nil {
		tm.CancelTrade(tradeID)
		return err
	}

	// Check if accepter still has the available resources - could have accepted a different offer in the mean time
	for resource, quantity := range trade.Request {
		totalFromRequestedResource := accepterState.Resources[resource]
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
			ownerState.RemoveResource(resource, quantity)
			accepterState.AddResource(resource, quantity)
		}
	}

	for resource, quantity := range trade.Request {
		if quantity > 0 {
			ownerState.AddResource(resource, quantity)
			accepterState.RemoveResource(resource, quantity)
		}
	}

	trade.Finalized = true
	trade.Status = TradeFinalized
	if trade.ParentID >= 0 {
		// If finalizing a counter offer, close the parent and all siblings
		parentID := trade.ParentID
		if parentTrade, ok := tm.trades[parentID]; ok {
			parentTrade.Status = TradeClosed
		}
		for _, siblingID := range tm.parentToChildMap[parentID] {
			if siblingTrade, ok := tm.trades[siblingID]; ok {
				siblingTrade.Status = TradeClosed
			}
		}
	} else {
		// If finalizing the parent trade, close all its children
		for _, childID := range tm.parentToChildMap[tradeID] {
			if childTrade, ok := tm.trades[childID]; ok {
				childTrade.Status = TradeClosed
			}
		}
	}
	return nil
}

func (tm *Instance) RejectTradeOffer(playerState *player.Instance, tradeID int) error {
	playerID := playerState.ID
	trade, exists := tm.trades[tradeID]
	if !exists {
		err := fmt.Errorf("Cannot reject trade offer: invalid tradeID %d", tradeID)
		return err
	}

	if trade.Status != TradeOpen {
		err := fmt.Errorf("Cannot reject trade offer: trade#%d is not in Open status", tradeID)
		return err
	}

	if trade.Responses[playerID] == nil || trade.Responses[playerID].Blocked {
		err := fmt.Errorf("Cannot reject offer: not part of trade#%d opponents", tradeID)
		return err
	}

	if trade.Requester == playerID {
		_, ok := tm.trades[trade.ParentID]
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

func (tm *Instance) CancelTradeOffer(playerState *player.Instance, tradeID int) error {
	playerID := playerState.ID
	_, exists := tm.trades[tradeID]
	if !exists {
		err := fmt.Errorf("Cannot cancel trade offer: invalid tradeID %d", tradeID)
		return err
	}

	if playerID != tm.trades[tradeID].Creator {
		err := fmt.Errorf("Cannot cancel offer: not owned trade")
		return err
	}

	return tm.CancelTrade(tradeID)
}
