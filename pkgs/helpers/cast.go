package helpers

func SafeCast[T any](target interface{}) (bool, T) {
	value, isSuccess := target.(T)
	return isSuccess, value
}
