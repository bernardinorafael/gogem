// Package fault provides a standardized REST error type implementing Go's error
// interface with HTTP status codes, semantic tags, and field-level validation errors.
//
// Fault uses the functional options pattern for flexible error construction.
// Each error carries an HTTP code, a human-readable message, a semantic tag
// for programmatic handling, and optionally field-level validation details.
//
// Basic usage:
//
//	err := fault.NewNotFound("user not found")
//	err := fault.NewBadRequest("invalid email format")
//
// Custom errors with options:
//
//	err := fault.New("payment failed",
//	    fault.WithHTTPCode(http.StatusConflict),
//	    fault.WithTag(fault.Conflict),
//	    fault.WithErr(originalErr),
//	)
//
// Validation errors with field details:
//
//	err := fault.NewValidation("invalid body", validationErr)
//	// Produces: { "status": 422, "message": "invalid body", "fields": [{"field": "email", "message": "required"}] }
//
// Extracting tags from wrapped errors:
//
//	switch fault.GetTag(err) {
//	case fault.NotFound:
//	    // handle not found
//	case fault.BadRequest:
//	    // handle bad request
//	}
//
// Fault implements Is(), Unwrap(), and Error() for seamless integration
// with Go's errors package and error chains.
package fault
