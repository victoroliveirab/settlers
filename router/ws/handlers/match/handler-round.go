package match

import (
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func handleEndRound(player *entities.GamePlayer, message *types.WebSocketMessage) (bool, error) {
	room := player.Room
	game := room.Game
	err := game.EndRound(player.Username)
	if err != nil {
		wsErr := sendEndRoundError(player.Connection, player.ID, err)
		if wsErr != nil {
			return true, wsErr
		}
		return true, nil
	}

	nextPlayer := room.Participants[game.CurrentRoundPlayerIndex()].Player
	err = SendPlayerRound(room, nextPlayer)
	room.EnqueueBroadcastMessage(BuildPlayerRoundOpponentsBroadcast(room), []int64{nextPlayer.ID}, nil)
	return true, err
}
