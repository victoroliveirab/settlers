package match

import (
	"github.com/victoroliveirab/settlers/core"
	coreT "github.com/victoroliveirab/settlers/core/types"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

type currentRoundPlayerStateUpdateResponsePayload struct {
	Player string `json:"player"`
}

type verticesStateUpdateResponsePayload struct {
	AvailableCityVertices       []int `json:"availableCityVertices"`
	AvailableSettlementVertices []int `json:"availableSettlementVertices"`
	Enabled                     bool  `json:"enabled"`
	Highlight                   bool  `json:"highlight"`
}

type edgesStateUpdateResponsePayload struct {
	AvailableEdges []int `json:"availableEdges"`
	Enabled        bool  `json:"enabled"`
	Highlight      bool  `json:"highlight"`
}

type mapStateUpdateResponsePayload struct {
	BlockedTiles []int                 `json:"blockedTiles"`
	Cities       map[int]core.Building `json:"cities"`
	Roads        map[int]core.Building `json:"roads"`
	Settlements  map[int]core.Building `json:"settlements"`
}

type diceStateUpdateResponsePayload struct {
	Dice    [2]int `json:"dice"`
	Enabled bool   `json:"enabled"`
}

type handStateUpdateResponsePayload struct {
	Hand map[string]int `json:"hand"`
}

type devHandStateUpdateResponsePayload struct {
	DevHand map[string]int `json:"devHand"`
}

type devHandPermissionsStateUpdateResponsePayload struct {
	DevHandPermissions map[string]bool `json:"devHandPermissions"`
}

type resourceCountStateUpdateResponsePayload struct {
	ResourceCount map[string]int `json:"resourceCount"`
}

type moveRobberStateUpdateResponsePayload struct {
	AvailableTiles []int `json:"availableTiles"`
	Enabled        bool  `json:"enabled"`
	Highlight      bool  `json:"highlight"`
}

type discardPhaseStateUpdateResponsePayload struct {
	DiscardAmounts map[string]int `json:"discardAmounts"`
	Enabled        bool           `json:"enabled"`
}

type passStateUpdateResponsePayload struct {
	Enabled bool `json:"enabled"`
}

type startTradeStateUpdateResponsePayload struct {
	Enabled bool `json:"enabled"`
}

type buyDevCardStateUpdateResponsePayload struct {
	Enabled bool `json:"enabled"`
}

type pickRobbedStateUpdate struct {
	Enabled bool     `json:"enabled"`
	Options []string `json:"options"`
}

type updateActiveTradeOffersStateUpdate struct {
	Offers []core.Trade `json:"offers"`
}

type hydrateSetupMatchResponsePayload struct {
	EdgeUpdate        *types.WebSocketServerResponse `json:"edgeUpdate"`
	Map               []*coreT.MapBlock              `json:"map"`
	MapName           string                         `json:"mapName"`
	MapUpdate         *types.WebSocketServerResponse `json:"mapUpdate"`
	Players           []coreT.Player                 `json:"players"`
	Ports             map[int]string                 `json:"ports"`
	ResourceCount     map[string]int                 `json:"resourceCount"`
	RoundPlayerUpdate *types.WebSocketServerResponse `json:"roundPlayerUpdate"`
	VertexUpdate      *types.WebSocketServerResponse `json:"vertexUpdate"`
}

type hydrateOngoingMatchResponsePayload struct {
	BuyDevCardUpdate         *types.WebSocketServerResponse `json:"buyDevCardUpdate"`
	DevHandUpdate            *types.WebSocketServerResponse `json:"devHandUpdate"`
	DevHandPermissionsUpdate *types.WebSocketServerResponse `json:"devHandPermissionsUpdate"`
	DiceUpdate               *types.WebSocketServerResponse `json:"diceUpdate"`
	DiscardUpdate            *types.WebSocketServerResponse `json:"discardUpdate"`
	EdgeUpdate               *types.WebSocketServerResponse `json:"edgeUpdate"`
	HandUpdate               *types.WebSocketServerResponse `json:"handUpdate"`
	Map                      []*coreT.MapBlock              `json:"map"`
	MapName                  string                         `json:"mapName"`
	MapUpdate                *types.WebSocketServerResponse `json:"mapUpdate"`
	PassActionState          *types.WebSocketServerResponse `json:"passActionState"`
	Players                  []coreT.Player                 `json:"players"`
	Ports                    map[int]string                 `json:"ports"`
	ResourceCount            map[string]int                 `json:"resourceCount"`
	RobbablePlayersUpdate    *types.WebSocketServerResponse `json:"robbablePlayersUpdate"`
	RobberUpdate             *types.WebSocketServerResponse `json:"robberMovementUpdate"`
	RoundPlayerUpdate        *types.WebSocketServerResponse `json:"roundPlayerUpdate"`
	TradeActionState         *types.WebSocketServerResponse `json:"tradeActionState"`
	TradeOffersUpdate        *types.WebSocketServerResponse `json:"tradeOffersUpdate"`
	VertexUpdate             *types.WebSocketServerResponse `json:"vertexUpdate"`
}
