//go:build bdd

// Package contextkeys provides context keys for BDD test infrastructure.
// This package exists to avoid import cycles between support and steps packages.
package contextkeys

// WorldContextKeyType is a custom type for the context key to avoid collisions.
// This follows Go best practices for context keys (SA1029).
type WorldContextKeyType string

// WorldContextKey is the context key used to store and retrieve the World object.
// It must be used consistently across all packages that access the world from context.
const WorldContextKey WorldContextKeyType = "world"
