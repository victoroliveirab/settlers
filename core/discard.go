package core

import "fmt"

func (state *GameState) DiscardPlayerCards(playerID string, resources map[string]int) error {
	if state.roundType != DiscardPhase {
		err := fmt.Errorf("Cannot discard during %s", RoundTypeTranslation[state.roundType])
		return err
	}

	playerDiscardAmount := state.discardAmountByPlayer(playerID)
	if playerDiscardAmount == 0 {
		err := fmt.Errorf("Player mustn't discard")
		return err
	}

	playerState := state.playersStates[playerID]

	if playerState.HasDiscardedThisRound {
		err := fmt.Errorf("Player already discarded this round")
		return err
	}

	playerHand := playerState.Resources
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
	}
	playerState.HasDiscardedThisRound = true

	for _, player := range state.players {
		playerState := state.playersStates[player.ID]
		if playerState.HasDiscardedThisRound {
			continue
		}
		if playerState.DiscardAmount > 0 {
			return nil
		}
	}

	state.roundType = MoveRobberDue7
	return nil
}
