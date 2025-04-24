//go:build test

package player

import (
	coreT "github.com/victoroliveirab/settlers/core/types"
)

func (p *Instance) SetResources(resources map[string]int) {
	p.resources = resources
}

func (p *Instance) SetDevelopmentCards(cards map[string][]*coreT.DevelopmentCard) {
	p.developmentCards = cards
}

func (p *Instance) SetUsedDevelopmentCards(cards map[string]int) {
	p.usedDevelopmentCards = cards
}
