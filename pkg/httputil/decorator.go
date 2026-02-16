package httputil

import (
	"context"
	"net/http"

	"github.com/bernardinorafael/gogem/pkg/fault"
)

// Validator is an interface that defines a contract for validating data structures (DTOs).
// Any DTO representing the body of a request must implement this Validate method,
// returning an error if any validation rule is not met.
// This allows middlewares and handlers to generically ensure the integrity of received data.
//
// Example:
//
//	type dto.UserRequest struct {
//		Name string `json:"name"`
//		Email string `json:"email"`
//	}
//
//	func (u dto.UserRequest) Validate() error {
//		if u.Name == "" {
//			return errors.New("name is required")
//		}
//		return nil
//	}
type Validator interface {
	Validate() error
}

// WithValidation is a decorator that validates the request body before passing control to the next handler.
// It performs the following operations:
//  1. Reads and deserializes the request body into the specified DTO type T
//  2. Validates the DTO using its Validate() method
//  3. If validation passes, stores the validated DTO in the request context
//  4. If validation fails, returns an appropriate error response
//  5. Calls the next handler with the enriched request context
//
// Usage:
//
//	func createUserHandler(w http.ResponseWriter, r *http.Request) {
//		userData := GetBodyFromContext[dto.UserRequest](r)
//		// Process the validated user data...
//	}
//
//	// In your router setup:
//	router.HandleFunc("/users", WithValidation[dto.UserRequest](createUserHandler))
//
// The decorator ensures that only valid data reaches your handlers, reducing boilerplate
// validation code and improving code maintainability.
func WithValidation[T Validator](done http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body T
		if err := ReadRequestBody(w, r, &body); err != nil {
			WriteError(w, err)
			return
		}

		if err := body.Validate(); err != nil {
			WriteError(w, fault.NewValidation("invalid body", err))
			return
		}

		ctx := r.Context()
		// ctx = context.WithValue(ctx, RequestKey{}, body)
		ctx = context.WithValue(ctx, "body", body)
		r = r.WithContext(ctx)

		done(w, r)
	}
}
