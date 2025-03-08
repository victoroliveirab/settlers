package connect

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core"
	"github.com/victoroliveirab/settlers/db/models"
	"github.com/victoroliveirab/settlers/router/ws/entities"
	matchsetup "github.com/victoroliveirab/settlers/router/ws/handlers/match-setup"
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
	fmt.Println("HandleConnection", playerID, room.Status)

	if room.Status == "prematch" {
		if alreadyPartOfRoom {
			player, err := prematch.ReconnectPlayer(room, playerID, conn, func(player *entities.GamePlayer) {
				room.RemovePlayer(playerID)
			})
			if err != nil {
				return nil, err
			}
			room.EnqueueBroadcastMessage(buildPlayerReconnectedPreMatchBroadcast(player), []int64{player.ID}, nil)
			return player, nil
		} else {
			player, err := prematch.ConnectPlayer(room, user, conn, func(player *entities.GamePlayer) {
				room.RemovePlayer(playerID)
			})
			if err != nil {
				return nil, err
			}
			room.EnqueueBroadcastMessage(buildPlayerConnectedPreMatchBroadcast(player), []int64{player.ID}, nil)
			return player, nil
		}
	}

	if !alreadyPartOfRoom {
		err := fmt.Errorf("Cannot connect to room %s right now: room at %s", room.ID, room.Status)
		return nil, err
	}

	if room.Status == "setup" {
		player, err := matchsetup.ReconnectPlayer(room, playerID, conn, func(player *entities.GamePlayer) {
			room.RemovePlayer(playerID)
		})
		if err != nil {
			return nil, err
		}
		return player, nil
	}

	if room.Status == "match" {

	}

	return nil, nil
}

func TryReconnectPlayer(player *entities.GamePlayer) error {
	room := player.Room

	game := room.Game
	// Not player's turn, nothing to do
	if game.CurrentRoundPlayer().ID != player.Username {
		return nil
	}

	if game.RoundType() == core.SetupSettlement1 || game.RoundType() == core.SetupSettlement2 {
		err := matchsetup.SendBuildSetupSettlementRequest(player)
		return err
	} else if game.RoundType() == core.SetupRoad1 || game.RoundType() == core.SetupRoad2 {
		err := matchsetup.SendBuildSetupRoadRequest(player)
		return err
	} else if game.RoundType() == core.FirstRound || game.RoundType() == core.Regular || game.RoundType() == core.BetweenTurns {
		return nil // hydrate will take care of it
	}

	err := fmt.Errorf("Cannot reconnect player#%s: not known round type %s", player.Username, core.RoundTypeTranslation[game.RoundType()])
	return err
}
