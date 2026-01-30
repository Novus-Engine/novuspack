package testhelpers

import (
	"context"
	"time"
)

// CancelledContext returns a cancelled context for testing.
// Useful for testing context cancellation handling.
//
// Example usage:
//
//	ctx := testhelpers.CancelledContext()
//	err := someFunc(ctx)
//	// Should return context.Canceled error
func CancelledContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

// TimeoutContext returns a context that has already timed out.
// Useful for testing timeout handling.
//
// Example usage:
//
//	ctx := testhelpers.TimeoutContext()
//	err := someFunc(ctx)
//	// Should return context.DeadlineExceeded error
func TimeoutContext() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), -1*time.Second)
	defer cancel()
	time.Sleep(10 * time.Millisecond)
	return ctx
}
