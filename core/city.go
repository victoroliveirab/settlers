package core

import (
	"fmt"

	"github.com/victoroliveirab/settlers/utils"
)

func (state *GameState) BuildCity(playerID string, vertexID int) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot build city during other player's turn")
		return err
	}

	if state.roundType != Regular {
		err := fmt.Errorf("Cannot build city during %s", RoundTypeTranslation[state.roundType])
		return err
	}

	resources := state.playerResourceHandMap[playerID]
	if resources["Grain"] < 2 || resources["Ore"] < 3 {
		err := fmt.Errorf("Insufficient resources to build a city")
		return err
	}

	numberOfCities := len(state.playerCityMap[playerID])
	if numberOfCities >= state.maxCities {
		err := fmt.Errorf("Cannot have more than %d cities at once", state.maxCities)
		return err
	}

	for index, settlementVertexID := range state.playerSettlementMap[playerID] {
		if settlementVertexID == vertexID {
			entry := Building{
				ID:    vertexID,
				Owner: playerID,
			}
			state.settlementMap[vertexID] = Building{}
			state.cityMap[vertexID] = entry

			settlements := state.playerSettlementMap[playerID]
			utils.SliceRemove(&settlements, index)
			state.playerSettlementMap[playerID] = settlements
			state.playerCityMap[playerID] = append(state.playerCityMap[playerID], vertexID)

			state.playerResourceHandMap[playerID]["Grain"] -= 2
			state.playerResourceHandMap[playerID]["Ore"] -= 3

			state.updatePoints()
			return nil
		}
	}

	err := fmt.Errorf("Cannot build city at vertex#%d since player doesn't have a settlement in that spot", vertexID)
	return err
}
