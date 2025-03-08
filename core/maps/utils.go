package maps

import (
	"math/rand"
	"sort"

	"github.com/victoroliveirab/settlers/utils"
)

func MapToShuffledSlice[T any](instance map[string]int, transformer func(el string) T, rand *rand.Rand) []T {
	keys := make([]string, 0)
	for key := range instance {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	slice := make([]T, 0)
	for _, key := range keys {
		quantity := instance[key]
		for i := 0; i < quantity; i++ {
			slice = append(slice, transformer(key))
		}
	}
	utils.SliceShuffle(slice, rand)
	return slice
}
