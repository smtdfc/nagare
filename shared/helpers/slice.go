package helpers

func Map[T any, R any](collection []T, iteratee func(T) (R, error)) ([]R, error) {
	result := make([]R, len(collection))
	for i, item := range collection {
		r, err := iteratee(item)
		if err != nil {
			return nil, err
		}
		result[i] = r
	}
	return result, nil
}
