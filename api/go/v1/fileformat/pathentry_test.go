package fileformat

import (
	"bytes"
	"encoding/binary"
	"io"
	"strings"
	"testing"
)

// TestPathEntrySize verifies PathEntry size calculation
func TestPathEntrySize(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wantSize int
	}{
		{"Short path", "test.txt", 46},                                           // 2 (PathLength) + 8 (Path) + 36 (metadata)
		{"Medium path", "path/to/file.txt", 54},                                  // 2 + 16 + 36
		{"Long path", "very/long/path/to/some/file/in/a/deep/directory.txt", 89}, // 2 + 51 + 36
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
// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.2 - Path Entries
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
				Mode:       0644,
				UserID:     1000,
				GroupID:    1000,
				ModTime:    1638360000000000000,
				CreateTime: 1638360000000000000,
				AccessTime: 1638360000000000000,
			},
			false,
		},
		{
			"Path entry with long path",
			PathEntry{
				PathLength: 20,
				Path:       "path/to/test/file.txt",
				Mode:       0755,
				UserID:     2000,
				GroupID:    2000,
				ModTime:    1638361000000000000,
				CreateTime: 1638360000000000000,
				AccessTime: 1638362000000000000,
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
				// Calculate expected size: 2 (PathLength) + pathLen + 36 (metadata)
				expectedSize := 2 + int(tt.entry.PathLength) + 36
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
				if entry.Mode != tt.entry.Mode {
					t.Errorf("Mode = %d, want %d", entry.Mode, tt.entry.Mode)
				}
				if entry.UserID != tt.entry.UserID {
					t.Errorf("UserID = %d, want %d", entry.UserID, tt.entry.UserID)
				}
				if entry.GroupID != tt.entry.GroupID {
					t.Errorf("GroupID = %d, want %d", entry.GroupID, tt.entry.GroupID)
				}
				if entry.ModTime != tt.entry.ModTime {
					t.Errorf("ModTime = %d, want %d", entry.ModTime, tt.entry.ModTime)
				}
				if entry.CreateTime != tt.entry.CreateTime {
					t.Errorf("CreateTime = %d, want %d", entry.CreateTime, tt.entry.CreateTime)
				}
				if entry.AccessTime != tt.entry.AccessTime {
					t.Errorf("AccessTime = %d, want %d", entry.AccessTime, tt.entry.AccessTime)
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
		{"Only PathLength", []byte{0x08, 0x00}},
		{"Incomplete path", []byte{0x08, 0x00, 0x74, 0x65}}, // "te"
		{"Missing metadata", []byte{0x08, 0x00, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x74, 0x78, 0x74}}, // path only
		{"Incomplete Mode", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint16(8))
			buf.WriteString("test.txt")
			_ = binary.Write(buf, binary.LittleEndian, uint32(0644))
			// Missing UserID, GroupID, times
			return buf.Bytes()[:14] // PathLength + Path + partial Mode
		}()},
		{"Incomplete UserID", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint16(8))
			buf.WriteString("test.txt")
			_ = binary.Write(buf, binary.LittleEndian, uint32(0644))
			_ = binary.Write(buf, binary.LittleEndian, uint32(1000))
			// Missing GroupID, times
			return buf.Bytes()[:18] // Up to partial UserID
		}()},
		{"Incomplete GroupID", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint16(8))
			buf.WriteString("test.txt")
			_ = binary.Write(buf, binary.LittleEndian, uint32(0644))
			_ = binary.Write(buf, binary.LittleEndian, uint32(1000))
			_ = binary.Write(buf, binary.LittleEndian, uint32(1000))
			// Missing times
			return buf.Bytes()[:22] // Up to partial GroupID
		}()},
		{"Incomplete ModTime", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint16(8))
			buf.WriteString("test.txt")
			_ = binary.Write(buf, binary.LittleEndian, uint32(0644))
			_ = binary.Write(buf, binary.LittleEndian, uint32(1000))
			_ = binary.Write(buf, binary.LittleEndian, uint32(1000))
			_ = binary.Write(buf, binary.LittleEndian, uint64(1638360000000000000))
			// Missing CreateTime, AccessTime
			return buf.Bytes()[:30] // Up to partial ModTime
		}()},
		{"Incomplete CreateTime", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint16(8))
			buf.WriteString("test.txt")
			_ = binary.Write(buf, binary.LittleEndian, uint32(0644))
			_ = binary.Write(buf, binary.LittleEndian, uint32(1000))
			_ = binary.Write(buf, binary.LittleEndian, uint32(1000))
			_ = binary.Write(buf, binary.LittleEndian, uint64(1638360000000000000))
			_ = binary.Write(buf, binary.LittleEndian, uint64(1638360000000000000))
			// Missing AccessTime
			return buf.Bytes()[:38] // Up to partial CreateTime
		}()},
		{"Incomplete AccessTime", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint16(8))
			buf.WriteString("test.txt")
			_ = binary.Write(buf, binary.LittleEndian, uint32(0644))
			_ = binary.Write(buf, binary.LittleEndian, uint32(1000))
			_ = binary.Write(buf, binary.LittleEndian, uint32(1000))
			_ = binary.Write(buf, binary.LittleEndian, uint64(1638360000000000000))
			_ = binary.Write(buf, binary.LittleEndian, uint64(1638360000000000000))
			// Partial AccessTime (need 8 bytes, only have 7)
			return buf.Bytes()[:45] // Up to partial AccessTime
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
		Mode:       0644,
		UserID:     1000,
		GroupID:    1000,
		ModTime:    1638360000000000000,
		CreateTime: 1638360000000000000,
		AccessTime: 1638362000000000000,
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
// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.2 - Path Entries
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
				Mode:       0644,
				UserID:     1000,
				GroupID:    1000,
				ModTime:    1638360000000000000,
				CreateTime: 1638360000000000000,
				AccessTime: 1638360000000000000,
			},
			false,
		},
		{
			"Path entry with long path",
			PathEntry{
				PathLength: 20,
				Path:       "path/to/test/file.txt",
				Mode:       0755,
				UserID:     2000,
				GroupID:    2000,
				ModTime:    1638361000000000000,
				CreateTime: 1638360000000000000,
				AccessTime: 1638362000000000000,
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
				if entry.Mode != tt.entry.Mode {
					t.Errorf("Mode mismatch: %d != %d", entry.Mode, tt.entry.Mode)
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
				Mode:       0644,
				UserID:     1000,
				GroupID:    1000,
				ModTime:    1638360000000000000,
				CreateTime: 1638360000000000000,
				AccessTime: 1638360000000000000,
			},
		},
		{
			"Long path",
			PathEntry{
				PathLength: 31, // "very/long/path/to/test/file.txt" is 31 bytes
				Path:       "very/long/path/to/test/file.txt",
				Mode:       0755,
				UserID:     2000,
				GroupID:    2000,
				ModTime:    1638361000000000000,
				CreateTime: 1638360000000000000,
				AccessTime: 1638362000000000000,
			},
		},
		{
			"Unicode path",
			PathEntry{
				PathLength: 17, // "测试/文件.txt" is 17 bytes in UTF-8
				Path:       "测试/文件.txt",
				Mode:       0644,
				UserID:     1000,
				GroupID:    1000,
				ModTime:    1638360000000000000,
				CreateTime: 1638360000000000000,
				AccessTime: 1638360000000000000,
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
			if entry.Mode != tt.entry.Mode {
				t.Errorf("Mode mismatch: %d != %d", entry.Mode, tt.entry.Mode)
			}
			if entry.UserID != tt.entry.UserID {
				t.Errorf("UserID mismatch: %d != %d", entry.UserID, tt.entry.UserID)
			}
			if entry.GroupID != tt.entry.GroupID {
				t.Errorf("GroupID mismatch: %d != %d", entry.GroupID, tt.entry.GroupID)
			}
			if entry.ModTime != tt.entry.ModTime {
				t.Errorf("ModTime mismatch: %d != %d", entry.ModTime, tt.entry.ModTime)
			}
			if entry.CreateTime != tt.entry.CreateTime {
				t.Errorf("CreateTime mismatch: %d != %d", entry.CreateTime, tt.entry.CreateTime)
			}
			if entry.AccessTime != tt.entry.AccessTime {
				t.Errorf("AccessTime mismatch: %d != %d", entry.AccessTime, tt.entry.AccessTime)
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
			&errorWriter{},
			true,
			"failed to write",
		},
		{
			"Failing writer during path write",
			PathEntry{
				PathLength: 8,
				Path:       "test.txt",
			},
			&failingWriter{maxBytes: 1},
			true,
			"failed to write",
		},
		{
			"Incomplete path write",
			PathEntry{
				PathLength: 8,
				Path:       "test.txt",
			},
			&incompleteWriter{maxWrite: 5},
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
		{
			"Failing writer during Mode write",
			PathEntry{
				PathLength: 8,
				Path:       "test.txt",
			},
			&failingWriter{maxBytes: 10},
			true,
			"failed to write",
		},
		{
			"Failing writer during UserID write",
			PathEntry{
				PathLength: 8,
				Path:       "test.txt",
			},
			&failingWriter{maxBytes: 14},
			true,
			"failed to write",
		},
		{
			"Failing writer during GroupID write",
			PathEntry{
				PathLength: 8,
				Path:       "test.txt",
			},
			&failingWriter{maxBytes: 18},
			true,
			"failed to write",
		},
		{
			"Failing writer during ModTime write",
			PathEntry{
				PathLength: 8,
				Path:       "test.txt",
			},
			&failingWriter{maxBytes: 22},
			true,
			"failed to write",
		},
		{
			"Failing writer during CreateTime write",
			PathEntry{
				PathLength: 8,
				Path:       "test.txt",
			},
			&failingWriter{maxBytes: 30},
			true,
			"failed to write",
		},
		{
			"Failing writer during AccessTime write",
			PathEntry{
				PathLength: 8,
				Path:       "test.txt",
			},
			&failingWriter{maxBytes: 38},
			true,
			"failed to write",
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
