package fault

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type FieldError struct {
	Field   string `json:"field" example:"email"`
	Message string `json:"message" example:"invalid email format"`
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

func WithValidationError(err error) func(*Fault) {
	if err == nil {
		return func(_ *Fault) {}
	}

	var validations []FieldError
	split := strings.SplitSeq(err.Error(), ";")

	for validation := range split {
		validation = strings.TrimSpace(validation)
		if validation == "" {
			continue
		}

		parts := strings.SplitN(validation, ":", 2)
		if len(parts) != 2 {
			validations = append(validations, FieldError{
				Field:   "general",
				Message: validation,
			})
			continue
		}

		field := strings.TrimSpace(parts[0])
		msg := strings.TrimSpace(parts[1])

		// If the field or message is empty, skip it
		if field == "" || msg == "" {
			continue
		}

		removePeriod := func(s string) string {
			if strings.HasSuffix(s, ".") {
				return s[:len(s)-1]
			}
			return s
		}

		validations = append(validations, FieldError{
			Field:   field,
			Message: removePeriod(msg),
		})
	}

	return func(f *Fault) {
		f.FieldError = validations
	}
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
	var fault *Fault
	return errors.As(target, &fault)
}

func (f *Fault) Unwrap() error {
	return f.Err
}
