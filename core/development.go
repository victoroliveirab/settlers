package core

import (
	"fmt"

	"github.com/victoroliveirab/settlers/utils"
)

func (state *GameState) BuyDevelopmentCard(playerID string) error {
	err := state.IsBuyDevelopmentCardAvailable(playerID)
	if err != nil {
		return err
	}

	state.playerResourceHandMap[playerID]["Sheep"]--
	state.playerResourceHandMap[playerID]["Grain"]--
	state.playerResourceHandMap[playerID]["Ore"]--

	card := state.developmentCards[state.developmentCardHeadIndex]
	card.RoundBought = state.roundNumber
	state.developmentCardHeadIndex++

	state.playerDevelopmentHandMap[playerID][card.Name] = append(state.playerDevelopmentHandMap[playerID][card.Name], card)

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

	state.currentPlayerNumberOfPlayedDevCards++
	state.playerDevelopmentCardUsedMap[playerID]["Knight"]++
	changed := state.recountKnights()
	if changed {
		state.updatePoints()
	}
	// If game is over, no need to make the player move robber
	if state.roundType != GameOver {
		state.roundType = MoveRobberDueKnight
		return nil
	}
	return nil
}

func (state *GameState) UseMonopoly(playerID string) error {
	err := state.consumeDevelopmentCardByPlayer(playerID, "Monopoly")
	if err != nil {
		return err
	}

	state.currentPlayerNumberOfPlayedDevCards++
	state.playerDevelopmentCardUsedMap[playerID]["Monopoly"]++
	state.roundType = MonopolyPickResource
	return nil
}

func (state *GameState) PickMonopolyResource(playerID, resourceName string) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot pick monopoly resource during other player's round")
		return err
	}

	if state.roundType != MonopolyPickResource {
		err := fmt.Errorf("Cannot pick monopoly resource during %s", RoundTypeTranslation[state.roundType])
		return err
	}
	// TODO: make resource name typesafe
	for opponentID, resources := range state.playerResourceHandMap {
		if opponentID == playerID {
			continue
		}
		quantity := resources[resourceName]
		state.playerResourceHandMap[opponentID][resourceName] = 0
		state.playerResourceHandMap[playerID][resourceName] += quantity
	}

	state.roundType = Regular

	return nil
}

func (state *GameState) UseRoadBuilding(playerID string) error {
	if len(state.playerRoadMap[playerID]) >= state.maxRoads {
		err := fmt.Errorf("Player cannot build any more roads")
		return err
	}

	err := state.consumeDevelopmentCardByPlayer(playerID, "Road Building")
	if err != nil {
		return err
	}

	state.currentPlayerNumberOfPlayedDevCards++
	state.playerDevelopmentCardUsedMap[playerID]["Road Building"]++
	state.roundType = BuildRoad1Development
	return nil
}

func (state *GameState) PickRoadBuildingSpot(playerID string, edgeID int) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot pick road building spot during other player's round")
		return err
	}

	if state.roundType != BuildRoad1Development && state.roundType != BuildRoad2Development {
		err := fmt.Errorf("Cannot pick road building spot during %s", RoundTypeTranslation[state.roundType])
		return err
	}

	// REFACTOR: this is repeating road.go code
	edge, exists := state.roadMap[edgeID]
	if exists {
		owner := state.findPlayer(edge.Owner)
		err := fmt.Errorf("Player %s already has road at edge #%d", owner, edgeID)
		return err
	}

	if len(state.playerRoadMap[playerID]) >= state.maxRoads {
		err := fmt.Errorf("Player cannot build any more roads")
		return err
	}

	if !state.ownsBuildingApproaching(playerID, edgeID) {
		err := fmt.Errorf("Cannot build isolated road (edge#%d)", edgeID)
		return err
	}
	// END REFACTOR
	state.handleNewRoad(playerID, edgeID)

	if state.roundType == BuildRoad2Development {
		state.roundType = Regular
		return nil
	}

	// Player built last available road during the first build phase of development card
	if len(state.playerRoadMap[playerID]) >= state.maxCards {
		state.roundType = Regular
		return nil
	}

	state.roundType = BuildRoad2Development
	return nil
}

