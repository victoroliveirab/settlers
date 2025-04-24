package core

import (
	"github.com/victoroliveirab/settlers/core/packages/board"
	"github.com/victoroliveirab/settlers/core/packages/round"
	coreT "github.com/victoroliveirab/settlers/core/types"
)

func (state *GameState) IsPlayerTurn(playerID string) bool {
	return state.currentPlayer().ID == playerID
}

func (state *GameState) IsPassTurnAllowed(playerID string) bool {
	isPlayerRound := state.IsPlayerTurn(playerID)
	roundType := state.round.GetRoundType()
	return isPlayerRound && roundType == round.Regular
}

func (state *GameState) IsPickingRobbedAllowed(playerID string) bool {
	isPlayerRound := state.IsPlayerTurn(playerID)
	roundType := state.round.GetRoundType()
	return isPlayerRound && roundType == round.PickRobbed
}

func (state *GameState) IsPickingMonopolyAllowed(playerID string) bool {
	isPlayerRound := state.IsPlayerTurn(playerID)
	roundType := state.round.GetRoundType()
	return isPlayerRound && roundType == round.MonopolyPickResource
}

func (state *GameState) IsPickingYOPAllowed(playerID string) bool {
	isPlayerRound := state.IsPlayerTurn(playerID)
	roundType := state.round.GetRoundType()
	return isPlayerRound && roundType == round.YearOfPlentyPickResources
}

func (state *GameState) IsStartTradeAllowed(playerID string) bool {
	isPlayerRound := state.IsPlayerTurn(playerID)
	roundType := state.round.GetRoundType()
	return isPlayerRound && roundType == round.Regular
}

// FIXME: get hand to see if it has the required resources
func (state *GameState) IsBuyDevCardAllowed(playerID string) bool {
	isPlayerRound := state.IsPlayerTurn(playerID)
	roundType := state.round.GetRoundType()
	return isPlayerRound && roundType == round.Regular
}

func (state *GameState) IsDevCardPlayable(playerID string, devCardType string) bool {
	playerState := state.playersStates[playerID]
	cards := playerState.GetDevelopmentCards()[devCardType]
	if len(cards) == 0 {
		return false
	}
	for _, card := range cards {
		if card.RoundBought < state.round.GetRoundNumber() {
			return true
		}
	}
	return false
}

func (state *GameState) IsRoadBuilding() bool {
	roundType := state.round.GetRoundType()
	return roundType == round.SetupRoad1 ||
		roundType == round.SetupRoad2 ||
		roundType == round.BuildRoad1Development ||
		roundType == round.BuildRoad2Development
}

func (state *GameState) IsSettlementBuilding() bool {
	roundType := state.round.GetRoundType()
	return roundType == round.SetupSettlement1 || roundType == round.SetupSettlement2
}

func (state *GameState) IsRobberTurn() bool {
	roundType := state.round.GetRoundType()
	return roundType == round.MoveRobberDue7 || roundType == round.MoveRobberDueKnight
}

func (state *GameState) GetAllCities() map[int]board.Building {
	return state.board.GetCities()
}

func (state *GameState) GetAllRoads() map[int]board.Building {
	return state.board.GetRoads()
}

func (state *GameState) GetAllSettlements() map[int]board.Building {
	return state.board.GetSettlements()
}

func (state *GameState) GetBoard() []coreT.MapBlock {
	return state.board.GetTiles()
}
