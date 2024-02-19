package database

import (
	"fmt"
	"strings"
)

// ConvertMapToArgsStr converts a map to a string of args for a named query.
func ConvertMapToArgsStr(filters map[string]any, separator string) string {
	args := []string{}
	for key := range filters {
		args = append(args, fmt.Sprintf("%s = :%s", key, key))
	}
	return strings.Join(args, separator)
}
