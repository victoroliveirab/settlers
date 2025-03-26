package match

import (
	"fmt"

	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

type moveRobberPayload struct {
	tileID int
}

func parseMoveRobberPayload(payload map[string]interface{}) (*moveRobberPayload, error) {
	tileID, ok := payload["tile"].(float64)
	if !ok {
		err := fmt.Errorf("malformed data: tile")
		return nil, err
	}

	return &moveRobberPayload{
		tileID: int(tileID),
	}, nil
}

func handleMoveRobber(player *entities.GamePlayer, message *types.WebSocketMessage) (bool, error) {
	payload, err := parseMoveRobberPayload(message.Payload)
	if err != nil {
		// wsErr := sendMoveRobberError(player.Connection, player.ID, err)
		// return true, wsErr
		return true, err
	}

	tileID := payload.tileID
	room := player.Room
	game := room.Game
	err = game.MoveRobber(player.Username, tileID)
	if err != nil {
		// wsErr := sendMoveRobberError(player.Connection, player.ID, err)
		// return true, wsErr
		return true, err
	}

	logs := []string{fmt.Sprintf("%s moved the robber", player.Username)}

	robbablePlayers, _ := game.RobbablePlayers(player.Username)
	if len(robbablePlayers) == 0 {
		logs = append(logs, fmt.Sprintf("No player to rob"))
		// room.EnqueueBroadcastMessage(buildRobberMovedBroadcast(tileID, logs), []int64{}, nil)
		return true, nil
	} else if len(robbablePlayers) == 1 {
		err := game.RobPlayer(player.Username, robbablePlayers[0])
		if err != nil { // The only player in tile (besides the current player) doesn't have any resources
			// NOTE: This should never happen since the core logic
			logs = append(logs, fmt.Sprintf("No player to rob"))
		} else {
			logs = append(logs, fmt.Sprintf("%s robbed %s", player.Username, robbablePlayers[0]))
		}
		// room.EnqueueBroadcastMessage()
	}
	return true, nil
}
