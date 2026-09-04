package rest

var (
	NotFoundErr = NewApiError("NOT_FOUND", "Not found", 404)
	InternalErr = NewApiError("INTERNAL_ERR", "Internal error", 500)
)
