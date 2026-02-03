package testhelpers

import (
	"errors"
	"io"
	"testing"
)

// assertReadResult checks n and err from a Read call; fails t if err != nil or n != wantN.
func assertReadResult(t *testing.T, n int, err error, wantN int, op string) {
	t.Helper()
	if err != nil {
		t.Errorf("%s should succeed, got error: %v", op, err)
	}
	if n != wantN {
		t.Errorf("%s: expected %d bytes, got %d", op, wantN, n)
	}
}

// testWriterTwoWrites performs two writes and checks n and err for each; used to reduce dupl in tests.
func testWriterTwoWrites(t *testing.T, w io.Writer, first, second []byte, wantFirstN, wantSecondN int) {
	t.Helper()
	n, err := w.Write(first)
	if err != nil {
		t.Errorf("write should succeed, got error: %v", err)
	}
	if n != wantFirstN {
		t.Errorf("expected %d bytes written, got %d", wantFirstN, n)
	}
	n, err = w.Write(second)
	if err != nil {
		t.Errorf("second write should succeed, got error: %v", err)
	}
	if n != wantSecondN {
		t.Errorf("expected %d bytes written, got %d", wantSecondN, n)
	}
}

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
		testFailingWriterAtLimit(t, 5, 1)
	})

	t.Run("fails immediately when at limit", func(t *testing.T) {
		testFailingWriterAtLimit(t, 4, 0)
	})
}

// testFailingWriterAtLimit writes "test" then "more" to a FailingWriter with the given limit; expects write to fail and n to match expectedN.
func testFailingWriterAtLimit(t *testing.T, limit, expectedN int) {
	t.Helper()
	w := NewFailingWriter(limit)
	_, _ = w.Write([]byte("test"))
	n, err := w.Write([]byte("more"))
	if err == nil {
		t.Error("second write should fail")
	}
	if n != expectedN {
		t.Errorf("expected %d byte(s) written, got %d", expectedN, n)
	}
}

func TestIncompleteWriter(t *testing.T) {
	t.Run("writes partial data", func(t *testing.T) {
		testWriterTwoWrites(t, NewIncompleteWriter(10), []byte("hello"), []byte("world"), 5, 5)
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
		testWriterTwoWrites(t, NewPartialWriter(10), []byte("test"), []byte("hello world"), 4, 6)
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
		buf := make([]byte, 5)
		n, err := r.Read(buf)
		assertReadResult(t, n, err, 5, "first read")
		if string(buf[:n]) != "test " {
			t.Errorf("expected 'test ', got %q", string(buf[:n]))
		}
		n, err = r.Read(buf)
		assertReadResult(t, n, err, 4, "second read")
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
