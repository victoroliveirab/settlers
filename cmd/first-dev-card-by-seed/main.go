package main

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core"
	"github.com/victoroliveirab/settlers/core/types"
	"github.com/victoroliveirab/settlers/utils"
)

func main() {
	for i := 1; i < 50; i++ {
		randGenerator := utils.RandNew(int64(i))
		var getter func() []*types.DevelopmentCard
		core.CreateTestGameWithRand(randGenerator, core.MockWithPeekDevCards(&getter))
		devCards := getter()
		fmt.Printf("Seed: %d, First Dev Card: %s\n", i, devCards[0].Name)
	}
}
