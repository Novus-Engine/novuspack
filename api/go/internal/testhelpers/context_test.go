package testhelpers

import (
	"context"
	"testing"
	"time"
)

func TestCancelledContext(t *testing.T) {
	t.Run("returns cancelled context", func(t *testing.T) {
		ctx := CancelledContext()
		if ctx == nil {
			t.Fatal("CancelledContext should not return nil")
		}

		// Context should be cancelled
		select {
		case <-ctx.Done():
			// Expected - context is done
		default:
			t.Error("CancelledContext should return a done context")
		}

		// Should return context.Canceled error
		err := ctx.Err()
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	})

	t.Run("context error in operations", func(t *testing.T) {
		ctx := CancelledContext()

		// Simulate using the context in an operation
		select {
		case <-ctx.Done():
			err := ctx.Err()
			if err != context.Canceled {
				t.Errorf("expected context.Canceled error, got %v", err)
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("context should be immediately cancelled")
		}
	})
}

func TestTimeoutContext(t *testing.T) {
	t.Run("returns timed out context", func(t *testing.T) {
		ctx := TimeoutContext()
		if ctx == nil {
			t.Fatal("TimeoutContext should not return nil")
		}

		// Context should be done (timed out)
		select {
		case <-ctx.Done():
			// Expected - context is done
		case <-time.After(100 * time.Millisecond):
			t.Error("TimeoutContext should return a done context")
		}

		// Should return context.DeadlineExceeded error
		err := ctx.Err()
		if err != context.DeadlineExceeded {
			t.Errorf("expected context.DeadlineExceeded, got %v", err)
		}
	})

	t.Run("context deadline exceeded in operations", func(t *testing.T) {
		ctx := TimeoutContext()

		// Simulate using the context in an operation
		select {
		case <-ctx.Done():
			err := ctx.Err()
			if err != context.DeadlineExceeded {
				t.Errorf("expected context.DeadlineExceeded error, got %v", err)
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("context should be timed out")
		}
	})

	t.Run("deadline has passed", func(t *testing.T) {
		ctx := TimeoutContext()

		// Check that deadline has already passed
		deadline, ok := ctx.Deadline()
		if !ok {
			t.Error("TimeoutContext should have a deadline")
		}

		if time.Now().Before(deadline) {
			t.Error("deadline should be in the past")
		}
	})
}
