// Package apiutil provides generic utilities for API response construction.
//
// The Expandable[T] type represents a field that can be serialized as either
// a plain ID string or a fully expanded nested object, depending on whether
// the data was loaded. This is useful for API responses where related resources
// can optionally be expanded inline.
//
// Collapsed (data not loaded):
//
//	field := apiutil.NewExpandableField[User]("user_0a1b2c...", nil)
//	json.Marshal(field)
//	// "user_0a1b2c..."
//
// Expanded (data loaded):
//
//	field := apiutil.NewExpandableField("user_0a1b2c...", &user)
//	json.Marshal(field)
//	// {"id": "user_0a1b2c...", "name": "Alice", "email": "alice@example.com"}
//
// Typical usage in a response struct:
//
//	type OrderResponse struct {
//	    ID       string                       `json:"id"`
//	    Total    int                          `json:"total"`
//	    Customer *apiutil.Expandable[Customer] `json:"customer"`
//	}
//
//	// Without expand: { "id": "...", "total": 100, "customer": "cust_abc123" }
//	// With expand:    { "id": "...", "total": 100, "customer": { "name": "Alice", ... } }
package apiutil
