package generics

import (
	"bytes"
	"encoding/binary"
	"io"
	"strings"
	"testing"

	"github.com/novus-engine/novuspack/api/go/internal/testhelpers"
)

// TestPathEntrySize verifies PathEntry size calculation
func TestPathEntrySize(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wantSize int
	}{
		{"Short path", "test.txt", 10},                                           // 2 (PathLength) + 8 (Path)
		{"Medium path", "path/to/file.txt", 18},                                  // 2 (PathLength) + 16 (Path)
		{"Long path", "very/long/path/to/some/file/in/a/deep/directory.txt", 53}, // 2 (PathLength) + 51 (Path)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := PathEntry{
				PathLength: uint16(len(tt.path)),
				Path:       tt.path,
			}

			if entry.Size() != tt.wantSize {
				t.Errorf("Size() = %d, want %d", entry.Size(), tt.wantSize)
			}
		})
	}
}

// TestPathEntryValidation verifies validation logic
func TestPathEntryValidation(t *testing.T) {
	tests := []struct {
		name    string
		entry   PathEntry
		wantErr bool
	}{
		{
			"Valid path",
			PathEntry{
				PathLength: 8,
				Path:       "test.txt",
			},
			false,
		},
		{
			"Empty path",
			PathEntry{
				PathLength: 0,
				Path:       "",
			},
			true,
		},
		{
			"Whitespace-only path",
			PathEntry{
				PathLength: 3,
				Path:       "   ",
			},
			true,
		},
		{
			"Tab-only path",
			PathEntry{
				PathLength: 2,
				Path:       "\t\t",
			},
			true,
		},
		{
			"Newline-only path",
			PathEntry{
				PathLength: 1,
				Path:       "\n",
			},
			true,
		},
		{
			"Length mismatch",
			PathEntry{
				PathLength: 10,
				Path:       "test.txt",
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.entry.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestPathEntryReadFrom verifies ReadFrom deserialization
// Specification: package_file_format.md: 4.1.4.2 Path Entries
func TestPathEntryReadFrom(t *testing.T) {
	tests := []struct {
		name    string
		entry   PathEntry
		wantErr bool
	}{
		{
			"Valid path entry",
			PathEntry{
				PathLength: 8,
				Path:       "test.txt",
			},
			false,
		},
		{
			"Path entry with long path",
			PathEntry{
				PathLength: 20,
				Path:       "path/to/test/file.txt",
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure PathLength matches actual path length
			tt.entry.PathLength = uint16(len(tt.entry.Path))

			// Serialize original entry using WriteTo
			buf := new(bytes.Buffer)
			_, err := tt.entry.WriteTo(buf)
			if err != nil {
				t.Fatalf("Failed to serialize entry: %v", err)
			}

			// Deserialize using ReadFrom
			var entry PathEntry
			n, err := entry.ReadFrom(buf)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// After updating PathLength, the expected size should match what was written
				// Calculate expected size: 2 (PathLength) + PathLength (path bytes)
				expectedSize := 2 + int(tt.entry.PathLength)
				if n != int64(expectedSize) {
					t.Errorf("ReadFrom() read %d bytes, want %d", n, expectedSize)
				}

				// Verify all fields match
				if entry.PathLength != tt.entry.PathLength {
					t.Errorf("PathLength = %d, want %d", entry.PathLength, tt.entry.PathLength)
				}
				if entry.Path != tt.entry.Path {
					t.Errorf("Path = %q, want %q", entry.Path, tt.entry.Path)
				}

				// Verify validation passes
				if err := entry.Validate(); err != nil {
					t.Errorf("ReadFrom() entry validation failed: %v", err)
				}
			}
		})
	}
}

// TestPathEntryReadFromIncompleteData verifies ReadFrom handles incomplete data
func TestPathEntryReadFromIncompleteData(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"No data", []byte{}},
		{"Incomplete PathLength", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint16(8))
			// Missing second byte of PathLength and Path
			return buf.Bytes()[:1] // Only 1 byte of PathLength
		}()},
		{"Incomplete Path", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint16(8))
			buf.WriteString("test") // Only 4 bytes of 8-byte path
			return buf.Bytes()
		}()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var entry PathEntry
			r := bytes.NewReader(tt.data)
			_, err := entry.ReadFrom(r)

			if err == nil {
				t.Errorf("ReadFrom() expected error for incomplete data, got nil")
			}
		})
	}
}

