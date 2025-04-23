package development

import (
	"fmt"

	coreT "github.com/victoroliveirab/settlers/core/types"
)

type Instance struct {
	cards         []*coreT.DevelopmentCard
	nextCardIndex int // Index of the *next* card to be drawn
}

func New(cards []*coreT.DevelopmentCard) *Instance {
	return &Instance{
		cards:         cards,
		nextCardIndex: 0,
	}
}

func (d *Instance) Draw() (*coreT.DevelopmentCard, error) {
	if d.IsEmpty() {
		return nil, fmt.Errorf("cannot draw card: deck is empty")
	}
	card := d.cards[d.nextCardIndex]
	d.nextCardIndex++
	return card, nil
}

func (d *Instance) IsEmpty() bool {
	return d.nextCardIndex >= len(d.cards)
}

func (d *Instance) Remaining() int {
	return len(d.cards) - d.nextCardIndex
}
