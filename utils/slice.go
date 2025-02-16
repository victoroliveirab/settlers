package utils

import (
	"math/rand"
)

func SliceShuffle[T any](slice []T, randGenerator *rand.Rand) {
	for i := len(slice) - 1; i > 0; i-- {
		j := randGenerator.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}
