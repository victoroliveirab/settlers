package postmatch

import (
	"time"

	"github.com/victoroliveirab/settlers/core/packages/summary"
)

type postMatchDataResponsePayload struct {
	Report        summary.ReportOutput `json:"report"`
	RoomStatus    string               `json:"roomStatus"`
	RoundsPlayed  int                  `json:"roundsPlayed"`
	StartDatetime time.Time            `json:"startDatetime"`
	EndDatetime   time.Time            `json:"endDatetime"`
}
