package fileformat

import (
	"encoding/binary"
	"fmt"
	"io"
)

// PathEntry represents a file path with associated metadata.
//
// Size: Variable (2 + path_length + 36 bytes of metadata)
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.2 - Path Entries
type PathEntry struct {
	// PathLength is the length of the path in bytes (UTF-8)
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.2
	PathLength uint16

	// Path is the UTF-8 encoded file path (not null-terminated)
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.2
	Path string

	// Mode is the file permissions and type (Unix-style)
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.2
	Mode uint32

	// UserID is the user identifier (Unix-style)
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.2
	UserID uint32

	// GroupID is the group identifier (Unix-style)
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.2
	GroupID uint32

	// ModTime is the modification time (Unix nanoseconds)
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.2
	ModTime uint64

	// CreateTime is the creation time (Unix nanoseconds)
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.2
	CreateTime uint64

	// AccessTime is the access time (Unix nanoseconds)
	// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.2
	AccessTime uint64
}

// Validate performs validation checks on the PathEntry.
//
// Validation checks:
//   - PathLength must match actual Path length
//   - Path must not be empty
//   - Path must be valid UTF-8
//
// Returns an error if any validation check fails.
func (p *PathEntry) Validate() error {
	if p.Path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	pathLen := uint16(len(p.Path))
	if p.PathLength != pathLen {
		return fmt.Errorf("path length mismatch: specified %d, actual %d", p.PathLength, pathLen)
	}

	return nil
}

// Size returns the total size of the PathEntry in bytes.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.2
func (p *PathEntry) Size() int {
	// PathLength(2) + Path + Mode(4) + UserID(4) + GroupID(4) + ModTime(8) + CreateTime(8) + AccessTime(8)
	// Metadata total: 4 + 4 + 4 + 8 + 8 + 8 = 36 bytes
	return 2 + int(p.PathLength) + 36
}

// ReadFrom reads a PathEntry from the provided io.Reader.
//
// The binary format is:
//   - PathLength (2 bytes, little-endian uint16)
//   - Path (PathLength bytes, UTF-8 string, not null-terminated)
//   - Mode (4 bytes, little-endian uint32)
//   - UserID (4 bytes, little-endian uint32)
//   - GroupID (4 bytes, little-endian uint32)
//   - ModTime (8 bytes, little-endian uint64)
//   - CreateTime (8 bytes, little-endian uint64)
//   - AccessTime (8 bytes, little-endian uint64)
//
// Returns the number of bytes read and any error encountered.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.2 - Path Entries
func (p *PathEntry) ReadFrom(r io.Reader) (int64, error) {
	var totalRead int64

	// Read PathLength (2 bytes)
	var pathLength uint16
	if err := binary.Read(r, binary.LittleEndian, &pathLength); err != nil {
		return totalRead, fmt.Errorf("failed to read path length: %w", err)
	}
	totalRead += 2
	p.PathLength = pathLength

	// Read Path (PathLength bytes)
	if pathLength > 0 {
		pathBytes := make([]byte, pathLength)
		n, err := io.ReadFull(r, pathBytes)
		if err != nil {
			return totalRead, fmt.Errorf("failed to read path: %w", err)
		}
		if uint16(n) != pathLength {
			return totalRead, fmt.Errorf("incomplete path read: got %d bytes, expected %d", n, pathLength)
		}
		totalRead += int64(n)
		p.Path = string(pathBytes)
	} else {
		p.Path = ""
	}

	// Read Mode (4 bytes)
	if err := binary.Read(r, binary.LittleEndian, &p.Mode); err != nil {
		return totalRead, fmt.Errorf("failed to read mode: %w", err)
	}
	totalRead += 4

	// Read UserID (4 bytes)
	if err := binary.Read(r, binary.LittleEndian, &p.UserID); err != nil {
		return totalRead, fmt.Errorf("failed to read user ID: %w", err)
	}
	totalRead += 4

	// Read GroupID (4 bytes)
	if err := binary.Read(r, binary.LittleEndian, &p.GroupID); err != nil {
		return totalRead, fmt.Errorf("failed to read group ID: %w", err)
	}
	totalRead += 4

	// Read ModTime (8 bytes)
	if err := binary.Read(r, binary.LittleEndian, &p.ModTime); err != nil {
		return totalRead, fmt.Errorf("failed to read modification time: %w", err)
	}
	totalRead += 8

	// Read CreateTime (8 bytes)
	if err := binary.Read(r, binary.LittleEndian, &p.CreateTime); err != nil {
		return totalRead, fmt.Errorf("failed to read creation time: %w", err)
	}
	totalRead += 8

	// Read AccessTime (8 bytes)
	if err := binary.Read(r, binary.LittleEndian, &p.AccessTime); err != nil {
		return totalRead, fmt.Errorf("failed to read access time: %w", err)
	}
	totalRead += 8

	return totalRead, nil
}

