package match

import (
	"fmt"

	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func diffResourceHands(before, after map[string]int) (map[string]int, error) {
	diff := map[string]int{}
	for resource, afterQuantity := range after {
		beforeQuantity, exists := before[resource]
		if !exists {
			return nil, fmt.Errorf("cannot compare two different shapes of hands - mismatch in key %s", resource)
		}
		diff[resource] = afterQuantity - beforeQuantity
	}

	return diff, nil
}

func hasDiff(diff map[string]int) bool {
	for _, quantity := range diff {
		if quantity != 0 {
			return true
		}
	}
	return false
}

// TODO: get rid of this and use the utils fn directly
func formatResourceCollection(collection map[string]int) string {
	return utils.FormatResources(collection)
}
