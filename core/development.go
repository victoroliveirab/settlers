package core

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core/packages/round"
	coreT "github.com/victoroliveirab/settlers/core/types"
)

func (state *GameState) BuyDevelopmentCard(playerID string) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot buy development card during other player's round")
		return err
	}

	if state.round.GetRoundType() != round.Regular {
		err := fmt.Errorf("Cannot buy development card during %s", state.round.GetCurrentRoundTypeDescription())
		return err
	}

	playerState := state.playersStates[playerID]
	if playerState.Resources["Sheep"] < 1 || playerState.Resources["Grain"] < 1 || playerState.Resources["Ore"] < 1 {
		err := fmt.Errorf("Cannot buy development card: insufficient resources")
		return err
	}

	card, err := state.development.Draw()
	if err != nil {
		return err
	}
	card.RoundBought = state.round.GetRoundNumber()

	playerState.RemoveResource("Sheep", 1)
	playerState.RemoveResource("Grain", 1)
	playerState.RemoveResource("Ore", 1)
	state.stats.AddResourcesUsed(playerID, "Sheep", 1)
	state.stats.AddResourcesUsed(playerID, "Grain", 1)
	state.stats.AddResourcesUsed(playerID, "Ore", 1)
	playerState.AddDevelopmentCard(card)
	state.stats.AddDevCardDrawn(playerID, card.Name)

	if card.Name == "Victory Point" {
		state.updatePoints()
	}

	return nil
}

func (state *GameState) UseDevelopmentCard(playerID, devCardType string) error {
	switch devCardType {
	case "Knight":
		return state.UseKnight(playerID)
	case "Monopoly":
		return state.UseMonopoly(playerID)
	case "Road Building":
		return state.UseRoadBuilding(playerID)
	case "Year of Plenty":
		return state.UseYearOfPlenty(playerID)
	case "Victory Point":
		err := fmt.Errorf("Cannot use Victory Point card")
		return err
	default:
		err := fmt.Errorf("Unknown dev card: %s", devCardType)
		return err
	}
}

func (state *GameState) UseKnight(playerID string) error {
	err := state.consumeDevelopmentCardByPlayer(playerID, "Knight")
	if err != nil {
		return err
	}

	changed := state.recountKnights()
	if changed {
		state.updatePoints()
	}
	// If game is over, no need to make the player move robber
	if state.round.GetRoundType() != round.GameOver {
		state.round.SetRoundType(round.MoveRobberDueKnight)
	}
	return nil
}

func (state *GameState) UseMonopoly(playerID string) error {
	err := state.consumeDevelopmentCardByPlayer(playerID, "Monopoly")
	if err != nil {
		return err
	}

	state.round.SetRoundType(round.MonopolyPickResource)
	return nil
}

func (state *GameState) PickMonopolyResource(playerID, resourceName string) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot pick monopoly resource during other player's round")
		return err
	}

	if state.round.GetRoundType() != round.MonopolyPickResource {
		err := fmt.Errorf("Cannot pick monopoly resource during %s", state.round.GetCurrentRoundTypeDescription())
		return err
	}

	monopolyPlayerState := state.playersStates[playerID]

	// TODO: make resource name typesafe
	for _, player := range state.players {
		if player.ID == playerID {
			continue
		}
		playerState := state.playersStates[player.ID]
		quantity := playerState.Resources[resourceName]
		if quantity > 0 {
			playerState.RemoveResource(resourceName, quantity)
			monopolyPlayerState.AddResource(resourceName, quantity)
		}
	}

	state.round.SetRoundType(round.Regular)

	return nil
}

func (state *GameState) UseRoadBuilding(playerID string) error {
	playerState := state.playersStates[playerID]
	if len(playerState.Roads) >= state.maxRoads {
		err := fmt.Errorf("Player cannot build any more roads")
		return err
	}

	err := state.consumeDevelopmentCardByPlayer(playerID, "Road Building")
	if err != nil {
		return err
	}

	state.round.SetRoundType(round.BuildRoad1Development)
	return nil
}