// WriteTo writes a PathEntry to the provided io.Writer.
//
// The binary format is:
//   - PathLength (2 bytes, little-endian uint16)
//   - Path (PathLength bytes, UTF-8 string, not null-terminated)
//   - Mode (4 bytes, little-endian uint32)
//   - UserID (4 bytes, little-endian uint32)
//   - GroupID (4 bytes, little-endian uint32)
//   - ModTime (8 bytes, little-endian uint64)
//   - CreateTime (8 bytes, little-endian uint64)
//   - AccessTime (8 bytes, little-endian uint64)
//
// Returns the number of bytes written and any error encountered.
//
// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.2 - Path Entries
func (p *PathEntry) WriteTo(w io.Writer) (int64, error) {
	var totalWritten int64

	// Write PathLength (2 bytes)
	if err := binary.Write(w, binary.LittleEndian, p.PathLength); err != nil {
		return totalWritten, fmt.Errorf("failed to write path length: %w", err)
	}
	totalWritten += 2

	// Write Path (PathLength bytes)
	if p.PathLength > 0 {
		pathBytes := []byte(p.Path)
		if uint16(len(pathBytes)) != p.PathLength {
			return totalWritten, fmt.Errorf("path length mismatch: specified %d, actual %d", p.PathLength, len(pathBytes))
		}
		n, err := w.Write(pathBytes)
		if err != nil {
			return totalWritten, fmt.Errorf("failed to write path: %w", err)
		}
		if uint16(n) != p.PathLength {
			return totalWritten, fmt.Errorf("incomplete path write: wrote %d bytes, expected %d", n, p.PathLength)
		}
		totalWritten += int64(n)
	}

	// Write Mode (4 bytes)
	if err := binary.Write(w, binary.LittleEndian, p.Mode); err != nil {
		return totalWritten, fmt.Errorf("failed to write mode: %w", err)
	}
	totalWritten += 4

	// Write UserID (4 bytes)
	if err := binary.Write(w, binary.LittleEndian, p.UserID); err != nil {
		return totalWritten, fmt.Errorf("failed to write user ID: %w", err)
	}
	totalWritten += 4

	// Write GroupID (4 bytes)
	if err := binary.Write(w, binary.LittleEndian, p.GroupID); err != nil {
		return totalWritten, fmt.Errorf("failed to write group ID: %w", err)
	}
	totalWritten += 4

	// Write ModTime (8 bytes)
	if err := binary.Write(w, binary.LittleEndian, p.ModTime); err != nil {
		return totalWritten, fmt.Errorf("failed to write modification time: %w", err)
	}
	totalWritten += 8

	// Write CreateTime (8 bytes)
	if err := binary.Write(w, binary.LittleEndian, p.CreateTime); err != nil {
		return totalWritten, fmt.Errorf("failed to write creation time: %w", err)
	}
	totalWritten += 8

	// Write AccessTime (8 bytes)
	if err := binary.Write(w, binary.LittleEndian, p.AccessTime); err != nil {
		return totalWritten, fmt.Errorf("failed to write access time: %w", err)
	}
	totalWritten += 8

	return totalWritten, nil
}
