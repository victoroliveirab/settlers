package match

import (
	"time"

	"github.com/victoroliveirab/settlers/core/packages/board"
	"github.com/victoroliveirab/settlers/core/packages/trade"
	coreT "github.com/victoroliveirab/settlers/core/types"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

type currentRoundPlayerStateUpdateResponsePayload struct {
	Deadline    *time.Time `json:"deadline"`
	Player      string     `json:"player"`
	ServerTime  time.Time  `json:"serverNow"`
	SubDeadline *time.Time `json:"subDeadline"`
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
	BlockedTiles []int                  `json:"blockedTiles"`
	Cities       map[int]board.Building `json:"cities"`
	Roads        map[int]board.Building `json:"roads"`
	Settlements  map[int]board.Building `json:"settlements"`
}

type portStateUpdateResponsePayload struct {
	Ports []string `json:"ports"`
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

type devHandCountStateUpdateResponsePayload struct {
	DevHandCount map[string]int `json:"devHandCount"`
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
	Offers []trade.Trade `json:"offers"`
}

type pointsStateUpdate struct {
	Points map[string]int `json:"points"`
}

type longestRoadStateUpdate struct {
	LongestRoadSizeByPlayer map[string]int `json:"longestRoadSizeByPlayer"`
}

type knightUsageStateUpdate struct {
	KnightUsesByPlayer map[string]int `json:"knightUsesByPlayer"`
}

type monopolyStateUpdate struct {
	Enabled bool `json:"enabled"`
}

type yearOfPlentyStateUpdate struct {
	Enabled bool `json:"enabled"`
}

type hydrateSetupMatchResponsePayload struct {
	DevHandCount      map[string]int                 `json:"devHandCount"`
	EdgeUpdate        *types.WebSocketServerResponse `json:"edgeUpdate"`
	Map               []coreT.MapBlock               `json:"map"`
	MapName           string                         `json:"mapName"`
	MapUpdate         *types.WebSocketServerResponse `json:"mapUpdate"`
	Players           []coreT.Player                 `json:"players"`
	Ports             []coreT.Port                   `json:"ports"`
	ResourceCount     map[string]int                 `json:"resourceCount"`
	RoomStatus        string                         `json:"roomStatus"`
	RoundPlayerUpdate *types.WebSocketServerResponse `json:"roundPlayerUpdate"`
	VertexUpdate      *types.WebSocketServerResponse `json:"vertexUpdate"`
}

type hydrateOngoingMatchResponsePayload struct {
	BuyDevCardUpdate         *types.WebSocketServerResponse `json:"buyDevCardUpdate"`
	DevHandCount             map[string]int                 `json:"devHandCount"`
	DevHandUpdate            *types.WebSocketServerResponse `json:"devHandUpdate"`
	DevHandPermissionsUpdate *types.WebSocketServerResponse `json:"devHandPermissionsUpdate"`
	DiceUpdate               *types.WebSocketServerResponse `json:"diceUpdate"`
	DiscardUpdate            *types.WebSocketServerResponse `json:"discardUpdate"`
	EdgeUpdate               *types.WebSocketServerResponse `json:"edgeUpdate"`
	HandUpdate               *types.WebSocketServerResponse `json:"handUpdate"`
	KnightsUsageUpdate       *types.WebSocketServerResponse `json:"knightsUsageUpdate"`
	LongestRoadUpdate        *types.WebSocketServerResponse `json:"longestRoadUpdate"`
	Map                      []coreT.MapBlock               `json:"map"`
	MapName                  string                         `json:"mapName"`
	MapUpdate                *types.WebSocketServerResponse `json:"mapUpdate"`
	MonopolyUpdate           *types.WebSocketServerResponse `json:"monopolyUpdate"`
	PassActionState          *types.WebSocketServerResponse `json:"passActionState"`
	Players                  []coreT.Player                 `json:"players"`
	PointsUpdate             *types.WebSocketServerResponse `json:"pointsUpdate"`
	Ports                    []coreT.Port                   `json:"ports"`
	PortsUpdate              *types.WebSocketServerResponse `json:"portsUpdate"`
	ResourceCount            map[string]int                 `json:"resourceCount"`
	RobbablePlayersUpdate    *types.WebSocketServerResponse `json:"robbablePlayersUpdate"`
	RobberUpdate             *types.WebSocketServerResponse `json:"robberMovementUpdate"`
	RoomStatus               string                         `json:"roomStatus"`
	RoundPlayerUpdate        *types.WebSocketServerResponse `json:"roundPlayerUpdate"`
	TradeActionState         *types.WebSocketServerResponse `json:"tradeActionState"`
	TradeOffersUpdate        *types.WebSocketServerResponse `json:"tradeOffersUpdate"`
	VertexUpdate             *types.WebSocketServerResponse `json:"vertexUpdate"`
	YearOfPlentyUpdate       *types.WebSocketServerResponse `json:"yearOfPlentyUpdate"`
}
