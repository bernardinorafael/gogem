package function

// Map applies a function to each element of a slice and returns a new slice with the results.
//
//   - iterable: slice of elements of type T
//   - callback: function that takes an element of type T and returns an element of type U
//
// Example:
//
//	names := []string{"alice", "bob", "charlie"}
//	uppercase := F.Map(names, func(name string) string {
//		return strings.ToUpper(name)
//	})
//	// uppercase = []string{"ALICE", "BOB", "CHARLIE"}
func Map[T, R any](iterable []T, callback func(T) R) []R {
	if iterable == nil {
		return nil
	}
	output := make([]R, len(iterable))
	for i, item := range iterable {
		output[i] = callback(item)
	}
	return output
}

// ForEach executes a function for each element of a slice.
//
//   - iterable: slice of elements of type T
//   - callback: function that takes an element of type T (no return value)
//
// Example:
//
//	users := []User{{Name: "Alice"}, {Name: "Bob"}}
//	F.ForEach(users, func(user User) {
//		fmt.Printf("Processing user: %s\n", user.Name)
//	})
func ForEach[T any](iterable []T, callback func(T)) {
	for _, item := range iterable {
		callback(item)
	}
}
