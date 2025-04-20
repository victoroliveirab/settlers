package connect

import (
	"fmt"

	"github.com/victoroliveirab/settlers/db/models"
	"github.com/victoroliveirab/settlers/logger"
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/handlers/match"
	prematch "github.com/victoroliveirab/settlers/router/ws/handlers/pre-match"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func isUserAtRoom(user *models.User, room *entities.Room) bool {
	for _, participant := range room.Participants {
		if participant.Player != nil && participant.Player.ID == user.ID {
			return true
		}
	}
	return false
}

func HandleConnection(conn *types.WebSocketConnection, user *models.User, room *entities.Room) (*entities.GamePlayer, error) {
	alreadyPartOfRoom := isUserAtRoom(user, room)
	playerID := user.ID

	if room.Status == "prematch" {
		if alreadyPartOfRoom {
			player, err := prematch.ReconnectPlayer(room, playerID, conn, func(player *entities.GamePlayer) {
				player.Connection.Close()
				room.RemovePlayer(playerID)
			})
			if err != nil {
				return nil, err
			}
			logger.LogMessage(playerID, "prematch.ReconnectPlayer", fmt.Sprintf("Player %d reconnected to room#%s", playerID, room.ID))
			room.EnqueueOutgoingMessage(prematch.BuildRoomMessage(room, "room.new-update"), nil, nil)
			return player, nil
		} else {
			player, err := prematch.ConnectPlayer(room, user, conn, func(player *entities.GamePlayer) {
				room.RemovePlayer(playerID)
			})
			if err != nil {
				return nil, err
			}
			logger.LogMessage(playerID, "prematch.ConnectPlayer", fmt.Sprintf("Player %d connected to room#%s", playerID, room.ID))
			room.EnqueueOutgoingMessage(prematch.BuildRoomMessage(room, "room.new-update"), nil, nil)
			return player, nil
		}
	}

	if !alreadyPartOfRoom || room.Status == "over" {
		err := fmt.Errorf("Cannot connect to room %s right now: room at %s", room.ID, room.Status)
		return nil, err
	}

	if room.Status == "setup" || room.Status == "match" {
		player, err := match.ReconnectPlayer(room, playerID, conn, func(player *entities.GamePlayer) {
			room.RemovePlayer(playerID)
		})
		if err != nil {
			return nil, err
		}
		return player, nil
	}

	err := fmt.Errorf("Cannot connect to room %s right now: unknown status %s", room.ID, room.Status)
	return nil, err
}
