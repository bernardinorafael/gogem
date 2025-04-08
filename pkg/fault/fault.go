package fault

import (
	"fmt"
	"net/http"
)

type Fault struct {
	HTTPCode int    `json:"-"`
	Err      error  `json:"-"`
	Tag      Tag    `json:"tag"`
	Message  string `json:"message"`
}

// New instantiates a new Fault with the given message
// The message is used to describe the error in detail
//
// The default HTTP code is 400.
func New(message string) *Fault {
	return &Fault{
		HTTPCode: http.StatusBadRequest,
		Err:      nil,
		Tag:      UNTAGGED,
		Message:  message,
	}
}

// WithHTTPCode sets the HTTP code for the fault
func (f *Fault) WithHTTPCode(code int) *Fault {
	f.HTTPCode = code
	return f
}

// WithError sets the error for the fault
func (f *Fault) WithError(err error) *Fault {
	if err == nil {
		return f
	}
	f.Err = err
	return f
}

// WithTag sets the tag for the fault
func (f *Fault) WithTag(tag Tag) *Fault {
	f.Tag = tag
	return f
}

// GetHTTPCode returns the HTTP code for the fault
func (f *Fault) GetHTTPCode() int {
	return f.HTTPCode
}

func (f *Fault) Error() string {
	if f.Err != nil {
		return fmt.Sprintf("%s:%s (caused by: %v)", f.Tag, f.Message, f.Err)
	}
	return fmt.Sprintf("%s:%s", f.Tag, f.Message)
}

func (f *Fault) Is(target error) bool {
	_, ok := target.(*Fault)
	return ok
}

func (f *Fault) Unwrap() error {
	return f.Err
}
