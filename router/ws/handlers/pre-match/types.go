package prematch

import (
	coreT "github.com/victoroliveirab/settlers/core/types"
	"github.com/victoroliveirab/settlers/router/ws/entities"
)

type roomConnectResponsePayload struct {
	Room *entities.Room `json:"room"`
}

type roomUpdateParamRequestPayload struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

type roomPlayerChangeColorRequestPayload struct {
	Color string `json:"color"`
}

type roomPlayerReadyRequestPayload struct {
	Ready bool `json:"ready"`
}

type roomUpdateResponsePayload struct {
	Room       *entities.Room                 `json:"room"`
	RoomParams []entities.RoomParamsMetaEntry `json:"params"`
}

type roomStartMatchPayload struct {
	Map           []*coreT.MapBlock `json:"map"`
	MapName       string            `json:"mapName"`
	Players       []coreT.Player    `json:"players"`
	ResourceCount map[string]int    `json:"resourceCount"`
	Logs          []string          `json:"logs"`
}
