// Package function provides generic utility functions for slice operations.
//
// This package is typically imported with the alias F for concise usage:
//
//	import F "github.com/bernardinorafael/gogem/function"
//
// Map transforms each element of a slice using a callback:
//
//	names := F.Map(users, func(u User) string {
//	    return u.Name
//	})
//	// []string{"Alice", "Bob", "Charlie"}
//
//	doubled := F.Map([]int{1, 2, 3}, func(n int) int {
//	    return n * 2
//	})
//	// []int{2, 4, 6}
//
// ForEach executes a side-effect function for each element:
//
//	F.ForEach(users, func(u User) {
//	    fmt.Printf("Processing: %s\n", u.Name)
//	})
package function
