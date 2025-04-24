package trade

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core/packages/player"
	"github.com/victoroliveirab/settlers/utils"
)

func (tm *Instance) MakeGeneralPortTrade(
	playerState *player.Instance,
	generalPortTradeCost int,
	givenResources map[string]int,
	requestedResources map[string]int,
) error {
	ownedPorts := playerState.GetPortTypes()
	if !utils.SliceContains(ownedPorts, "General") {
		err := fmt.Errorf("Cannot trade in port General: doesn't own port")
		return err
	}

	ownedResources := playerState.GetResources()
	availableResourcesToRequest := 0
	for resource, quantity := range givenResources {
		if quantity == 0 {
			continue
		}
		if utils.SliceContains(ownedPorts, resource) {
			err := fmt.Errorf("Cannot trade %s in General port: owns specific port", resource)
			return err
		}
		if quantity%generalPortTradeCost != 0 {
			err := fmt.Errorf("Cannot trade %d of %s: not a multiple of %d", quantity, resource, generalPortTradeCost)
			return err
		}
		if ownedResources[resource] < quantity {
			err := fmt.Errorf("Cannot trade %d of %s with port: doesn't have that quantity available", quantity, resource)
			return err
		}
		availableResourcesToRequest += quantity / generalPortTradeCost
	}

	for resource, quantity := range requestedResources {
		givenQuantity, ok := givenResources[resource]
		if ok && givenQuantity > 0 && quantity > 0 {
			err := fmt.Errorf("Cannot complete port trade: giving and requesting %s", resource)
			return err
		}
		availableResourcesToRequest -= quantity
	}
	if availableResourcesToRequest != 0 {
		err := fmt.Errorf("Cannot complete port trade: wrong proportion of given and requested resorces")
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

func (tm *Instance) MakeResourcePortTrade(
	playerState *player.Instance,
	resourcePortTradeCost int,
	givenResources map[string]int,
	requestedResources map[string]int,
) error {
	ownedPorts := playerState.GetPortTypes()
	ownedResources := playerState.GetResources()
	availableResourcesToRequest := 0
	for resource, quantity := range givenResources {
		if quantity == 0 {
			continue
		}
		if !utils.SliceContains(ownedPorts, resource) {
			err := fmt.Errorf("Cannot trade in port %s: doesn't own port", resource)
			return err
		}
		if quantity%resourcePortTradeCost != 0 {
			err := fmt.Errorf("Cannot trade %d of %s: not a multiple of %d", quantity, resource, resourcePortTradeCost)
			return err
		}
		if ownedResources[resource] < quantity {
			err := fmt.Errorf("Cannot trade %d of %s with port: doesn't have that quantity available", quantity, resource)
			return err
		}
		availableResourcesToRequest += quantity / resourcePortTradeCost
	}

	for resource, quantity := range requestedResources {
		givenQuantity, ok := givenResources[resource]
		if ok && givenQuantity > 0 && quantity > 0 {
			err := fmt.Errorf("Cannot complete port trade: giving and requesting %s", resource)
			return err
		}
		availableResourcesToRequest -= quantity
	}
	if availableResourcesToRequest != 0 {
		err := fmt.Errorf("Cannot complete port trade: wrong proportion of given and requested resorces")
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
