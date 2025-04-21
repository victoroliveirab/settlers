package board

import (
	"math/rand"

	coreMaps "github.com/victoroliveirab/settlers/core/maps"
	coreT "github.com/victoroliveirab/settlers/core/types"
)

type Building struct {
	ID    int    `json:"id"`
	Owner string `json:"owner"`
}

type Instance struct {
	Cities         map[int]Building
	Definition     *coreMaps.MapDefinition
	MapName        string
	Ports          map[int]string
	Roads          map[int]Building
	RobberLocation int
	Settlements    map[int]Building
	Tiles          []*coreT.MapBlock
}

func New(mapName string, definitions *coreMaps.MapDefinition, randGenerator *rand.Rand) *Instance {
	data := coreMaps.GenerateMap(definitions, randGenerator)
	b := &Instance{
		Cities:         make(map[int]Building),
		Definition:     definitions,
		MapName:        mapName,
		Ports:          data.Ports,
		Roads:          make(map[int]Building),
		RobberLocation: data.RobberPosition,
		Settlements:    make(map[int]Building),
		Tiles:          data.Tiles,
	}
	return b
}

func (b *Instance) AddCity(playerID string, vertexID int) {
	delete(b.Settlements, vertexID)
	b.Cities[vertexID] = Building{Owner: playerID, ID: vertexID}
}

func (b *Instance) AddRoad(playerID string, edgeID int) {
	b.Roads[edgeID] = Building{Owner: playerID, ID: edgeID}
}

func (b *Instance) AddSettlement(playerID string, vertexID int) {
	b.Settlements[vertexID] = Building{Owner: playerID, ID: vertexID}
}
