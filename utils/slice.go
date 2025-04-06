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

func SliceFindIndex[T comparable](slice []T, callback func(val T) bool) int {
	for i, v := range slice {
		if callback(v) {
			return i
		}
	}
	return -1
}

func SliceLast[T any](slice []T) T {
	return slice[len(slice)-1]
}

func SliceEqual[T comparable](slice1 []T, slice2 []T) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i, itemSlice1 := range slice1 {
		if itemSlice1 != slice2[i] {
			return false
		}
	}
	return true
}

func SliceRemove[T any](slice *[]T, index int) {
	if index < 0 || index >= len(*slice) {
		return
	}
	*slice = append((*slice)[:index], (*slice)[index+1:]...)
}
