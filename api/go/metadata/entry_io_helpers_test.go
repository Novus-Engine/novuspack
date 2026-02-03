// Shared table runners for HashEntry and OptionalDataEntry WriteTo/ReadFrom tests.

package metadata

import (
	"bytes"
	"io"
	"testing"
)

type writeToEntry interface {
	writeTo(io.Writer) (int64, error)
}

type readFromEntry interface {
	readFrom(io.Reader) (int64, error)
}

type writeToCase struct {
	name    string
	entry   writeToEntry
	wantErr bool
	minSize int64
}

func runWriteToEntryTable(t *testing.T, tests []writeToCase) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			n, err := tt.entry.writeTo(&buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if n == 0 {
				t.Error("WriteTo() wrote 0 bytes")
				return
			}
			if tt.minSize > 0 && n < tt.minSize {
				t.Errorf("WriteTo() wrote %d bytes, want at least %d", n, tt.minSize)
			}
		})
	}
}

type writeToErrorCase struct {
	name    string
	writer  io.Writer
	wantErr bool
}

func runWriteToErrorPathsTable(t *testing.T, entry writeToEntry, tests []writeToErrorCase) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := entry.writeTo(tt.writer)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err == nil {
				t.Error("WriteTo() expected error but got nil")
			}
		})
	}
}

type readFromIncompleteCase struct {
	name string
	data []byte
}

func runReadFromIncompleteTable(t *testing.T, tests []readFromIncompleteCase, newEntry func() readFromEntry) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := newEntry()
			_, err := entry.readFrom(bytes.NewReader(tt.data))
			if err == nil {
				t.Errorf("ReadFrom() expected error for incomplete data, got nil")
			}
		})
	}
}
