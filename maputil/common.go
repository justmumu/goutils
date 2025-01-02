package maputil

import (
	extmaps "golang.org/x/exp/maps"
)

// GetKeys returns the map's keys.
func GetKeys[K comparable, V any](maps ...map[K]V) []K {
	var keys []K
	for _, m := range maps {
		keys = append(keys, extmaps.Keys(m)...)
	}
	return keys
}

// GetValues returns the map's values.
func GetValues[K comparable, V any](maps ...map[K]V) []V {
	var values []V
	for _, m := range maps {
		values = append(values, extmaps.Values(m)...)
	}
	return values
}

// Difference returns the inputted map without the keys specified as input.
func Difference[K comparable, V any](m map[K]V, keys ...K) map[K]V {
	for _, key := range keys {
		delete(m, key)
	}

	return m
}

// Flatten takes a map and returns a new one where nested maps are replaced
// by dot-delimited keys.
func Flatten(m map[string]any, separator string) map[string]any {
	if separator == "" {
		separator = "."
	}
	o := make(map[string]any)
	for k, v := range m {
		switch child := v.(type) {
		case map[string]any:
			nm := Flatten(child, separator)
			for nk, nv := range nm {
				o[k+separator+nk] = nv
			}
		default:
			o[k] = v
		}
	}
	return o
}

// Walk a map and visit all the edge key:value pairs
func Walk(m map[string]any, callback func(k string, v any)) {
	for k, v := range m {
		switch child := v.(type) {
		case map[string]any:
			Walk(child, callback)
		default:
			callback(k, v)
		}
	}
}
