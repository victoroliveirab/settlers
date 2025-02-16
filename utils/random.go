package utils

import (
	"math/rand"
	"time"
)

func RandNew(seed int64) *rand.Rand {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	source := rand.NewSource(seed)
	generator := rand.New(source)
	return generator
}
