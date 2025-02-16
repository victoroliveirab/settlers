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

func SliceContains[T comparable](slice []T, val T) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func SliceLast[T any](slice []T) T {
	return slice[len(slice)-1]
}
