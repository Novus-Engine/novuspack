package fileformat

import (
	"errors"
)

// failingWriter is a writer that fails after writing a specified number of bytes
type failingWriter struct {
	maxBytes int
	written  int
}

func (w *failingWriter) Write(p []byte) (int, error) {
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

// errorWriter is a writer that always returns an error
type errorWriter struct {
	err error
}

func (w *errorWriter) Write(p []byte) (int, error) {
	if w.err == nil {
		w.err = errors.New("write error")
	}
	return 0, w.err
}

// incompleteWriter is a writer that writes fewer bytes than requested
type incompleteWriter struct {
	maxWrite int
	written  int
}

func (w *incompleteWriter) Write(p []byte) (int, error) {
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