func (state *GameState) PickRoadBuildingSpot(playerID string, edgeID int) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot pick road building spot during other player's round")
		return err
	}

	roundType := state.round.GetRoundType()
	if roundType != round.BuildRoad1Development && roundType != round.BuildRoad2Development {
		err := fmt.Errorf("Cannot pick road building spot during %s", state.round.GetCurrentRoundTypeDescription())
		return err
	}

	// REFACTOR: this is repeating road.go code
	edge, exists := state.board.GetRoads()[edgeID]
	if exists {
		owner := state.findPlayer(edge.Owner)
		err := fmt.Errorf("Player %s already has road at edge #%d", owner, edgeID)
		return err
	}

	playerState := state.playersStates[playerID]
	if len(playerState.Roads) >= state.maxRoads {
		err := fmt.Errorf("Player cannot build any more roads")
		return err
	}

	if !state.ownsBuildingApproaching(playerID, edgeID) {
		err := fmt.Errorf("Cannot build isolated road (edge#%d)", edgeID)
		return err
	}
	// END REFACTOR
	state.handleNewRoad(playerID, edgeID)

	if state.round.GetRoundType() == round.BuildRoad2Development {
		state.round.SetRoundType(round.Regular)
		return nil
	}

	// Player built last available road during the first build phase of development card
	if len(playerState.Roads) >= state.maxCards {
		state.round.SetRoundType(round.Regular)
		return nil
	}

	state.round.SetRoundType(round.BuildRoad2Development)
	return nil
}

func (state *GameState) UseYearOfPlenty(playerID string) error {
	err := state.consumeDevelopmentCardByPlayer(playerID, "Year of Plenty")
	if err != nil {
		return err
	}

	state.round.SetRoundType(round.YearOfPlentyPickResources)
	return nil
}

func (state *GameState) PickYearOfPlentyResources(playerID, resource1, resource2 string) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot pick year of plenty resources during other player's round")
		return err
	}

	if state.round.GetRoundType() != round.YearOfPlentyPickResources {
		err := fmt.Errorf("Cannot pick year of plenty resources during %s", state.round.GetCurrentRoundTypeDescription())
		return err
	}

	playerState := state.playersStates[playerID]
	playerState.AddResource(resource1, 1)
	playerState.AddResource(resource2, 1)
	state.round.SetRoundType(round.Regular)
	return nil
}

func (state *GameState) consumeDevelopmentCardByPlayer(playerID, devCardType string) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot use knight card during other player's round")
		return err
	}

	roundType := state.round.GetRoundType()
	switch devCardType {
	case "Knight":
		if roundType != round.FirstRound && roundType != round.Regular && roundType != round.BetweenTurns {
			err := fmt.Errorf("Cannot use knight card during %s", state.round.GetCurrentRoundTypeDescription())
			return err
		}
	default:
		if roundType != round.Regular {
			err := fmt.Errorf("Cannot use %s card during %s", devCardType, state.round.GetCurrentRoundTypeDescription())
			return err
		}
	}

	playerState := state.playersStates[playerID]
	cards, exists := playerState.DevelopmentCards[devCardType]
	if !exists {
		err := fmt.Errorf("Cannot use %s card: not owned", devCardType)
		return err
	}

	if playerState.NumDevCardsPlayedTurn >= state.maxDevCardsPerRound {
		err := fmt.Errorf("Cannot use %s card: can only play %d development card(s) per turn", devCardType, state.maxDevCardsPerRound)
		return err
	}

	var cardToUse *coreT.DevelopmentCard
	for _, card := range cards {
		if card.RoundBought < state.round.GetRoundNumber() {
			cardToUse = card
			break
		}
	}
	if cardToUse == nil {
		err := fmt.Errorf("Cannot use %s card: bought this turn", devCardType)
		return err
	}
	playerState.ConsumeDevelopmentCard(cardToUse)

	return nil
}

func (state *GameState) NumberOfKnightsUsedByPlayer(playerID string) int {
	playerState := state.playersStates[playerID]
	return playerState.UsedDevelopmentCards["Knight"]
}
