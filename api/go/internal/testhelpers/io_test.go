package testhelpers

import (
	"errors"
	"testing"
)

func TestErrorWriter(t *testing.T) {
	t.Run("default error", func(t *testing.T) {
		w := NewErrorWriter()
		_, err := w.Write([]byte("test"))
		if err == nil {
			t.Error("ErrorWriter should return error")
		}
		if err.Error() != "write error" {
			t.Errorf("expected 'write error', got %q", err.Error())
		}
	})

	t.Run("custom error", func(t *testing.T) {
		customErr := errors.New("custom error")
		w := NewErrorWriter(customErr)
		_, err := w.Write([]byte("test"))
		if err != customErr {
			t.Errorf("expected custom error, got %v", err)
		}
	})
}

func TestFailingWriter(t *testing.T) {
	t.Run("writes until limit", func(t *testing.T) {
		w := NewFailingWriter(5)
		n, err := w.Write([]byte("test"))
		if err != nil {
			t.Errorf("first write should succeed, got error: %v", err)
		}
		if n != 4 {
			t.Errorf("expected 4 bytes written, got %d", n)
		}
	})

	t.Run("fails after limit", func(t *testing.T) {
		w := NewFailingWriter(5)
		_, _ = w.Write([]byte("test"))
		n, err := w.Write([]byte("more"))
		if err == nil {
			t.Error("second write should fail")
		}
		// Should write 1 byte (5-4=1) then fail
		if n != 1 {
			t.Errorf("expected 1 byte written, got %d", n)
		}
	})

	t.Run("fails immediately when at limit", func(t *testing.T) {
		w := NewFailingWriter(4)
		_, _ = w.Write([]byte("test"))
		n, err := w.Write([]byte("more"))
		if err == nil {
			t.Error("write at limit should fail")
		}
		if n != 0 {
			t.Errorf("expected 0 bytes written, got %d", n)
		}
	})
}

func TestIncompleteWriter(t *testing.T) {
	t.Run("writes partial data", func(t *testing.T) {
		w := NewIncompleteWriter(10)

		n, err := w.Write([]byte("hello"))
		if err != nil {
			t.Errorf("write should succeed, got error: %v", err)
		}
		if n != 5 {
			t.Errorf("expected 5 bytes written, got %d", n)
		}

		n, err = w.Write([]byte("world"))
		if err != nil {
			t.Errorf("partial write should succeed, got error: %v", err)
		}
		if n != 5 {
			t.Errorf("expected 5 bytes written (reaching limit), got %d", n)
		}
	})

	t.Run("fails beyond limit", func(t *testing.T) {
		w := NewIncompleteWriter(10)
		_, _ = w.Write([]byte("hello"))
		_, _ = w.Write([]byte("world"))

		_, err := w.Write([]byte("more"))
		if err == nil {
			t.Error("write beyond limit should fail")
		}
	})
}

func TestPartialWriter(t *testing.T) {
	t.Run("writes until limit", func(t *testing.T) {
		w := NewPartialWriter(10)

		// First write succeeds
		n, err := w.Write([]byte("test"))
		if err != nil {
			t.Errorf("write should succeed, got error: %v", err)
		}
		if n != 4 {
			t.Errorf("expected 4 bytes, got %d", n)
		}

		// Write up to limit
		n, err = w.Write([]byte("hello world"))
		if err != nil {
			t.Errorf("write should succeed, got error: %v", err)
		}
		if n != 6 {
			t.Errorf("expected 6 bytes (to reach 10 total), got %d", n)
		}
	})

	t.Run("returns zero without error beyond limit", func(t *testing.T) {
		w := NewPartialWriter(10)
		_, _ = w.Write([]byte("test"))
		_, _ = w.Write([]byte("hello world"))

		// Beyond limit returns (0, nil)
		n, err := w.Write([]byte("more"))
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if n != 0 {
			t.Errorf("expected 0 bytes written, got %d", n)
		}
	})
}

func TestErrorReader(t *testing.T) {
	t.Run("default error", func(t *testing.T) {
		r := NewErrorReader()
		buf := make([]byte, 10)
		_, err := r.Read(buf)
		if err == nil {
			t.Error("ErrorReader should return error")
		}
		if err.Error() != "read error" {
			t.Errorf("expected 'read error', got %q", err.Error())
		}
	})

	t.Run("custom error", func(t *testing.T) {
		customErr := errors.New("custom error")
		r := NewErrorReader(customErr)
		buf := make([]byte, 10)
		_, err := r.Read(buf)
		if err != customErr {
			t.Errorf("expected custom error, got %v", err)
		}
	})
}

func TestPartialReader(t *testing.T) {
	t.Run("reads data then errors", func(t *testing.T) {
		data := []byte("test data")
		r := NewPartialReader(data)

		// First read succeeds
		buf := make([]byte, 5)
		n, err := r.Read(buf)
		if err != nil {
			t.Errorf("first read should succeed, got error: %v", err)
		}
		if n != 5 {
			t.Errorf("expected 5 bytes, got %d", n)
		}
		if string(buf) != "test " {
			t.Errorf("expected 'test ', got %q", string(buf))
		}

		// Second read gets remaining data
		n, err = r.Read(buf)
		if err != nil {
			t.Errorf("second read should succeed, got error: %v", err)
		}
		if n != 4 {
			t.Errorf("expected 4 bytes, got %d", n)
		}

		// Third read returns error
		n, err = r.Read(buf)
		if err == nil {
			t.Error("read beyond data should return error")
		}
		if n != 0 {
			t.Errorf("expected 0 bytes, got %d", n)
		}
	})

	t.Run("custom error", func(t *testing.T) {
		customErr := errors.New("custom read error")
		r := NewPartialReader([]byte("test"), customErr)

		// Read all data
		buf := make([]byte, 10)
		_, _ = r.Read(buf)

		// Next read returns custom error
		_, err := r.Read(buf)
		if err != customErr {
			t.Errorf("expected custom error, got %v", err)
		}
	})
}
