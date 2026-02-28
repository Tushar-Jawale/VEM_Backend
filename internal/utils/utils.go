package utils

import (
	"strconv"
	"strings"
)

// ParseInt converts a string to an integer, returning a default value on error
func ParseInt(data string, defaultVal int) int {
	if data == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(data)
	if err != nil {
		return defaultVal
	}
	return val
}

// ParseFloat64 converts a string to a float64, returning a default value on error
func ParseFloat64(data string, defaultVal float64) float64 {
	if data == "" {
		return defaultVal
	}
	val, err := strconv.ParseFloat(data, 64)
	if err != nil {
		return defaultVal
	}
	return val
}

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

// MinInts returns the minimum integer in a slice
func MinInts(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	min := nums[0]
	for _, n := range nums {
		if n < min {
			min = n
		}
	}
	return min
}
