package prematch

import "github.com/victoroliveirab/settlers/router/ws/entities"

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
