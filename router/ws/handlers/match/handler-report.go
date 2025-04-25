package match

import (
	"github.com/victoroliveirab/settlers/core/packages/summary"
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

type reportResponsePayload struct {
	Report summary.ReportOutput `json:"report"`
}

func handleReportRequest(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	room := player.Room
	game := room.Game
	report := game.GetReport()

	msg := &types.WebSocketServerResponse{
		Type: "match.report.success",
		Payload: reportResponsePayload{
			Report: report,
		},
	}
	room.EnqueueOutgoingMessage(msg, []string{player.Username}, nil)
	return true, nil
}
