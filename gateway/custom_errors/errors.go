package custom_errors

var (
	ErrMissingRequestBody = NewGatewayError("MISSING_REQUEST_BODY", "Request body is missing", 400)
	ErrInvalidJsonBody    = NewGatewayError("INVALID_JSON_BODY", "Invalid JSON format in request body", 400)
	ErrInvalidBody        = NewGatewayError("INVALID_BODY", "Request body data is invalid or missing required fields", 400)
	ErrValidationFailed   = NewGatewayError("VALIDATION_FAILED", "Request data validation failed", 400)
	ErrUnauthorized       = NewGatewayError("UNAUTHORIZED", "Unauthorized access", 401)
	ErrForbidden          = NewGatewayError("FORBIDDEN", "Access forbidden", 403)

	ErrInternalServer     = NewGatewayError("INTERNAL_SERVER_ERROR", "Internal server error", 500)
	ErrServiceUnavailable = NewGatewayError("SERVICE_UNAVAILABLE", "Service temporarily unavailable", 503)
)
