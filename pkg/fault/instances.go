package fault

import (
	"encoding/json"
	"net/http"
)

// NewHTTPError receives an error and writes it to the response writer
// It sets the content type to application/json and writes the error
// If the error is not a Fault, it writes a new InternalServerError
func NewHTTPError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	if err, ok := err.(*Fault); ok {
		w.WriteHeader(err.GetHTTPCode())
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(NewInternalServerError("an unexpected error occurred"))
}

func NewBadRequest(msg string) *Fault {
	return &Fault{
		HTTPCode: http.StatusBadRequest,
		Err:      nil,
		Tag:      BAD_REQUEST,
		Message:  msg,
	}
}

func NewNotFound(msg string) *Fault {
	return &Fault{
		HTTPCode: http.StatusNotFound,
		Err:      nil,
		Tag:      NOT_FOUND,
		Message:  msg,
	}
}

func NewInternalServerError(msg string) *Fault {
	return &Fault{
		HTTPCode: http.StatusInternalServerError,
		Err:      nil,
		Tag:      INTERNAL_SERVER_ERROR,
		Message:  msg,
	}
}

func NewUnauthorized(msg string) *Fault {
	return &Fault{
		HTTPCode: http.StatusUnauthorized,
		Err:      nil,
		Tag:      UNAUTHORIZED,
		Message:  msg,
	}
}

func NewForbidden(msg string) *Fault {
	return &Fault{
		HTTPCode: http.StatusForbidden,
		Err:      nil,
		Tag:      FORBIDDEN,
		Message:  msg,
	}
}

func NewConflict(msg string) *Fault {
	return &Fault{
		HTTPCode: http.StatusConflict,
		Err:      nil,
		Tag:      CONFLICT,
		Message:  msg,
	}
}

func NewTooManyRequests(msg string) *Fault {
	return &Fault{
		HTTPCode: http.StatusTooManyRequests,
		Err:      nil,
		Tag:      TOO_MANY_REQUESTS,
		Message:  msg,
	}
}
