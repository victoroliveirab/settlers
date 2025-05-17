package postmatch

import (
	"github.com/victoroliveirab/settlers/router/ws/entities"
)

func SendCurrentGameState(player *entities.GamePlayer) error {
	room := player.Room
	msg := BuildPostMatchHydrateMessage(room)
	wsErr := player.WriteJSON(msg)
	return wsErr
}
