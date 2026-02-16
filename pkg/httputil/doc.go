// Package httputil provides utilities for HTTP request parsing, response writing,
// and a generic validation middleware decorator for Go web handlers.
//
// Request body parsing with size limits and detailed error messages:
//
//	var body CreateUserRequest
//	if err := httputil.ReadRequestBody(w, r, &body); err != nil {
//	    httputil.WriteError(w, err)
//	    return
//	}
//
// Query parameter readers with zero-value defaults:
//
//	page := httputil.ReadQueryInt(r.URL.Query(), "page")           // 0 if missing
//	sort := httputil.ReadQueryString(r.URL.Query(), "sort")        // "" if missing
//	active := httputil.ReadQueryBool(r.URL.Query(), "active")      // false if missing
//	tags := httputil.ReadQueryArray(r.URL.Query(), "tags")         // [] if missing
//
// JSON response helpers:
//
//	httputil.WriteJSON(w, http.StatusOK, user)
//	httputil.WriteSuccess(w, http.StatusCreated)
//	httputil.WriteError(w, err) // automatically extracts HTTP code from fault.Fault
//
// Generic validation middleware using WithValidation[T]:
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
//	    body := r.Context().Value("body").(CreateUserDTO)
//	    // body is already validated
//	}))
package httputil
