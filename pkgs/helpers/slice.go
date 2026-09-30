package helpers

func SliceMap[I any, O any](inputs []I, fn func(I) *O) []*O {
	outputs := make([]*O, 0, len(inputs))
	for _, item := range inputs {
		output := fn(item)
		if output == nil {
			continue
		}
		outputs = append(outputs, output)
	}
	return outputs
}
