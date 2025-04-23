package utils

import (
	"fmt"
	"strings"

	"github.com/victoroliveirab/settlers/core"
)

func FormatResources(collection map[string]int) string {
	var sb strings.Builder
	for _, resource := range core.ResourcesOrder {
		quantity, ok := collection[resource]
		if !ok || quantity == 0 {
			continue
		}
		sb.WriteString(fmt.Sprintf("[res q=%d v=%s] ", quantity, resource))
	}
	return strings.Trim(sb.String(), " ")
}
