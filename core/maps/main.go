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
	Tokens               []int           `json:"tokens"`
	Resources            []ResourceEntry `json:"resources"`
	VerticesByTile       map[int][6]int  `json:"verticesByTile"`
	EdgesByTile          map[int][6]int  `json:"edgesByTile"`
	HexCoordinatesByTile map[int][3]int  `json:"hexCoordinatesByTile"`
	TilesByVertex        map[int][]int   `json:"tilesByVertex"`
	VerticesByEdge       map[int][2]int  `json:"verticesByEdge"`
	EdgesByVertex        map[int][]int   `json:"edgesByVertex"`
}

type meta struct {
	Id      int
	Name    string
	Players struct {
		Min int
		Max int
	}
}

type jsonStructure struct {
	Data MapDefinition
	Meta meta
}

type generateMapReturnType struct {
	Definition MapDefinition
	Tiles      []*coreT.MapBlock
}

var MapCollection map[string]MapDefinition = make(map[string]MapDefinition)

func LoadMap(name string) error {
	filename := fmt.Sprintf("%s.json", name)
	filePath := filepath.Join("engine", "maps-definitions", "data", filename)
	action := fmt.Sprintf("LoadMap.%s", name)

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		logger.LogError("system", action, -1, err)
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read file contents
	content, err := io.ReadAll(file)
	if err != nil {
		logger.LogError("system", action, -1, err)
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Unmarshal JSON into Definition
	var data jsonStructure
	if err := json.Unmarshal(content, &data); err != nil {
		logger.LogError("system", action, -1, err)
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	MapCollection[name] = data.Data
	logger.LogMessage("system", action, "loaded map successfully")
	return nil
}

func GenerateMap(name string, rand *rand.Rand) (*generateMapReturnType, error) {
	definitions, exists := MapCollection[name]
	if !exists {
		err := fmt.Errorf("cannot generate unknown map %s", name)
		return nil, err
	}

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

	return &generateMapReturnType{
		Definition: definitions,
		Tiles:      instance,
	}, nil
}
