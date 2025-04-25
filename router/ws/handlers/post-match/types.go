package postmatch

import (
	"time"

	"github.com/victoroliveirab/settlers/core/packages/summary"
	coreT "github.com/victoroliveirab/settlers/core/types"
)

type postMatchDataResponsePayload struct {
	Report        summary.ReportOutput `json:"report"`
	RoomStatus    string               `json:"roomStatus"`
	RoundsPlayed  int                  `json:"roundsPlayed"`
	StartDatetime time.Time            `json:"startDatetime"`
	EndDatetime   time.Time            `json:"endDatetime"`
}

type postMatchHydrateResponsePayload struct {
	Report        summary.ReportOutput `json:"report"`
	RoomName      string               `json:"roomName"`
	RoomStatus    string               `json:"roomStatus"`
	RoundsPlayed  int                  `json:"roundsPlayed"`
	StartDatetime time.Time            `json:"startDatetime"`
	EndDatetime   time.Time            `json:"endDatetime"`
	Players       []coreT.Player       `json:"players"`
	MapName       string               `json:"mapName"`
}
