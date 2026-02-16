package fault

import (
	"errors"
	"fmt"
	"net/http"
)

type FieldError struct {
	Field   string `json:"field" example:"email"`
	Message string `json:"message" example:"invalid email format"`
}

func NewFieldError(field, message string) FieldError {
	return FieldError{Field: field, Message: message}
}

type Fault struct {
	HTTPCode   int          `json:"status" example:"400"`
	Message    string       `json:"message" example:"validation failed"`
	FieldError []FieldError `json:"fields"`

	Tag Tag   `json:"-"`
	Err error `json:"-"`
}

// New instantiates a new Fault with the given message
// The message is used to describe the error in detail
//
// The default HTTP code is 400.
func New(msg string, options ...func(*Fault)) *Fault {
	var validations = make([]FieldError, 0)

	fault := Fault{
		Err:        nil,
		Tag:        Untagged,
		HTTPCode:   http.StatusBadRequest,
		Message:    msg,
		FieldError: validations,
	}

	for _, fn := range options {
		fn(&fault)
	}

	return &fault
}

// WithHTTPCode sets the HTTP code for the fault
func WithHTTPCode(code int) func(*Fault) {
	return func(f *Fault) {
		f.HTTPCode = code
	}
}

// WithErr WithError sets the error for the fault
func WithErr(err error) func(*Fault) {
	return func(f *Fault) {
		if err == nil {
			return
		}
		f.Err = err
	}
}

// WithTag sets the tag for the fault
func WithTag(tag Tag) func(*Fault) {
	return func(f *Fault) {
		f.Tag = tag
	}
}

// WithFieldError sets the field errors for the fault
func WithFieldError(args ...FieldError) func(*Fault) {
	return func(f *Fault) {
		f.FieldError = args
	}
}

// GetHTTPCode returns the HTTP code for the fault
func (f *Fault) GetHTTPCode() int {
	return f.HTTPCode
}

func (f *Fault) Error() string {
	if f.Err != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", f.Tag, f.Message, f.Err)
	}
	return fmt.Sprintf("%s: %s", f.Tag, f.Message)
}

func (f *Fault) Is(target error) bool {
	var t *Fault
	if !errors.As(target, &t) {
		return false
	}
	return f.Tag == t.Tag
}

func (f *Fault) Unwrap() error {
	return f.Err
}