func (state *GameState) UseYearOfPlenty(playerID string) error {
	err := state.consumeDevelopmentCardByPlayer(playerID, "Year of Plenty")
	if err != nil {
		return err
	}

	state.currentPlayerNumberOfPlayedDevCards++
	state.playerDevelopmentCardUsedMap[playerID]["Year of Plenty"]++
	state.roundType = YearOfPlentyPickResources
	return nil
}

func (state *GameState) PickYearOfPlentyResources(playerID, resource1, resource2 string) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot pick year of plenty resources during other player's round")
		return err
	}

	if state.roundType != YearOfPlentyPickResources {
		err := fmt.Errorf("Cannot pick year of plenty resources during %s", RoundTypeTranslation[state.roundType])
		return err
	}

	state.playerResourceHandMap[playerID][resource1]++
	state.playerResourceHandMap[playerID][resource2]++
	state.roundType = Regular
	return nil
}

func (state *GameState) consumeDevelopmentCardByPlayer(playerID, devCardType string) error {
	err := state.IsDevCardPlayable(playerID, devCardType)
	if err != nil {
		return err
	}

	// REFACTOR: logic repeated from IsDevCardPlayable
	index := -1
	for i, card := range state.playerDevelopmentHandMap[playerID][devCardType] {
		if card.RoundBought < state.roundNumber {
			index = i
			break
		}
	}
	if index == -1 {
		err := fmt.Errorf("Cannot play development card bought this turn")
		return err
	}

	cards := state.playerDevelopmentHandMap[playerID][devCardType]
	utils.SliceRemove(&cards, index)
	state.playerDevelopmentHandMap[playerID][devCardType] = cards
	return nil
}

// REFACTOR: it's strange to return an error from this function, but makes things easier
func (state *GameState) IsDevCardPlayable(playerID, devCardType string) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot use knight card during other player's round")
		return err
	}

	if state.roundType != FirstRound && state.roundType != Regular && state.roundType != BetweenTurns {
		err := fmt.Errorf("Cannot use knight card during %s", RoundTypeTranslation[state.roundType])
		return err
	}

	if len(state.playerDevelopmentHandMap[playerID][devCardType]) == 0 {
		err := fmt.Errorf("Player %s doesn't have a %s card", playerID, devCardType)
		return err
	}

	if state.currentPlayerNumberOfPlayedDevCards >= state.maxDevCardsPerRound {
		err := fmt.Errorf("Can only play %d development card(s) per turn", state.maxDevCardsPerRound)
		return err
	}

	index := -1
	for i, card := range state.playerDevelopmentHandMap[playerID][devCardType] {
		if card.RoundBought < state.roundNumber {
			index = i
			break
		}
	}
	if index == -1 {
		err := fmt.Errorf("Cannot play development card bought this turn")
		return err
	}

	return nil
}

func (state *GameState) IsBuyDevelopmentCardAvailable(playerID string) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot buy development card during other player's round")
		return err
	}

	if state.roundType != Regular {
		err := fmt.Errorf("Cannot buy development card during %s", RoundTypeTranslation[state.roundType])
		return err
	}

	if state.developmentCardHeadIndex >= len(state.developmentCards) {
		err := fmt.Errorf("There is no development card left")
		return err
	}

	resources := state.playerResourceHandMap[playerID]
	if resources["Sheep"] < 1 || resources["Grain"] < 1 || resources["Ore"] < 1 {
		err := fmt.Errorf("Insufficient resources to buy a development card")
		return err
	}

	return nil
}

func (state *GameState) NumberOfKnightsUsedByPlayer(playerID string) int {
	return state.playerDevelopmentCardUsedMap[playerID]["Knight"]
}
