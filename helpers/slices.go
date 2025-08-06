package helpers

// SliceToUniqMap convert slice to map with T key and struct as value
func SliceToUniqMap[T comparable](t []T) map[T]struct{} {
	m := make(map[T]struct{}, len(t))

	for _, val := range t {
		m[val] = struct{}{}
	}

	return m
}
