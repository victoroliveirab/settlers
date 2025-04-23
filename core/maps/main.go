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
	Tiles          []coreT.MapBlock
	RobberPosition int
	Ports          map[int]string
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

func GetMapDefinitions(name string) (*MapDefinition, error) {
	data, exists := MapCollection[name]
	if !exists {
		err := fmt.Errorf("unknown map: %s", name)
		return nil, err
	}
	return &data.Data, nil
}

func GenerateMap(definitions *MapDefinition, rand *rand.Rand) *generateMapReturnType {
	robberPosition := -1

	instance := make([]coreT.MapBlock, 0)
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
		token := 0
		// Prevents index out of range error when Desert is the last tile
		if i+desertShift < len(definitions.Tokens) {
			token = definitions.Tokens[i+desertShift]
		}
		block := coreT.MapBlock{
			Resource: resourceName,
			Token:    token,
			Blocked:  false,
		}

		if resourceName == "Desert" {
			block.Token = 0
			desertShift--
			if !alreadyBlocked {
				block.Blocked = true
				robberPosition = i
				alreadyBlocked = true
			}
		}

		instance = append(instance, block)
	}
	utils.SliceShuffle(instance, rand)

	for index := range instance {
		tileID := index + 1
		instance[index].ID = tileID
		instance[index].Vertices = definitions.VerticesByTile[tileID]
		instance[index].Edges = definitions.EdgesByTile[tileID]
		instance[index].Coordinates = coreT.HexCoordinate{
			Q: definitions.HexCoordinatesByTile[tileID][0],
			R: definitions.HexCoordinatesByTile[tileID][1],
			S: definitions.HexCoordinatesByTile[tileID][2],
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

	return &generateMapReturnType{
		RobberPosition: robberPosition,
		Ports:          ports,
		Tiles:          instance,
	}
}

func GetMetadata(mapName string) (*meta, error) {
	data, exists := MapCollection[mapName]
	if !exists {
		err := fmt.Errorf("%s doesn't exist in memory", mapName)
		return nil, err
	}
	return &data.Meta, nil
}
