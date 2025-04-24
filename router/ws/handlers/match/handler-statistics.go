package match

import (
	"github.com/victoroliveirab/settlers/core/packages/summary"
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

type statisticsResponsePayload struct {
	Statistics summary.Statistics `json:"statistics"`
}

func handleStatisticsRequest(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	room := player.Room
	game := room.Game
	stats := game.GetStatistics()

	msg := &types.WebSocketServerResponse{
		Type: "match.statistics.success",
		Payload: statisticsResponsePayload{
			Statistics: stats,
		},
	}
	room.EnqueueOutgoingMessage(msg, []string{player.Username}, nil)
	return true, nil
}
