package models

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return e.Message
}

// Common error constructors
func NewBadRequestError(message string) APIError {
	return APIError{Code: 400, Message: message}
}

func NewUnauthorizedError(message string) APIError {
	return APIError{Code: 401, Message: message}
}

func NewForbiddenError(message string) APIError {
	return APIError{Code: 403, Message: message}
}

func NewNotFoundError(message string) APIError {
	return APIError{Code: 404, Message: message}
}

func NewConflictError(message string) APIError {
	return APIError{Code: 409, Message: message}
}

func NewInternalServerError(message string) APIError {
	return APIError{Code: 500, Message: message}
}

func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	apiErr, ok := err.(APIError)
	return ok && apiErr.Code == 404
}
