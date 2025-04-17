package utils

import (
	"fmt"
	"strings"
)

var resourcesOrder [5]string = [5]string{"Lumber", "Brick", "Sheep", "Grain", "Ore"}

func FormatResources(collection map[string]int) string {
	var sb strings.Builder
	for _, resource := range resourcesOrder {
		quantity, ok := collection[resource]
		if !ok || quantity == 0 {
			continue
		}
		sb.WriteString(fmt.Sprintf("[res q=%d v=%s] ", quantity, resource))
	}
	return strings.Trim(sb.String(), " ")
}
