// Package uid generates K-Sortable unique identifiers with optional prefixes.
//
// Generated IDs are composed of a 4-byte timestamp (seconds since a custom epoch)
// followed by 8 bytes of cryptographic randomness, encoded as hex. This makes them
// naturally sortable by creation time while remaining globally unique.
//
// Basic usage:
//
//	id := uid.New("user")      // "user_0a1b2c3d4e5f6a7b8c9d0e1f"
//	id := uid.New("tok")       // "tok_0a1b2c3d4e5f6a7b8c9d0e1f"
//	id := uid.New("")          // "0a1b2c3d4e5f6a7b8c9d0e1f" (no prefix)
//
// Validation:
//
//	uid.IsValid("user_0a1b2c3d4e5f6a7b8c9d0e1f") // true
//	uid.IsValid("0a1b2c3d4e5f6a7b8c9d0e1f")      // true (no prefix)
//	uid.IsValid("invalid")                         // false
//
// ID format:
//
//	With prefix:    {prefix}_{24 hex chars}
//	Without prefix: {24 hex chars}
package uid
