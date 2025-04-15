package match

import (
	"fmt"
	"strconv"
	"strings"
)

var resourcesOrder [5]string = [5]string{"Lumber", "Brick", "Sheep", "Grain", "Ore"}

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

func serializeHandDiff(hand map[string]int) string {
	var sb strings.Builder
	first := true
	for resource, quantity := range hand {
		if quantity != 0 {
			if !first {
				sb.WriteString(", ")
			}
			sb.WriteString(strconv.FormatInt(int64(quantity), 10))
			sb.WriteString(" ")
			sb.WriteString(resource)
		}
	}
	return sb.String()
}

func formatResourceCollection(collection map[string]int) string {
	var sb strings.Builder
	for _, resource := range resourcesOrder {
		quantity, ok := collection[resource]
		if !ok || quantity == 0 {
			continue
		}
		sb.WriteString(fmt.Sprintf(" [res q=%d]%s[/res]", quantity, resource))
	}
	return strings.Trim(sb.String(), " ")
}
