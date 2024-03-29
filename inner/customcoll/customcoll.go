package customcoll

func Filter[T any](slice []T, filterFunc func(T) bool) []T {
	var result []T
	for _, value := range slice {
		if filterFunc(value) {
			result = append(result, value)
		}
	}
	return result
}

func Contains[T comparable](slice []T, value T) (bool, int) {
	for i, item := range slice {
		if item == value {
			return true, i
		}
	}
	return false, -1
}

func MapKeys[T comparable, U any](m map[T]U) []T {
	keys := make([]T, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}
