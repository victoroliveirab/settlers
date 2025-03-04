package utils

import (
	"strings"
)

func FormatResources(resources map[string]int) string {
	var builder strings.Builder
	for resource, quantity := range resources {
		for i := 0; i < quantity; i++ {
			builder.WriteString("[res]")
			builder.WriteString(resource)
			builder.WriteString("[/res]")
		}
	}
	return builder.String()
}
