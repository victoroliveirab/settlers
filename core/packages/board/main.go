package board

import (
	"fmt"
	"maps"
	"math/rand"

	coreMaps "github.com/victoroliveirab/settlers/core/maps"
	coreT "github.com/victoroliveirab/settlers/core/types"
)

type Building struct {
	ID    int    `json:"id"`
	Owner string `json:"owner"`
}

type Instance struct {
	cities         map[int]Building
	Definition     *coreMaps.MapDefinition
	MapName        string
	Ports          map[int]string
	roads          map[int]Building
	RobberLocation int
	settlements    map[int]Building
	tiles          []coreT.MapBlock
}

func New(mapName string, definitions *coreMaps.MapDefinition, randGenerator *rand.Rand) *Instance {
	data := coreMaps.GenerateMap(definitions, randGenerator)
	b := &Instance{
		cities:         make(map[int]Building),
		Definition:     definitions,
		MapName:        mapName,
		Ports:          data.Ports,
		roads:          make(map[int]Building),
		RobberLocation: data.RobberPosition,
		settlements:    make(map[int]Building),
		tiles:          data.Tiles,
	}
	return b
}

func (b *Instance) AddCity(playerID string, vertexID int) {
	delete(b.settlements, vertexID)
	b.cities[vertexID] = Building{Owner: playerID, ID: vertexID}
}

func (b *Instance) AddRoad(playerID string, edgeID int) {
	b.roads[edgeID] = Building{Owner: playerID, ID: edgeID}
}

func (b *Instance) AddSettlement(playerID string, vertexID int) {
	b.settlements[vertexID] = Building{Owner: playerID, ID: vertexID}
}

func (b *Instance) GetCities() map[int]Building {
	return maps.Clone(b.cities)
}

func (b *Instance) GetRoads() map[int]Building {
	return maps.Clone(b.roads)
}

func (b *Instance) GetSettlements() map[int]Building {
	return maps.Clone(b.settlements)
}

func (b *Instance) GetTiles() []coreT.MapBlock {
	tiles := make([]coreT.MapBlock, len(b.tiles))
	copy(tiles, b.tiles)
	return tiles
}

func (b *Instance) GetBlockedTilesIDs() []int {
	tileIDs := make([]int, 0)
	for _, tile := range b.tiles {
		if tile.Blocked {
			tileIDs = append(tileIDs, tile.ID)
		}
	}
	return tileIDs
}

func (b *Instance) GetUnblockedTilesIDs() []int {
	tileIDs := make([]int, 0)
	for _, tile := range b.tiles {
		if !tile.Blocked {
			tileIDs = append(tileIDs, tile.ID)
		}
	}
	return tileIDs
}

func (b *Instance) BlockTileByIndex(index int) error {
	if b.tiles[index].Blocked {
		err := fmt.Errorf("Cannot block tile at index %d: already blocked", index)
		return err
	}
	b.tiles[index].Blocked = true
	return nil
}

func (b *Instance) BlockTile(tileID int) error {
	for i, tile := range b.tiles {
		if tile.ID == tileID {
			err := b.BlockTileByIndex(i)
			if err != nil {
				err := fmt.Errorf("Cannot block tile#%d: already blocked", tileID)
				return err
			}
			return nil
		}
	}
	err := fmt.Errorf("Cannot block tile#%d: tile not found", tileID)
	return err
}

func (b *Instance) UnblockTileByIndex(index int) error {
	if !b.tiles[index].Blocked {
		err := fmt.Errorf("Cannot unblock tile at index %d: not blocked", index)
		return err
	}
	b.tiles[index].Blocked = false
	return nil
}

func (b *Instance) UnblockTile(tileID int) error {
	for i, tile := range b.tiles {
		if tile.ID == tileID {
			err := b.UnblockTileByIndex(i)
			if err != nil {
				err := fmt.Errorf("Cannot unblock tile#%d: not blocked", tileID)
				return err
			}
			return nil
		}
	}
	err := fmt.Errorf("Cannot unblock tile#%d: tile not found", tileID)
	return err
}
