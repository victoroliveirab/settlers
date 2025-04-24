package trade

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core/packages/player"
)

// TODO: register bank trade to the trade manager
func (tm *Instance) MakeBankTrade(
	playerState *player.Instance,
	bankTradeCost int,
	givenResources map[string]int,
	requestedResources map[string]int,
) error {
	availableResourcesToRequest := 0
	playerResources := playerState.GetResources()
	for resource, quantity := range givenResources {
		if quantity == 0 {
			continue
		}
		if quantity%bankTradeCost != 0 {
			err := fmt.Errorf("Cannot trade %d of %s: not a multiple of %d", quantity, resource, bankTradeCost)
			return err
		}
		if playerResources[resource] < quantity {
			err := fmt.Errorf("Cannot trade %d of %s with bank: doesn't have that quantity available", quantity, resource)
			return err
		}
		availableResourcesToRequest += quantity / bankTradeCost
	}

	for resource, quantity := range requestedResources {
		givenQuantity, ok := givenResources[resource]
		if ok && givenQuantity > 0 && quantity > 0 {
			err := fmt.Errorf("Cannot complete bank trade: giving and requesting %s", resource)
			return err
		}
		availableResourcesToRequest -= quantity
	}
	if availableResourcesToRequest != 0 {
		err := fmt.Errorf("Cannot complete bank trade: wrong proportion of given and requested resorces")
		return err
	}
	for resource, quantity := range givenResources {
		playerState.RemoveResource(resource, quantity)
	}
	for resource, quantity := range requestedResources {
		playerState.AddResource(resource, quantity)
	}
	return nil
}
