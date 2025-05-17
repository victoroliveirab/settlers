package prematch

import (
	"fmt"

	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func SendCurrentGameState(player *entities.GamePlayer) error {
	room := player.Room
	requestType := types.RequestType("room.connect")
	msg := BuildRoomMessage(room, fmt.Sprintf("%s.success", requestType))
	wsErr := player.WriteJSON(msg)
	return wsErr
}
