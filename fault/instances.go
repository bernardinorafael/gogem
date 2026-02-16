package fault

import (
	"net/http"
)

func NewValidation(message string, err error) *Fault {
	return New(
		message,
		WithHTTPCode(http.StatusUnprocessableEntity),
		WithTag(ValidationError),
		WithValidationError(err),
	)
}

func NewBadRequest(message string) *Fault {
	return New(
		message,
		WithHTTPCode(http.StatusBadRequest),
		WithTag(BadRequest),
	)
}

func NewNotFound(message string) *Fault {
	return New(
		message,
		WithHTTPCode(http.StatusNotFound),
		WithTag(NotFound),
	)
}

func NewInternalServerError(message string) *Fault {
	return New(
		message,
		WithHTTPCode(http.StatusInternalServerError),
		WithTag(InternalServerError),
	)
}

func NewUnauthorized(message string) *Fault {
	return New(
		message,
		WithHTTPCode(http.StatusUnauthorized),
		WithTag(Unauthorized),
	)
}

func NewForbidden(message string) *Fault {
	return New(
		message,
		WithHTTPCode(http.StatusForbidden),
		WithTag(Forbidden),
	)
}

func NewConflict(message string) *Fault {
	return New(
		message,
		WithHTTPCode(http.StatusConflict),
		WithTag(Conflict),
	)
}

func NewTooManyRequests(message string) *Fault {
	return New(
		message,
		WithHTTPCode(http.StatusTooManyRequests),
		WithTag(TooManyRequests),
	)
}

func NewUnprocessableEntity(message string) *Fault {
	return New(
		message,
		WithHTTPCode(http.StatusUnprocessableEntity),
		WithTag(UnprocessableEntity),
	)
}
