package core

import "fmt"

func (state *GameState) BuyDevelopmentCard(playerID string) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot buy development card during other player's round")
		return err
	}

	if state.roundType != FirstRound && state.roundType != Regular {
		err := fmt.Errorf("Cannot buy development card during %s", RoundTypeTranslation[state.roundType])
		return err
	}

	if state.developmentCardHeadIndex >= len(state.developmentCards) {
		err := fmt.Errorf("There is no development card left")
		return err
	}

	card := state.developmentCards[state.developmentCardHeadIndex]
	state.developmentCardHeadIndex++

	state.playerDevelopmentHandMap[playerID][card.Name]++

	if card.Name == "Victory Point" {
		state.updatePoints()
	}

	return nil
}

func (state *GameState) UseKnight(playerID string) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot use knight card during other player's round")
		return err
	}

	if state.roundType != FirstRound && state.roundType != Regular && state.roundType != BetweenTurns {
		err := fmt.Errorf("Cannot use knight card during %s", RoundTypeTranslation[state.roundType])
		return err
	}

	if state.playerDevelopmentHandMap[playerID]["Knight"] == 0 {
		err := fmt.Errorf("Player %s doesn't have a knight card", playerID)
		return err
	}

	if state.currentPlayerNumberOfPlayedDevCards >= state.maxDevCardsPerRound {
		err := fmt.Errorf("Can only play %d development card(s) per turn", state.maxDevCardsPerRound)
		return err
	}

	state.playerDevelopmentHandMap[playerID]["Knight"]--
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
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot use monopoly card during other player's round")
		return err
	}

	if state.roundType != FirstRound && state.roundType != Regular && state.roundType != BetweenTurns {
		err := fmt.Errorf("Cannot use monopoly card during %s", RoundTypeTranslation[state.roundType])
		return err
	}

	if state.playerDevelopmentHandMap[playerID]["Monopoly"] == 0 {
		err := fmt.Errorf("Player %s doesn't have a monopoly card", playerID)
		return err
	}

	if state.currentPlayerNumberOfPlayedDevCards >= state.maxDevCardsPerRound {
		err := fmt.Errorf("Can only play %d development card(s) per turn", state.maxDevCardsPerRound)
		return err
	}

	state.playerDevelopmentHandMap[playerID]["Monopoly"]--
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
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot use road building card during other player's round")
		return err
	}

	if state.roundType != FirstRound && state.roundType != Regular {
		err := fmt.Errorf("Cannot use road building card during %s", RoundTypeTranslation[state.roundType])
		return err
	}

	if state.playerDevelopmentHandMap[playerID]["Road Building"] == 0 {
		err := fmt.Errorf("Player %s doesn't have a road building card", playerID)
		return err
	}

	if state.currentPlayerNumberOfPlayedDevCards >= state.maxDevCardsPerRound {
		err := fmt.Errorf("Can only play %d development card(s) per turn", state.maxDevCardsPerRound)
		return err
	}

	if len(state.playerRoadMap[playerID]) >= state.maxRoads {
		err := fmt.Errorf("Player cannot build any more roads")
		return err
	}

	state.playerDevelopmentHandMap[playerID]["Road Building"]--
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
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot use monopoly card during other player's round")
		return err
	}

	if state.roundType != FirstRound && state.roundType != Regular {
		err := fmt.Errorf("Cannot use monopoly card during %s", RoundTypeTranslation[state.roundType])
		return err
	}

	if state.playerDevelopmentHandMap[playerID]["Year of Plenty"] == 0 {
		err := fmt.Errorf("Player %s doesn't have a Year of Plenty card", playerID)
		return err
	}

	if state.currentPlayerNumberOfPlayedDevCards >= state.maxDevCardsPerRound {
		err := fmt.Errorf("Can only play %d development card(s) per turn", state.maxDevCardsPerRound)
		return err
	}

	state.playerDevelopmentHandMap[playerID]["Year of Plenty"]--
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
