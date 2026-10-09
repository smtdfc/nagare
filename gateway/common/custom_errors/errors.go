package custom_errors

var (
	ErrMissingRequestBody = NewGatewayError("MISSING_REQUEST_BODY", "Request body is missing", 400)
	ErrInvalidBody        = NewGatewayError("INVALID_BODY", "Request body data is invalid or missing required fields", 400)
	ErrUnauthorized       = NewGatewayError("UNAUTHORIZED", "Unauthorized access", 401)
)
