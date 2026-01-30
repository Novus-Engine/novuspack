package testhelpers

import "strings"

// Contains checks if s contains substr (case-sensitive).
// This is a convenience wrapper around strings.Contains from the standard library.
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// ContainsIgnoreCase checks if s contains substr (case-insensitive).
func ContainsIgnoreCase(s, substr string) bool {
	return strings.Contains(
		strings.ToLower(s),
		strings.ToLower(substr),
	)
}

// IndexIgnoreCase finds the index of substr in s (case-insensitive).
// Returns -1 if not found.
func IndexIgnoreCase(s, substr string) int {
	sLower := strings.ToLower(s)
	substrLower := strings.ToLower(substr)
	return strings.Index(sLower, substrLower)
}
