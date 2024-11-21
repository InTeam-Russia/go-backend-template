package apierr

type ApiError struct {
	Error string `json:"error"`
}

func NewApiError(error string) ApiError {
	return ApiError{Error: error}
}

var InternalServerError = ApiError{
	Error: "INTERNAL_SERVER",
}

var InvalidJsonError = ApiError{
	Error: "INVALID_JSON",
}

var NotFoundError = ApiError{
	Error: "NOT_FOUND",
}

var WrongCredentials = ApiError{
	Error: "WRONG_CREDENTIALS",
}
