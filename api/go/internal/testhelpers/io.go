// Package testhelpers provides test utilities for the NovusPack v1 API.
// This package contains mock I/O types, context helpers, and test data creators
// used across all test files in the v1 API.
package testhelpers

import "errors"

// ErrorWriter is a writer that always returns an error.
// Useful for testing error handling in Write operations.
//
// Example usage:
//
//	w := testhelpers.NewErrorWriter()
//	_, err := someFunc.WriteTo(w)
//	if err == nil {
//	    t.Error("expected error")
//	}
type ErrorWriter struct {
	err error
}

// NewErrorWriter creates a new ErrorWriter.
// If no custom error is provided, it returns "write error" by default.
func NewErrorWriter(customErr ...error) *ErrorWriter {
	w := &ErrorWriter{}
	if len(customErr) > 0 {
		w.err = customErr[0]
	}
	return w
}

func (w *ErrorWriter) Write(p []byte) (int, error) {
	if w.err == nil {
		w.err = errors.New("write error")
	}
	return 0, w.err
}

// FailingWriter is a writer that fails after writing a specified number of bytes.
// Useful for testing partial write scenarios and error recovery.
//
// Example usage:
//
//	w := testhelpers.NewFailingWriter(10) // Fail after 10 bytes
//	n, err := someFunc.WriteTo(w)
//	// First 10 bytes succeed, then error
type FailingWriter struct {
	maxBytes int
	written  int
}

// NewFailingWriter creates a writer that fails after maxBytes have been written.
func NewFailingWriter(maxBytes int) *FailingWriter {
	return &FailingWriter{maxBytes: maxBytes}
}

func (w *FailingWriter) Write(p []byte) (int, error) {
	if w.written >= w.maxBytes {
		return 0, errors.New("write failed")
	}
	remaining := w.maxBytes - w.written
	if len(p) > remaining {
		w.written = w.maxBytes
		return remaining, errors.New("write failed")
	}
	w.written += len(p)
	return len(p), nil
}

// IncompleteWriter is a writer that successfully writes partial data.
// Unlike FailingWriter, it returns no error but writes fewer bytes than requested.
// Useful for testing incomplete write handling.
//
// Example usage:
//
//	w := testhelpers.NewIncompleteWriter(10)
//	n, err := w.Write([]byte("hello world")) // writes "hello worl", returns 10, nil
type IncompleteWriter struct {
	maxWrite int
	written  int
}

// NewIncompleteWriter creates a writer that writes at most maxWrite bytes total.
func NewIncompleteWriter(maxWrite int) *IncompleteWriter {
	return &IncompleteWriter{maxWrite: maxWrite}
}

func (w *IncompleteWriter) Write(p []byte) (int, error) {
	if w.written >= w.maxWrite {
		return 0, errors.New("write failed")
	}
	toWrite := len(p)
	if toWrite > w.maxWrite-w.written {
		toWrite = w.maxWrite - w.written
	}
	w.written += toWrite
	return toWrite, nil
}

// PartialWriter is a writer that writes partial data and returns 0 bytes without error.
// This triggers incomplete write detection logic (when n != len(p) but err == nil).
//
// Example usage:
//
//	w := testhelpers.NewPartialWriter(10)
//	n, err := w.Write([]byte("test")) // Returns (4, nil) until 10 bytes written
type PartialWriter struct {
	maxBytes int
	written  int
}

// NewPartialWriter creates a writer that writes up to maxBytes, then returns (0, nil).
func NewPartialWriter(maxBytes int) *PartialWriter {
	return &PartialWriter{maxBytes: maxBytes}
}

func (w *PartialWriter) Write(p []byte) (int, error) {
	if w.written >= w.maxBytes {
		return 0, nil // Return 0 bytes written, no error
	}
	n := len(p)
	if w.written+n > w.maxBytes {
		n = w.maxBytes - w.written
	}
	w.written += n
	return n, nil
}

// ErrorReader is a reader that always returns a non-EOF error.
// Useful for testing error handling in Read operations.
//
// Example usage:
//
//	r := testhelpers.NewErrorReader()
//	_, err := someFunc.ReadFrom(r)
//	if err == nil {
//	    t.Error("expected error")
//	}
type ErrorReader struct {
	err error
}

// NewErrorReader creates a new ErrorReader.
// If no custom error is provided, it returns "read error" by default.
func NewErrorReader(customErr ...error) *ErrorReader {
	r := &ErrorReader{}
	if len(customErr) > 0 {
		r.err = customErr[0]
	}
	return r
}

func (r *ErrorReader) Read(p []byte) (int, error) {
	if r.err == nil {
		r.err = errors.New("read error")
	}
	return 0, r.err
}

// PartialReader reads some bytes then returns an error.
// Useful for testing partial read scenarios where some data is read successfully
// before an error occurs.
//
// Example usage:
//
//	r := testhelpers.NewPartialReader([]byte("test"), errors.New("read failed"))
//	buf := make([]byte, 10)
//	n, err := r.Read(buf) // Reads "test", then next call returns error
type PartialReader struct {
	data []byte
	pos  int
	err  error
}

// NewPartialReader creates a reader that reads data then returns an error.
// If no error is provided, it returns "read error" after data is exhausted.
func NewPartialReader(data []byte, err ...error) *PartialReader {
	r := &PartialReader{data: data}
	if len(err) > 0 {
		r.err = err[0]
	}
	return r
}

func (r *PartialReader) Read(p []byte) (int, error) {
	if r.pos >= len(r.data) {
		if r.err == nil {
			r.err = errors.New("read error")
		}
		return 0, r.err
	}
	n := copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}
