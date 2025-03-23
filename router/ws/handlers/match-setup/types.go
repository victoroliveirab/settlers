package matchsetup

import (
	coreT "github.com/victoroliveirab/settlers/core/types"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

type newSettlementRequestPayload struct {
	VertexID int `json:"vertex"`
}

type newRoadRequestPayload struct {
	EdgeID int `json:"edge"`
}

type hydrateResponsePayload struct {
	EdgeUpdate        *types.WebSocketServerResponse `json:"edgeUpdate"`
	Map               []*coreT.MapBlock              `json:"map"`
	MapUpdate         *types.WebSocketServerResponse `json:"mapUpdate"`
	Players           []coreT.Player                 `json:"players"`
	ResourceCount     map[string]int                 `json:"resourceCount"`
	RoundPlayerUpdate *types.WebSocketServerResponse `json:"roundPlayerUpdate"`
	VertexUpdate      *types.WebSocketServerResponse `json:"vertexUpdate"`
}
