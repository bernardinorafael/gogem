// Package pagination provides a generic container for paginated API responses
// with automatically computed metadata such as total pages, navigation flags,
// and current page position.
//
// Basic usage:
//
//	users := []User{...} // items for the current page
//	result := pagination.New(users, totalItems, currentPage, itemsPerPage)
//
//	// result.Items           → []User
//	// result.Pagination      → Pagination metadata
//
// The Pagination metadata includes:
//
//	result.Pagination.TotalItems      // total number of items across all pages
//	result.Pagination.TotalPages      // computed total page count
//	result.Pagination.CurrentPage     // current page number
//	result.Pagination.ItemsPerPage    // items per page
//	result.Pagination.HasNextPage     // true if there are more pages ahead
//	result.Pagination.HasPreviousPage // true if there are pages before
//	result.Pagination.IsFirstPage     // true if current page is the first
//	result.Pagination.IsLastPage      // true if current page is the last
//
// The Paginated[T] type serializes directly to JSON:
//
//	httputil.WriteJSON(w, http.StatusOK, result)
//	// { "items": [...], "pagination": { "totalItems": 100, "currentPage": 1, ... } }
package pagination
