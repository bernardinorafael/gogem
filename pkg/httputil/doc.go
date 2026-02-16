// Package httputil provides utilities for HTTP request parsing, response writing,
// and a generic validation middleware decorator for Go web handlers.
//
// Request body parsing with size limits and detailed error messages:
//
//	var body CreateUserRequest
//	if err := httputil.ReadRequestBody(r, &body); err != nil {
//	    httputil.WriteError(w, err)
//	    return
//	}
//
// Query parameter readers with caller-defined defaults:
//
//	page := httputil.ReadQueryInt(r.URL.Query(), "page", 1)
//	sort := httputil.ReadQueryString(r.URL.Query(), "sort", "created_at")
//	active := httputil.ReadQueryBool(r.URL.Query(), "active", false)
//	tags := httputil.ReadQueryArray(r.URL.Query(), "tags")
//
// Optional variants return nil when the parameter is absent, allowing
// callers to distinguish "not provided" from "provided with zero value":
//
//	if page := httputil.ReadQueryIntOptional(r.URL.Query(), "page"); page != nil {
//	    // parameter was explicitly provided
//	}
//	if sort := httputil.ReadQueryStringOptional(r.URL.Query(), "sort"); sort != nil {
//	    // parameter was explicitly provided
//	}
//	if active := httputil.ReadQueryBoolOptional(r.URL.Query(), "active"); active != nil {
//	    // parameter was explicitly provided
//	}
//
// JSON response helpers:
//
//	_ = httputil.WriteJSON(w, http.StatusOK, user)
//	_ = httputil.WriteSuccess(w, http.StatusCreated)
//	httputil.WriteError(w, err)
//
// Generic validation middleware using WithValidation[T] and GetBody[T]:
//
//	type CreateUserDTO struct {
//	    Name  string `json:"name"`
//	    Email string `json:"email"`
//	}
//
//	func (d CreateUserDTO) Validate() error {
//	    // return validation errors
//	}
//
//	router.Post("/users", httputil.WithValidation[CreateUserDTO](func(w http.ResponseWriter, r *http.Request) {
//	    body := httputil.GetBody[CreateUserDTO](r)
//	    // body is already validated
//	}))
package httputil
