package main

import (
	"fmt"
	"math/rand"
)

var seeds map[int]int

func hasCompletedSeeds() bool {
	for _, value := range seeds {
		if value == -1 {
			return false
		}
	}
	return true
}

func main() {
	seeds = make(map[int]int)
	seeds[2] = -1
	seeds[3] = -1
	seeds[4] = -1
	seeds[5] = -1
	seeds[6] = -1
	seeds[7] = -1
	seeds[8] = -1
	seeds[9] = -1
	seeds[10] = -1
	seeds[11] = -1
	seeds[12] = -1

	i := 1

	for {
		if hasCompletedSeeds() {
			break
		}
		s := rand.NewSource(int64(i))
		r := rand.New(s)
		i++

		dice1 := r.Intn(6) + 1
		dice2 := r.Intn(6) + 1
		sum := dice1 + dice2

		el, exists := seeds[sum]
		if !exists {
			continue
		}
		if el == -1 {
			seeds[sum] = i - 1
		}
	}

	fmt.Println("Seeds:")
	fmt.Println(seeds)

	s := rand.NewSource(int64(56))
	r := rand.New(s)
	dice1 := r.Intn(6) + 1
	dice2 := r.Intn(6) + 1
	sum := dice1 + dice2
	fmt.Println("sum:", sum)
	fmt.Println("dice1:", dice1)
	fmt.Println("dice2:", dice2)

}