// TestPathEntryReadFromEmptyPath verifies ReadFrom handles empty path correctly
func TestPathEntryReadFromEmptyPath(t *testing.T) {
	entry := PathEntry{
		PathLength: 0,
		Path:       "",
	}

	var buf bytes.Buffer
	_, err := entry.WriteTo(&buf)
	if err != nil {
		t.Fatalf("Failed to serialize entry: %v", err)
	}

	var readEntry PathEntry
	_, err = readEntry.ReadFrom(&buf)
	if err != nil {
		t.Fatalf("ReadFrom() error = %v", err)
	}

	if readEntry.PathLength != 0 {
		t.Errorf("PathLength = %d, want 0", readEntry.PathLength)
	}
	if readEntry.Path != "" {
		t.Errorf("Path = %q, want empty", readEntry.Path)
	}
}

// TestPathEntryWriteTo verifies WriteTo serialization
// Specification: package_file_format.md: 4.1.4.2 Path Entries
func TestPathEntryWriteTo(t *testing.T) {
	tests := []struct {
		name    string
		entry   PathEntry
		wantErr bool
	}{
		{
			"Valid path entry",
			PathEntry{
				PathLength: 8,
				Path:       "test.txt",
			},
			false,
		},
		{
			"Path entry with long path",
			PathEntry{
				PathLength: 20,
				Path:       "path/to/test/file.txt",
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure PathLength matches actual path length
			tt.entry.PathLength = uint16(len(tt.entry.Path))

			var buf bytes.Buffer
			n, err := tt.entry.WriteTo(&buf)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				expectedSize := tt.entry.Size()
				if n != int64(expectedSize) {
					t.Errorf("WriteTo() wrote %d bytes, want %d", n, expectedSize)
				}

				if buf.Len() != expectedSize {
					t.Errorf("WriteTo() buffer size = %d bytes, want %d", buf.Len(), expectedSize)
				}

				// Verify we can read it back
				var entry PathEntry
				_, readErr := entry.ReadFrom(&buf)
				if readErr != nil {
					t.Errorf("Failed to read back written data: %v", readErr)
				}

				if entry.PathLength != tt.entry.PathLength {
					t.Errorf("PathLength mismatch: %d != %d", entry.PathLength, tt.entry.PathLength)
				}
				if entry.Path != tt.entry.Path {
					t.Errorf("Path mismatch: %q != %q", entry.Path, tt.entry.Path)
				}
			}
		})
	}
}

// TestPathEntryRoundTrip verifies round-trip serialization
func TestPathEntryRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		entry PathEntry
	}{
		{
			"Short path",
			PathEntry{
				PathLength: 8,
				Path:       "test.txt",
			},
		},
		{
			"Long path",
			PathEntry{
				PathLength: 31, // "very/long/path/to/test/file.txt" is 31 bytes
				Path:       "very/long/path/to/test/file.txt",
			},
		},
		{
			"Unicode path",
			PathEntry{
				PathLength: 17, // "测试/文件.txt" is 17 bytes in UTF-8
				Path:       "测试/文件.txt",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure PathLength matches actual path length
			tt.entry.PathLength = uint16(len(tt.entry.Path))

			// Write
			var buf bytes.Buffer
			if _, err := tt.entry.WriteTo(&buf); err != nil {
				t.Fatalf("WriteTo() error = %v", err)
			}

			// Read
			var entry PathEntry
			if _, err := entry.ReadFrom(&buf); err != nil {
				t.Fatalf("ReadFrom() error = %v", err)
			}

			// Compare all fields
			if entry.PathLength != tt.entry.PathLength {
				t.Errorf("PathLength mismatch: %d != %d", entry.PathLength, tt.entry.PathLength)
			}
			if entry.Path != tt.entry.Path {
				t.Errorf("Path mismatch: %q != %q", entry.Path, tt.entry.Path)
			}

			// Validate
			if err := entry.Validate(); err != nil {
				t.Errorf("Round-trip entry validation failed: %v", err)
			}
		})
	}
}

