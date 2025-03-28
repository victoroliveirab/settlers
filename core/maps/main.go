package maps

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"

	coreT "github.com/victoroliveirab/settlers/core/types"
	"github.com/victoroliveirab/settlers/logger"
	"github.com/victoroliveirab/settlers/utils"
)

type ResourceEntry struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type MapDefinition struct {
	Tiles                []int           `json:"tiles"`
	PortsLocations       [][2]int        `json:"portsLocations"`
	PortsByDefinition    map[string]int  `json:"portsByDefinition"`
	Tokens               []int           `json:"tokens"`
	Resources            []ResourceEntry `json:"resources"`
	VerticesByTile       map[int][6]int  `json:"verticesByTile"`
	EdgesByTile          map[int][6]int  `json:"edgesByTile"`
	HexCoordinatesByTile map[int][3]int  `json:"hexCoordinatesByTile"`
	TilesByVertex        map[int][]int   `json:"tilesByVertex"`
	VerticesByEdge       map[int][2]int  `json:"verticesByEdge"`
	EdgesByVertex        map[int][]int   `json:"edgesByVertex"`
	DevelopmentCards     map[string]int  `json:"developmentCards"`
}

type meta struct {
	Id      int
	Name    string
	Players struct {
		Min int
		Max int
	}
	Params map[string]struct {
		Default     int
		Description string
		Label       string
		Priority    int
		Values      []int
	}
}

type jsonStructure struct {
	Data MapDefinition
	Meta meta
}

type generateMapReturnType struct {
	Definition       MapDefinition
	Tiles            []*coreT.MapBlock
	Ports            map[int]string
	DevelopmentCards []*coreT.DevelopmentCard
}

var MapCollection map[string]jsonStructure = make(map[string]jsonStructure)

func LoadMap(name string) error {
	filename := fmt.Sprintf("%s.json", name)
	filePath := filepath.Join("core", "maps", "data", filename)
	action := fmt.Sprintf("LoadMap.%s", name)

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		logger.LogSystemError(action, -1, err)
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read file contents
	content, err := io.ReadAll(file)
	if err != nil {
		logger.LogSystemError(action, -1, err)
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Unmarshal JSON into Definition
	var data jsonStructure
	if err := json.Unmarshal(content, &data); err != nil {
		logger.LogSystemError(action, -1, err)
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	MapCollection[name] = data
	logger.LogSystemMessage(action, "loaded map successfully")
	return nil
}

func GenerateMap(name string, rand *rand.Rand) (*generateMapReturnType, error) {
	data, exists := MapCollection[name]
	if !exists {
		err := fmt.Errorf("cannot generate unknown map %s", name)
		return nil, err
	}

	definitions := data.Data

	instance := make([]*coreT.MapBlock, 0)
	resourcesLeft := make([]int, 0)
	typesOfResources := len(definitions.Resources)

	for _, resourceEntry := range definitions.Resources {
		resourcesLeft = append(resourcesLeft, resourceEntry.Count)
	}

	var resourceIndex int
	var resourceName string
	var desertShift int = 0

	alreadyBlocked := false

	for i := range definitions.Tiles {
		for {
			resourceIndex = rand.Intn(typesOfResources)
			if resourcesLeft[resourceIndex] > 0 {
				break
			}
		}
		resourcesLeft[resourceIndex]--

		resourceName = definitions.Resources[resourceIndex].Name
		token := definitions.Tokens[i+desertShift]
		instance = append(instance, &coreT.MapBlock{
			Resource: resourceName,
			Token:    token,
			Blocked:  false,
		})

		if resourceName == "Desert" {
			desertShift--
			if !alreadyBlocked {
				instance[i].Blocked = true
				alreadyBlocked = true
			}
		}
	}
	utils.SliceShuffle(instance, rand)

	for index, tile := range instance {
		tile.ID = index + 1
		tile.Vertices = definitions.VerticesByTile[tile.ID]
		tile.Edges = definitions.EdgesByTile[tile.ID]
		tile.Coordinates = coreT.HexCoordinate{
			Q: definitions.HexCoordinatesByTile[tile.ID][0],
			R: definitions.HexCoordinatesByTile[tile.ID][1],
			S: definitions.HexCoordinatesByTile[tile.ID][2],
		}
	}

	// NOTE: this is done to enforce ordering for tests (math/random seed)
	portsDefinitions := MapToShuffledSlice(
		definitions.PortsByDefinition,
		func(el string) string { return el },
		rand,
	)
	portsVertices := definitions.PortsLocations
	ports := make(map[int]string)
	for index, port := range portsVertices {
		vertex1 := port[0]
		vertex2 := port[1]
		ports[vertex1] = portsDefinitions[index]
		ports[vertex2] = portsDefinitions[index]
	}

	developmentCards := MapToShuffledSlice[*coreT.DevelopmentCard](
		definitions.DevelopmentCards,
		func(el string) *coreT.DevelopmentCard { return &coreT.DevelopmentCard{Name: el} },
		rand,
	)

	return &generateMapReturnType{
		Definition:       definitions,
		DevelopmentCards: developmentCards,
		Ports:            ports,
		Tiles:            instance,
	}, nil
}

func GetMetadata(mapName string) (*meta, error) {
	data, exists := MapCollection[mapName]
	if !exists {
		err := fmt.Errorf("%s doesn't exist in memory", mapName)
		return nil, err
	}
	return &data.Meta, nil
}
