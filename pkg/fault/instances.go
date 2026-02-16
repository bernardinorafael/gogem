package fault

import (
	"net/http"
)

func NewValidation(message string, fields ...FieldError) *Fault {
	return New(
		message,
		WithHTTPCode(http.StatusUnprocessableEntity),
		WithTag(ValidationError),
		WithFieldError(fields...),
	)
}

func NewBadRequest(message string, options ...func(*Fault)) *Fault {
	defaults := []func(*Fault){
		WithHTTPCode(http.StatusBadRequest),
		WithTag(BadRequest),
	}
	return New(message, append(defaults, options...)...)
}

func NewNotFound(message string, options ...func(*Fault)) *Fault {
	defaults := []func(*Fault){
		WithHTTPCode(http.StatusNotFound),
		WithTag(NotFound),
	}
	return New(message, append(defaults, options...)...)
}

func NewInternalServerError(message string, options ...func(*Fault)) *Fault {
	defaults := []func(*Fault){
		WithHTTPCode(http.StatusInternalServerError),
		WithTag(InternalServerError),
	}
	return New(message, append(defaults, options...)...)
}

func NewUnauthorized(message string, options ...func(*Fault)) *Fault {
	defaults := []func(*Fault){
		WithHTTPCode(http.StatusUnauthorized),
		WithTag(Unauthorized),
	}
	return New(message, append(defaults, options...)...)
}

func NewForbidden(message string, options ...func(*Fault)) *Fault {
	defaults := []func(*Fault){
		WithHTTPCode(http.StatusForbidden),
		WithTag(Forbidden),
	}
	return New(message, append(defaults, options...)...)
}

func NewConflict(message string, options ...func(*Fault)) *Fault {
	defaults := []func(*Fault){
		WithHTTPCode(http.StatusConflict),
		WithTag(Conflict),
	}
	return New(message, append(defaults, options...)...)
}

func NewTooManyRequests(message string, options ...func(*Fault)) *Fault {
	defaults := []func(*Fault){
		WithHTTPCode(http.StatusTooManyRequests),
		WithTag(TooManyRequests),
	}
	return New(message, append(defaults, options...)...)
}

func NewUnprocessableEntity(message string, options ...func(*Fault)) *Fault {
	defaults := []func(*Fault){
		WithHTTPCode(http.StatusUnprocessableEntity),
		WithTag(UnprocessableEntity),
	}
	return New(message, append(defaults, options...)...)
}
