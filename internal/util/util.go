package util

func MapContains[K comparable, V any](container map[K]V, key K) bool {
	_, ok := container[key]
	return ok
}
