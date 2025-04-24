package core

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core/packages/round"
)

func (state *GameState) DiscardPlayerCards(playerID string, resources map[string]int) error {
	if state.round.GetRoundType() != round.DiscardPhase {
		err := fmt.Errorf("Cannot discard cards during %s", state.round.GetCurrentRoundTypeDescription())
		return err
	}

	playerState := state.playersStates[playerID]
	playerDiscardAmount := playerState.GetDiscardAmount()
	if playerDiscardAmount == 0 {
		err := fmt.Errorf("Cannot discard cards: player mustn't discard")
		return err
	}

	if playerState.GetHasDiscardedThisTurn() {
		err := fmt.Errorf("Cannot discard cards: player already discarded this round")
		return err
	}

	playerHand := playerState.GetResources()
	discardingTotal := 0

	for resource, quantity := range resources {
		if quantity > playerHand[resource] {
			err := fmt.Errorf("Cannot discard %d %s: doesn't have that amount", quantity, resource)
			return err
		}
		discardingTotal += quantity
	}

	if discardingTotal != playerDiscardAmount {
		err := fmt.Errorf("Cannot discard %d cards: must discard %d", discardingTotal, playerDiscardAmount)
		return err
	}

	for resource, quantity := range resources {
		playerState.RemoveResource(resource, quantity)
		state.bookKeeping.AddResourceDiscarded(playerID, resource, quantity)
	}
	playerState.SetHasDiscardedThisTurn(true)

	for _, player := range state.players {
		playerState := state.playersStates[player.ID]
		if playerState.GetHasDiscardedThisTurn() {
			continue
		}
		if playerState.GetDiscardAmount() > 0 {
			return nil
		}
	}

	state.round.SetRoundType(round.MoveRobberDue7)
	return nil
}
