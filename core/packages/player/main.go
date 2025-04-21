package player

import (
	"fmt"

	coreT "github.com/victoroliveirab/settlers/core/types"
	"github.com/victoroliveirab/settlers/utils"
)

type Instance struct {
	ID                    string
	Resources             map[string]int
	DevelopmentCards      map[string][]*coreT.DevelopmentCard
	UsedDevelopmentCards  map[string]int
	Settlements           []int
	Cities                []int
	Roads                 []int
	Ports                 []int
	PortsTypes            []string
	Points                int
	KnightCount           int
	LongestRoadSegments   []int
	NumDevCardsPlayedTurn int
	DiscardAmount         int
	HasDiscardedThisRound bool
}

// REFACTOR: for when we introduce more dev cards
// func New(player *coreT.Player, availableDevCards []string) *Instance {
func New(player *coreT.Player) *Instance {
	resources := map[string]int{
		"Lumber": 0,
		"Brick":  0,
		"Grain":  0,
		"Sheep":  0,
		"Ore":    0,
	}
	return &Instance{
		ID:                    player.ID,
		Resources:             resources,
		DevelopmentCards:      make(map[string][]*coreT.DevelopmentCard),
		UsedDevelopmentCards:  make(map[string]int),
		Settlements:           make([]int, 0),
		Cities:                make([]int, 0),
		Roads:                 make([]int, 0),
		Ports:                 make([]int, 0),
		PortsTypes:            make([]string, 0),
		Points:                0,
		KnightCount:           0,
		LongestRoadSegments:   make([]int, 0),
		NumDevCardsPlayedTurn: 0,
		DiscardAmount:         0,
		HasDiscardedThisRound: false,
	}
}

func (p *Instance) HasResourcesToBuildCity() bool {
	return p.Resources["Grain"] >= 2 && p.Resources["Ore"] >= 3
}

func (p *Instance) AddResource(resource string, quantity int) {
	p.Resources[resource] += quantity
}

func (p *Instance) RemoveResource(resource string, quantity int) {
	p.Resources[resource] -= quantity
}

func (p *Instance) AddDevelopmentCard(card *coreT.DevelopmentCard) {
	_, exists := p.DevelopmentCards[card.Name]
	if !exists {
		p.DevelopmentCards[card.Name] = make([]*coreT.DevelopmentCard, 0)
	}
	p.DevelopmentCards[card.Name] = append(p.DevelopmentCards[card.Name], card)
}

func (p *Instance) ConsumeDevelopmentCard(card *coreT.DevelopmentCard) error {
	if card == nil {
		err := fmt.Errorf("cannot consume card: nil pointer")
		return err
	}

	index := -1
	for i, entry := range p.DevelopmentCards[card.Name] {
		if entry == card {
			index = i
			break
		}
	}

	if index == -1 {
		err := fmt.Errorf("cannot consume card: card not found")
		return err
	}
	prevDevCards := p.DevelopmentCards[card.Name]
	p.DevelopmentCards[card.Name] = append(prevDevCards[:index], prevDevCards[index+1:]...)

	_, exists := p.UsedDevelopmentCards[card.Name]
	if !exists {
		p.UsedDevelopmentCards[card.Name] = 0
	}
	p.UsedDevelopmentCards[card.Name]++
	p.NumDevCardsPlayedTurn++
	return nil
}

func (p *Instance) AddSettlement(vertexID int) {
	p.Settlements = append(p.Settlements, vertexID)
}

func (p *Instance) AddCity(vertexID int) {
	settlements := p.Settlements
	for index, settlementID := range settlements {
		if settlementID != vertexID {
			continue
		}
		utils.SliceRemove(&settlements, index)
		p.Settlements = settlements
		break
	}
	p.Cities = append(p.Cities, vertexID)
}

func (p *Instance) AddRoad(edgeID int) {
	p.Roads = append(p.Roads, edgeID)
}

func (p *Instance) AddPort(vertexID int, kind string) {
	p.Ports = append(p.Ports, vertexID)
	p.PortsTypes = append(p.PortsTypes, kind)
}
