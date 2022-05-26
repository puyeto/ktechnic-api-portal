package errors

import (
	"errors"
	"net/http"
	"sort"

	validation "github.com/go-ozzo/ozzo-validation"
)

type validationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// New ...
func New(err string) error {
	return errors.New(err)
}

// InternalServerError creates a new API error representing an internal server error (HTTP 500)
func InternalServerError(err string) *APIError {
	return NewAPIError(http.StatusInternalServerError, err, Params{"error": err})
}

// NotFound creates a new API error representing a resource-not-found error (HTTP 404)
func NotFound(resource string) *APIError {
	return NewAPIError(http.StatusNotFound, "NOT_FOUND : "+resource, Params{"resource": resource})
}

// Unauthorized creates a new API error representing an authentication failure (HTTP 401)
func Unauthorized(err string) *APIError {
	return NewAPIError(http.StatusUnauthorized, "UNAUTHORIZED : "+err, Params{"error": err})
}

// BadRequest creates a new API error representing a bad request (HTTP 400)
func BadRequest(err string) *APIError {
	return NewAPIError(http.StatusBadRequest, "BADREQUEST : "+err, Params{"error": err})
}

// ValidationRequest ...
func ValidationRequest(err string) *APIError {
	return NewAPIError(http.StatusBadRequest, "VALIDATION_FAILED : "+err, Params{"error": err})
}

// NotFound creates a new API error representing a resource-not-found error (HTTP 204)
func NoContentFound(err string) *APIError {
	return NewAPIError(http.StatusNoContent, err, Params{"error": err})
}

// InvalidData converts a data validation error into an API error (HTTP 400)
func InvalidData(errs validation.Errors) *APIError {
	result := []validationError{}
	fields := []string{}
	for field := range errs {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	for _, field := range fields {
		err := errs[field]
		result = append(result, validationError{
			Field: field,
			Error: err.Error(),
		})
	}

	err := NewAPIError(http.StatusBadRequest, "INVALID_DATA", nil)
	err.Details = result

	return err
}
