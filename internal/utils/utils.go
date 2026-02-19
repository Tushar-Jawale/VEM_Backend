package utils

import (
	"strings"
)

// ParseCSV splits a comma-separated string into a slice of strings
// e.g. "a,b,c" -> ["a", "b", "c"]
func ParseCSV(data string) []string {
	if data == "" {
		return []string{}
	}
	parts := strings.Split(data, ",")
	var result []string
	for _, p := range parts {
		if val := strings.TrimSpace(p); val != "" {
			result = append(result, val)
		}
	}
	return result
}

// ParseKeyVal splits a string like "key:val,key2:val2" into a map
func ParseKeyVal(data string) map[string]string {
	result := make(map[string]string)
	if data == "" {
		return result
	}
	pairs := strings.Split(data, ",")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, ":", 2)
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			val := strings.TrimSpace(kv[1])
			if key != "" {
				result[key] = val
			}
		}
	}
	return result
}
