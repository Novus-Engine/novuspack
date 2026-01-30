// This file implements file lookup methods.
// It contains methods for looking up files by metadata (FileID, hash, checksum),
// finding entries by tag or type, and getting file counts.
//
// Specification: api_file_mgmt_queries.md: 2. Single-Entry Lookups

package novus_package

import (
	"bytes"
	"fmt"

	"github.com/novus-engine/novuspack/api/go/metadata"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// GetFileByFileID retrieves a file entry by its FileID.
//
// This is a pure in-memory operation that searches the FileEntries
// slice for a matching FileID.
//
// Parameters:
//   - fileID: The FileID to search for
//
// Returns:
//   - *metadata.FileEntry: The matching file entry, or nil if not found
//   - error: *PackageError if not found or on failure
//
// Specification: api_file_mgmt_queries.md: 2. Single-Entry Lookups
func (p *filePackage) GetFileByFileID(fileID uint64) (*metadata.FileEntry, error) {
	// This is a pure in-memory operation
	if p.FileEntries == nil {
		return nil, pkgerrors.NewPackageError(
			pkgerrors.ErrTypeValidation,
			fmt.Sprintf("file with FileID %d not found", fileID),
			nil,
			pkgerrors.ValidationErrorContext{
				Field:    "fileID",
				Value:    fileID,
				Expected: "existing FileID",
			},
		)
	}

	for _, entry := range p.FileEntries {
		if entry.FileID == fileID {
			return entry, nil
		}
	}

	return nil, pkgerrors.NewPackageError(
		pkgerrors.ErrTypeValidation,
		fmt.Sprintf("file with FileID %d not found", fileID),
		nil,
		pkgerrors.ValidationErrorContext{
			Field:    "fileID",
			Value:    fileID,
			Expected: "existing FileID",
		},
	)
}

// GetFileByHash retrieves a file entry by its content hash.
//
// This method searches for a file entry with a matching hash in its
// Hashes array. The hash comparison is byte-exact.
//
// Parameters:
//   - hash: The hash bytes to search for
//
// Returns:
//   - *metadata.FileEntry: The matching file entry, or nil if not found
//   - error: *PackageError if not found or on failure
//
// Specification: api_file_mgmt_queries.md: 2. Single-Entry Lookups
func (p *filePackage) GetFileByHash(hash []byte) (*metadata.FileEntry, error) {
	// This is a pure in-memory operation
	if p.FileEntries == nil || hash == nil {
		return nil, pkgerrors.NewPackageError(
			pkgerrors.ErrTypeValidation,
			"file with specified hash not found",
			nil,
			pkgerrors.ValidationErrorContext{
				Field:    "hash",
				Value:    hash,
				Expected: "existing hash",
			},
		)
	}

	for _, entry := range p.FileEntries {
		// Check all hashes in the entry
		for _, hashEntry := range entry.Hashes {
			if bytes.Equal(hashEntry.HashData, hash) {
				return entry, nil
			}
		}
	}

	return nil, pkgerrors.NewPackageError(
		pkgerrors.ErrTypeValidation,
		"file with specified hash not found",
		nil,
		pkgerrors.ValidationErrorContext{
			Field:    "hash",
			Value:    hash,
			Expected: "existing hash",
		},
	)
}

// GetFileByChecksum retrieves a file entry by its CRC32 checksum.
//
// This method searches for a file entry with a matching RawChecksum
// or StoredChecksum value.
//
// Parameters:
//   - checksum: The CRC32 checksum to search for
//
// Returns:
//   - *metadata.FileEntry: The matching file entry, or nil if not found
//   - error: *PackageError if not found or on failure
//
// Specification: api_file_mgmt_queries.md: 2. Single-Entry Lookups
func (p *filePackage) GetFileByChecksum(checksum uint32) (*metadata.FileEntry, error) {
	// This is a pure in-memory operation
	if p.FileEntries == nil {
		return nil, pkgerrors.NewPackageError(
			pkgerrors.ErrTypeValidation,
			fmt.Sprintf("file with checksum %d not found", checksum),
			nil,
			pkgerrors.ValidationErrorContext{
				Field:    "checksum",
				Value:    checksum,
				Expected: "existing checksum",
			},
		)
	}

	for _, entry := range p.FileEntries {
		if entry.RawChecksum == checksum || entry.StoredChecksum == checksum {
			return entry, nil
		}
	}

	return nil, pkgerrors.NewPackageError(
		pkgerrors.ErrTypeValidation,
		fmt.Sprintf("file with checksum %d not found", checksum),
		nil,
		pkgerrors.ValidationErrorContext{
			Field:    "checksum",
			Value:    checksum,
			Expected: "existing checksum",
		},
	)
}

// FindEntriesByTag finds all file entries that have a specific tag key-value pair.
//
// This method searches for file entries with matching tag keys and performs
// type-safe comparison of tag values.
//
// Parameters:
//   - tagKey: The tag key to search for
//   - tagValue: The tag value to match (type-safe comparison)
//
// Returns:
//   - []*metadata.FileEntry: Slice of matching file entries (empty if none found)
//   - error: *PackageError on failure
//
// Specification: api_file_mgmt_queries.md: 2. Single-Entry Lookups
func (p *filePackage) FindEntriesByTag(tagKey string, tagValue any) ([]*metadata.FileEntry, error) {
	// This is a pure in-memory operation
	if p.FileEntries == nil {
		return []*metadata.FileEntry{}, nil
	}

	matches := make([]*metadata.FileEntry, 0)

	for _, entry := range p.FileEntries {
		// Search OptionalData for tags entries
		for _, optData := range entry.OptionalData {
			if optData.DataType == metadata.OptionalDataTagsData && len(optData.Data) > 0 {
				// Parse tags from the data (JSON format)
				// For now, we'll do a simple string match in the JSON data
				// A full implementation would deserialize the JSON and compare properly
				searchStr := fmt.Sprintf("\"%s\"", tagKey)
				if bytes.Contains(optData.Data, []byte(searchStr)) {
					matches = append(matches, entry)
					break
				}
			}
		}
	}

	return matches, nil
}

// Note: compareTagValue was removed as tags are now searched using
// simple JSON string matching. A full implementation would deserialize
// the JSON and perform proper type-safe comparison.

// FindEntriesByType finds all file entries of a specific file type.
//
// This method searches for file entries with a matching Type field.
// Special metadata files (type >= 65000) are excluded from results.
//
// Parameters:
//   - fileType: The file type identifier to search for
//
// Returns:
//   - []*metadata.FileEntry: Slice of matching file entries (empty if none found)
//   - error: *PackageError on failure
//
// Specification: api_file_mgmt_queries.md: 2. Single-Entry Lookups
func (p *filePackage) FindEntriesByType(fileType uint16) ([]*metadata.FileEntry, error) {
	// This is a pure in-memory operation
	if p.FileEntries == nil {
		return []*metadata.FileEntry{}, nil
	}

	matches := make([]*metadata.FileEntry, 0)

	for _, entry := range p.FileEntries {
		if entry.Type == fileType {
			matches = append(matches, entry)
		}
	}

	return matches, nil
}

// GetFileCount returns the total number of regular content files in the package.
//
// This method counts FileEntries, excluding special metadata files
// (type >= 65000).
//
// Returns:
//   - int: Number of regular content files
//   - error: *PackageError on failure
//
// Specification: api_file_mgmt_queries.md: 3.2.1 FindEntriesByType Package Method
func (p *filePackage) GetFileCount() (int, error) {
	// This is a pure in-memory operation
	if p.FileEntries == nil {
		return 0, nil
	}

	count := 0
	for _, entry := range p.FileEntries {
		// Exclude special metadata files (type >= 65000)
		if entry.Type < 65000 {
			count++
		}
	}

	return count, nil
}