// TestPathEntryWriteToErrorPaths verifies WriteTo error handling
func TestPathEntryWriteToErrorPaths(t *testing.T) {
	tests := []struct {
		name      string
		entry     PathEntry
		writer    io.Writer
		wantErr   bool
		errSubstr string
	}{
		{
			"Error writer during PathLength write",
			PathEntry{
				PathLength: 8,
				Path:       "test.txt",
			},
			testhelpers.NewErrorWriter(),
			true,
			"failed to write",
		},
		{
			"Failing writer during PathLength write",
			PathEntry{
				PathLength: 8,
				Path:       "test.txt",
			},
			testhelpers.NewFailingWriter(1),
			true,
			"failed to write",
		},
		{
			"Failing writer during path write",
			PathEntry{
				PathLength: 8,
				Path:       "test.txt",
			},
			testhelpers.NewFailingWriter(2), // Can write PathLength but not Path
			true,
			"failed to write",
		},
		{
			"Incomplete path write",
			PathEntry{
				PathLength: 8,
				Path:       "test.txt",
			},
			testhelpers.NewIncompleteWriter(5), // 2 (PathLength) + 3 bytes of path (need 8)
			true,
			"incomplete path write",
		},
		{
			"Path length mismatch",
			PathEntry{
				PathLength: 20,
				Path:       "test.txt", // 8 bytes, but PathLength says 20
			},
			&bytes.Buffer{},
			true,
			"path length mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Only set PathLength if not testing mismatch
			if !strings.Contains(tt.name, "mismatch") {
				tt.entry.PathLength = uint16(len(tt.entry.Path))
			}
			_, err := tt.entry.WriteTo(tt.writer)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				if tt.errSubstr != "" {
					errStr := err.Error()
					if !strings.Contains(errStr, tt.errSubstr) {
						t.Errorf("WriteTo() error = %q, want error containing %q", errStr, tt.errSubstr)
					}
				}
			}
		})
	}
}

// TestPathEntry_GetPath tests the GetPath method.
func TestPathEntry_GetPath(t *testing.T) {
	tests := []struct {
		name     string
		entry    PathEntry
		wantPath string
	}{
		{
			"Path with leading slash",
			PathEntry{
				PathLength: 13,
				Path:       "/path/to/file",
			},
			"/path/to/file",
		},
		{
			"Path without leading slash",
			PathEntry{
				PathLength: 12,
				Path:       "path/to/file",
			},
			"path/to/file",
		},
		{
			"Root path",
			PathEntry{
				PathLength: 1,
				Path:       "/",
			},
			"/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.entry.GetPath(); got != tt.wantPath {
				t.Errorf("GetPath() = %q, want %q", got, tt.wantPath)
			}
		})
	}
}

// TestPathEntry_GetPathForPlatform tests the GetPathForPlatform method.
func TestPathEntry_GetPathForPlatform(t *testing.T) {
	tests := []struct {
		name      string
		entry     PathEntry
		isWindows bool
		want      string
	}{
		{
			"Unix path with leading slash",
			PathEntry{
				PathLength: 13,
				Path:       "/path/to/file",
			},
			false,
			"path/to/file",
		},
		{
			"Windows path with leading slash",
			PathEntry{
				PathLength: 13,
				Path:       "/path/to/file",
			},
			true,
			"path\\to\\file",
		},
		{
			"Unix path without leading slash",
			PathEntry{
				PathLength: 12,
				Path:       "path/to/file",
			},
			false,
			"path/to/file",
		},
		{
			"Windows path without leading slash",
			PathEntry{
				PathLength: 12,
				Path:       "path/to/file",
			},
			true,
			"path\\to\\file",
		},
		{
			"Unix root path",
			PathEntry{
				PathLength: 1,
				Path:       "/",
			},
			false,
			"",
		},
		{
			"Windows root path",
			PathEntry{
				PathLength: 1,
				Path:       "/",
			},
			true,
			"",
		},
		{
			"Unix nested path",
			PathEntry{
				PathLength: 20,
				Path:       "/very/long/path/file",
			},
			false,
			"very/long/path/file",
		},
		{
			"Windows nested path",
			PathEntry{
				PathLength: 20,
				Path:       "/very/long/path/file",
			},
			true,
			"very\\long\\path\\file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.entry.GetPathForPlatform(tt.isWindows); got != tt.want {
				t.Errorf("GetPathForPlatform(%v) = %q, want %q", tt.isWindows, got, tt.want)
			}
		})
	}
}
