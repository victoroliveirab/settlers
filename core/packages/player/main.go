package player

import (
	"fmt"
	"maps"

	coreT "github.com/victoroliveirab/settlers/core/types"
	"github.com/victoroliveirab/settlers/utils"
)

type Instance struct {
	id                    string
	resources             map[string]int
	developmentCards      map[string][]*coreT.DevelopmentCard
	usedDevelopmentCards  map[string]int
	settlements           []int
	cities                []int
	roads                 []int
	ports                 []int
	portsTypes            []string
	points                int
	longestRoadSegments   []int
	numDevCardsPlayedTurn int
	discardAmount         int
	hasDiscardedThisRound bool
}

func New(
	player *coreT.Player,
	initialResources map[string]int,
	initialDevCards map[string][]*coreT.DevelopmentCard,
) *Instance {
	return &Instance{
		id:               player.ID,
		resources:        maps.Clone(initialResources),
		developmentCards: maps.Clone(initialDevCards),
		settlements:           make([]int, 0),
		cities:                make([]int, 0),
		roads:                 make([]int, 0),
		ports:                 make([]int, 0),
		portsTypes:            make([]string, 0),
		longestRoadSegments:   make([]int, 0),
		numDevCardsPlayedTurn: 0,
		discardAmount:         0,
		hasDiscardedThisRound: false,
	}
}

func (p *Instance) HasResourcesToBuildCity() bool {
	return p.resources["Grain"] >= 2 && p.resources["Ore"] >= 3
}

func (p *Instance) AddResource(resource string, quantity int) {
	p.resources[resource] += quantity
}

func (p *Instance) RemoveResource(resource string, quantity int) {
	p.resources[resource] -= quantity
}

func (p *Instance) AddDevelopmentCard(card *coreT.DevelopmentCard) {
	_, exists := p.developmentCards[card.Name]
	if !exists {
		p.developmentCards[card.Name] = make([]*coreT.DevelopmentCard, 0)
	}
	p.developmentCards[card.Name] = append(p.developmentCards[card.Name], card)
}

func (p *Instance) ConsumeDevelopmentCard(card *coreT.DevelopmentCard) error {
	if card == nil {
		err := fmt.Errorf("cannot consume card: nil pointer")
		return err
	}

	index := -1
	for i, entry := range p.developmentCards[card.Name] {
		if entry == card {
			index = i
			break
		}
	}

	if index == -1 {
		err := fmt.Errorf("cannot consume card: card not found")
		return err
	}
	prevDevCards := p.developmentCards[card.Name]
	p.developmentCards[card.Name] = append(prevDevCards[:index], prevDevCards[index+1:]...)

	_, exists := p.usedDevelopmentCards[card.Name]
	if !exists {
		p.usedDevelopmentCards[card.Name] = 0
	}
	p.usedDevelopmentCards[card.Name]++
	p.numDevCardsPlayedTurn++
	return nil
}

func (p *Instance) AddSettlement(vertexID int) {
	p.settlements = append(p.settlements, vertexID)
}

func (p *Instance) AddCity(vertexID int) {
	settlements := p.settlements
	for index, settlementID := range settlements {
		if settlementID != vertexID {
			continue
		}
		utils.SliceRemove(&settlements, index)
		p.settlements = settlements
		break
	}
	p.cities = append(p.cities, vertexID)
}

func (p *Instance) AddRoad(edgeID int) {
	p.roads = append(p.roads, edgeID)
}

func (p *Instance) AddPort(vertexID int, kind string) {
	p.ports = append(p.ports, vertexID)
	p.portsTypes = append(p.portsTypes, kind)
}

func (p *Instance) GetID() string {
	return p.id
}

func (p *Instance) GetResources() map[string]int {
	return maps.Clone(p.resources)
}

func (p *Instance) GetDevelopmentCards() map[string][]*coreT.DevelopmentCard {
	return maps.Clone(p.developmentCards)
}

func (p *Instance) GetNumberOfVictoryPoints() int {
	cards, exists := p.developmentCards["Victory Point"]
	if !exists {
		return 0
	}
	return len(cards)
}

func (p *Instance) GetUsedDevelopmentCards() map[string]int {
	return maps.Clone(p.usedDevelopmentCards)
}

func (p *Instance) GetSettlements() []int {
	return p.settlements
}

func (p *Instance) GetNumberOfSettlements() int {
	return len(p.settlements)
}

func (p *Instance) GetCities() []int {
	return p.cities
}

func (p *Instance) GetNumberOfCities() int {
	return len(p.cities)
}

func (p *Instance) GetRoads() []int {
	return p.roads
}

func (p *Instance) GetNumberOfRoads() int {
	return len(p.roads)
}

func (p *Instance) GetPortTypes() []string {
	return p.portsTypes
}

func (p *Instance) GetArmySize() int {
	return p.usedDevelopmentCards["Knight"]
}

func (p *Instance) GetLongestRoadSize() int {
	return len(p.longestRoadSegments)
}

func (p *Instance) SetLongestRoadSegments(segments []int) {
	p.longestRoadSegments = segments
}

func (p *Instance) GetNumberOfDevCardsPlayedCurrentTurn() int {
	return p.numDevCardsPlayedTurn
}

func (p *Instance) ResetNumberOfDevCardsPlayedCurrentTurn() {
	p.numDevCardsPlayedTurn = 0
}

func (p *Instance) GetDiscardAmount() int {
	return p.discardAmount
}

func (p *Instance) SetDiscardAmount(value int) {
	p.discardAmount = value
}

func (p *Instance) GetHasDiscardedThisTurn() bool {
	return p.hasDiscardedThisRound
}

func (p *Instance) SetHasDiscardedThisTurn(value bool) {
	p.hasDiscardedThisRound = value
}
