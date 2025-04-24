package postmatch

import (
	"time"

	"github.com/victoroliveirab/settlers/core"
)

type postMatchDataResponsePayload struct {
	Points        map[string]int  `json:"points"`
	RoomStatus    string          `json:"roomStatus"`
	RoundsPlayed  int             `json:"roundsPlayed"`
	Statistics    core.Statistics `json:"statistics"`
	StartDatetime time.Time       `json:"startDatetime"`
	EndDatetime   time.Time       `json:"endDatetime"`
}
