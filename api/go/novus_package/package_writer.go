// This file implements Package write operations: Write, SafeWrite, FastWrite.
// It contains all write operations for persisting package contents to disk as specified
// in api_core.md and api_writing.md.
//
// Note: StageFile and UnstageFile have been removed from the Package API
// as per Priority 0 requirements. File management now uses AddFile/AddFileFromMemory/
// RemoveFile methods instead.
//
// Specification: api_writing.md: 1. SafeWrite - Atomic Package Writing

// Package novuspack provides the NovusPack API v1 implementation.
//
// This file contains package writer operations: Write, SafeWrite, FastWrite, and Defragment.
package novus_package

import (
	"context"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/novus-engine/novuspack/api/go/fileformat"
	"github.com/novus-engine/novuspack/api/go/internal"
	"github.com/novus-engine/novuspack/api/go/metadata"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// Write writes the package to disk.
//
// This method writes the package to the path configured via Create() or CreateWithOptions().
// Compression and signing options are determined by package state rather than method parameters.
//
// This baseline implementation uses SafeWrite with overwrite=true.
//
// Specification: api_core.md: 1.3 Package Write Operations
// Specification: api_writing.md: 1. SafeWrite - Atomic Package Writing
func (p *filePackage) Write(ctx context.Context) error {
	// Validate context
	if err := internal.CheckContext(ctx, "Write"); err != nil {
		return err
	}

	// Ensure path metadata special file is up to date
	if err := p.SavePathMetadataFile(ctx); err != nil {
		return pkgerrors.WrapError(err, pkgerrors.ErrTypeIO, "failed to save path metadata before Write")
	}

	// Use SafeWrite with overwrite=true as the baseline strategy
	return p.SafeWrite(ctx, true)
}

// SafeWrite writes the package to disk safely with atomic operations.
//
// This method writes the package to the path configured via Create() or CreateWithOptions(),
// using atomic operations (temporary file + rename) to ensure data integrity.
// Compression and signing options are determined by package state.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - overwrite: Whether to overwrite existing file (false = fail if exists)
//
// This baseline implementation writes uncompressed and unencrypted files only.
//
// Specification: api_core.md: 1.3 Package Write Operations
// Specification: api_writing.md: 5.3.3 Package.Write Method
func (p *filePackage) SafeWrite(ctx context.Context, overwrite bool) error {
	// Validate context
	if err := internal.CheckContext(ctx, "SafeWrite"); err != nil {
		return err
	}

	// Validate package has a file path configured
	if p.FilePath == "" {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "package has no file path configured", nil, struct{}{})
	}

	// Check if file exists and handle overwrite flag
	if !overwrite {
		if _, err := os.Stat(p.FilePath); err == nil {
			return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "file already exists and overwrite is false", nil, pkgerrors.ValidationErrorContext{
				Field:    "FilePath",
				Value:    p.FilePath,
				Expected: "non-existing file or overwrite=true",
			})
		}
	}

	// Create temp file in same directory as target (for atomic rename)
	tempFile, err := os.CreateTemp(filepath.Dir(p.FilePath), ".nvpk-temp-*")
	if err != nil {
		return pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to create temp file", pkgerrors.ValidationErrorContext{
			Field: "FilePath",
			Value: p.FilePath,
		})
	}
	tempPath := tempFile.Name()

	// Ensure temp file is cleaned up on error
	var writeErr error
	defer func() {
		if writeErr != nil {
			_ = tempFile.Close()
			_ = os.Remove(tempPath)
		}
	}()

	// Write package to temp file
	if err := p.writePackageToFile(ctx, tempFile); err != nil {
		writeErr = err
		return err
	}

	// Close temp file before rename
	if err := tempFile.Close(); err != nil {
		writeErr = pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to close temp file", pkgerrors.ValidationErrorContext{
			Field: "TempFile",
			Value: tempPath,
		})
		return writeErr
	}

	// Atomic rename to target path
	if err := os.Rename(tempPath, p.FilePath); err != nil {
		writeErr = pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to rename temp file to target", pkgerrors.ValidationErrorContext{
			Field:    "FilePath",
			Value:    p.FilePath,
			Expected: "successful rename",
		})
		return writeErr
	}

	return nil
}

// FastWrite writes the package to disk quickly without atomic operations.
//
// This method writes the package to the path configured via Create() or CreateWithOptions(),
// using in-place updates for better performance. Compression and signing options are
// determined by package state.
//
// TODO: Implement fast package writing with state-driven compression and signing.
//
// Specification: api_core.md: 1.3 Package Write Operations
// Specification: api_writing.md: 1. SafeWrite - Atomic Package Writing
func (p *filePackage) FastWrite(ctx context.Context) error {
	// TODO: Implement fast package writing
	// TODO: Use package state for compression and signing configuration
	return pkgerrors.NewPackageError(pkgerrors.ErrTypeUnsupported, "FastWrite not yet implemented", nil, struct{}{})
}

// Defragment defragments the package to optimize storage.
//
// TODO: Implement package defragmentation.
//
// Specification: api_basic_operations.md: 16. Package.Defragment Method
func (p *filePackage) Defragment(ctx context.Context) error {
	// TODO: Implement defragmentation
	return pkgerrors.NewPackageError(pkgerrors.ErrTypeUnsupported, "Defragment not yet implemented", nil, struct{}{})
}

// writePackageToFile writes the complete package structure to a file.
//
// This helper method writes:
//   - Header (placeholder, updated at end)
//   - Interleaved file entry metadata and data
//   - File index
//   - Package comment (if present)
//   - Updated header with correct offsets
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - file: File to write to
//
// Returns:
//   - error: *PackageError on failure
//
//nolint:gocognit,gocyclo // write sequence branches
func (p *filePackage) writePackageToFile(ctx context.Context, file *os.File) error {
	// Update PackageInfo (canonical in-memory metadata)
	if p.Info == nil {
		p.Info = metadata.NewPackageInfo()
	}

	// Set timestamps in PackageInfo
	now := time.Now()
	if p.Info.Created.IsZero() {
		p.Info.Created = now
	}
	p.Info.Modified = now

	// Initialize or update header from PackageInfo
	if p.header == nil {
		p.header = fileformat.NewPackageHeader()
		p.header.Magic = fileformat.NVPKMagic
		p.header.FormatVersion = fileformat.FormatVersion
	}

	// Sync PackageInfo to header (PackageInfo is the source of truth)
	p.header.CreatedTime = uint64(p.Info.Created.UnixNano())
	p.header.ModifiedTime = uint64(p.Info.Modified.UnixNano())
	p.header.VendorID = p.Info.VendorID
	p.header.AppID = p.Info.AppID
	p.header.PackageDataVersion = p.Info.PackageDataVersion
	p.header.MetadataVersion = p.Info.MetadataVersion

	// Write placeholder header (we'll update it later with correct offsets)
	headerSize := int64(fileformat.PackageHeaderSize)
	if _, err := writePackageHeader(file, p.header); err != nil {
		return pkgerrors.WrapError(err, pkgerrors.ErrTypeIO, "failed to write header placeholder")
	}

	// Track current offset (after header)
	currentOffset := uint64(headerSize)

	// Build file index as we write entries
	index := fileformat.NewFileIndex()
	index.FirstEntryOffset = currentOffset
	index.Entries = make([]fileformat.IndexEntry, 0, len(p.FileEntries))

	// Write interleaved file entry metadata and file data
	for _, fe := range p.FileEntries {
		if fe == nil {
			continue
		}

		if fe.IsDataLoaded {
			p.syncStoredMetadataFromMemory(fe)
		}

		// Record entry offset in index
		index.Entries = append(index.Entries, fileformat.IndexEntry{
			FileID: fe.FileID,
			Offset: currentOffset,
		})

		// Write file entry metadata
		entryOffset := currentOffset
		metaWritten, err := fe.WriteMetaTo(file)
		if err != nil {
			return pkgerrors.WrapError(err, pkgerrors.ErrTypeIO, "failed to write file entry metadata")
		}
		currentOffset += uint64(metaWritten)

		// Write file data
		switch {
		case fe.IsDataLoaded:
			// Write in-memory data
			n, err := file.Write(fe.Data)
			if err != nil {
				return pkgerrors.WrapError(err, pkgerrors.ErrTypeIO, "failed to write file data")
			}
			currentOffset += uint64(n)
		case fe.SourceFile != nil:
			dataSize := fe.SourceSize
			if dataSize == 0 {
				dataSize = int64(fe.OriginalSize)
			}

			if _, err := fe.SourceFile.Seek(fe.SourceOffset, io.SeekStart); err != nil {
				return pkgerrors.WrapError(err, pkgerrors.ErrTypeIO, "failed to seek to source file data")
			}

			needsChecksums := fe.RawChecksum == 0 || fe.StoredChecksum == 0 || fe.StoredSize == 0
			if needsChecksums {
				hasher := crc32.NewIEEE()
				writer := io.MultiWriter(file, hasher)
				n, err := io.CopyN(writer, fe.SourceFile, dataSize)
				if err != nil {
					return pkgerrors.WrapError(err, pkgerrors.ErrTypeIO, "failed to copy file data from source")
				}
				if n != dataSize {
					return pkgerrors.NewPackageError(pkgerrors.ErrTypeCorruption, "source file size mismatch during write", nil, pkgerrors.ValidationErrorContext{
						Field:    "SourceSize",
						Value:    n,
						Expected: fmt.Sprintf("%d", dataSize),
					})
				}

				checksum := hasher.Sum32()
				if fe.RawChecksum == 0 {
					fe.RawChecksum = checksum
				}
				if fe.StoredChecksum == 0 {
					fe.StoredChecksum = checksum
				}
				if fe.StoredSize == 0 {
					fe.StoredSize = uint64(n)
				}

				currentOffset += uint64(n)

				if err := p.rewriteFileEntryMeta(file, entryOffset, fe); err != nil {
					return err
				}
				if _, err := file.Seek(int64(currentOffset), io.SeekStart); err != nil {
					return pkgerrors.WrapError(err, pkgerrors.ErrTypeIO, "failed to seek back to end after metadata rewrite")
				}
			} else {
				n, err := io.CopyN(file, fe.SourceFile, dataSize)
				if err != nil {
					return pkgerrors.WrapError(err, pkgerrors.ErrTypeIO, "failed to copy file data from source")
				}
				currentOffset += uint64(n)
			}
		case p.fileHandle != nil:
			// Stream existing file data from opened package
			// Find data offset in source file
			var sourceDataOffset uint64
			for _, indexEntry := range p.index.Entries {
				if indexEntry.FileID == fe.FileID {
					sourceDataOffset = indexEntry.Offset + uint64(fe.TotalSize())
					break
				}
			}

			if sourceDataOffset > 0 {
				// Seek to source data
				if _, err := p.fileHandle.Seek(int64(sourceDataOffset), 0); err != nil {
					return pkgerrors.WrapError(err, pkgerrors.ErrTypeIO, "failed to seek to source file data")
				}

				// Copy data from source to target
				n, err := io.CopyN(file, p.fileHandle, int64(fe.StoredSize))
				if err != nil {
					return pkgerrors.WrapError(err, pkgerrors.ErrTypeIO, "failed to copy file data from source")
				}
				currentOffset += uint64(n)
			}
		default:
			return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "file entry has no data source for writing", nil, pkgerrors.ValidationErrorContext{
				Field:    "FileEntry",
				Value:    fe.FileID,
				Expected: "data loaded, source file, or open package handle",
			})
		}

		// Check context cancellation periodically
		select {
		case <-ctx.Done():
			return pkgerrors.NewPackageError(pkgerrors.ErrTypeContext, "context cancelled during write", ctx.Err(), struct{}{})
		default:
		}
	}

	// Update index metadata
	index.EntryCount = uint32(len(index.Entries))

	// Write file index
	indexStart := currentOffset
	indexWritten, err := writeFileIndexTo(file, index)
	if err != nil {
		return pkgerrors.WrapError(err, pkgerrors.ErrTypeIO, "failed to write file index")
	}
	currentOffset += uint64(indexWritten)

	// Update header with index information
	p.header.IndexStart = indexStart
	p.header.IndexSize = uint64(indexWritten)

	// Write package comment if present (using PackageComment binary format)
	if p.Info != nil && p.Info.Comment != "" {
		commentStart := currentOffset

		comment, err := buildPackageComment(p.Info.Comment)
		if err != nil {
			return pkgerrors.WrapError(err, pkgerrors.ErrTypeValidation, "failed to build comment for writing")
		}

		commentWritten, err := comment.WriteTo(file)
		if err != nil {
			return pkgerrors.WrapError(err, pkgerrors.ErrTypeIO, "failed to write package comment")
		}

		// Update header with comment information
		p.header.CommentStart = commentStart
		p.header.CommentSize = uint32(commentWritten)
	} else {
		// No comment
		p.header.CommentStart = 0
		p.header.CommentSize = 0
	}

	// No signatures in baseline implementation
	p.header.SignatureOffset = 0

	// Seek back to beginning and write updated header
	if _, err := file.Seek(0, 0); err != nil {
		return pkgerrors.WrapError(err, pkgerrors.ErrTypeIO, "failed to seek to beginning for header update")
	}

	if _, err := writePackageHeader(file, p.header); err != nil {
		return pkgerrors.WrapError(err, pkgerrors.ErrTypeIO, "failed to write updated header")
	}

	// Update package state
	p.index = index

	// Update PackageInfo file count (already updated timestamps/identity above)
	if p.Info == nil {
		p.Info = metadata.NewPackageInfo()
	}
	fileCount := 0
	for _, fe := range p.FileEntries {
		if fe == nil {
			continue
		}
		if p.SpecialFiles != nil {
			if _, ok := p.SpecialFiles[fe.Type]; ok {
				continue
			}
		}
		fileCount++
	}
	p.Info.FileCount = fileCount

	// Note: FileCount doesn't sync to header because header doesn't have a FileCount field
	// The file count is derived from the file index when reading

	// Calculate total sizes
	var totalOriginalSize, totalStoredSize uint64
	for _, fe := range p.FileEntries {
		if fe != nil {
			totalOriginalSize += fe.OriginalSize
			totalStoredSize += fe.StoredSize
		}
	}
	p.Info.FilesUncompressedSize = int64(totalOriginalSize)
	p.Info.FilesCompressedSize = int64(totalStoredSize)

	return nil
}

func (p *filePackage) syncStoredMetadataFromMemory(fe *metadata.FileEntry) {
	if fe == nil || !fe.IsDataLoaded {
		return
	}

	if fe.StoredSize == 0 {
		fe.StoredSize = uint64(len(fe.Data))
	}
	if fe.OriginalSize == 0 {
		fe.OriginalSize = fe.StoredSize
	}
	if fe.RawChecksum == 0 || fe.StoredChecksum == 0 {
		checksum := internal.CalculateCRC32(fe.Data)
		if fe.RawChecksum == 0 {
			fe.RawChecksum = checksum
		}
		if fe.StoredChecksum == 0 {
			fe.StoredChecksum = checksum
		}
	}
}

func (p *filePackage) rewriteFileEntryMeta(file *os.File, entryOffset uint64, fe *metadata.FileEntry) error {
	if _, err := file.Seek(int64(entryOffset), io.SeekStart); err != nil {
		return pkgerrors.WrapError(err, pkgerrors.ErrTypeIO, "failed to seek to file entry metadata for rewrite")
	}
	if _, err := fe.WriteMetaTo(file); err != nil {
		return pkgerrors.WrapError(err, pkgerrors.ErrTypeIO, "failed to rewrite file entry metadata")
	}
	return nil
}

func writePackageHeader(w io.Writer, header *fileformat.PackageHeader) (int64, error) {
	if header == nil {
		return 0, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "package header is nil", nil, struct{}{})
	}
	if err := binary.Write(w, binary.LittleEndian, header); err != nil {
		return 0, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write header", pkgerrors.ValidationErrorContext{
			Field:    "Header",
			Value:    nil,
			Expected: fmt.Sprintf("%d bytes", fileformat.PackageHeaderSize),
		})
	}
	return fileformat.PackageHeaderSize, nil
}

func writeFileIndexTo(w io.Writer, index *fileformat.FileIndex) (int64, error) {
	if index == nil {
		return 0, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "file index is nil", nil, struct{}{})
	}

	var totalWritten int64
	index.EntryCount = uint32(len(index.Entries))

	if err := binary.Write(w, binary.LittleEndian, index.EntryCount); err != nil {
		return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write entry count", pkgerrors.ValidationErrorContext{
			Field:    "EntryCount",
			Value:    index.EntryCount,
			Expected: "written successfully",
		})
	}
	totalWritten += 4

	if err := binary.Write(w, binary.LittleEndian, index.Reserved); err != nil {
		return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write reserved", pkgerrors.ValidationErrorContext{
			Field:    "Reserved",
			Value:    index.Reserved,
			Expected: "written successfully",
		})
	}
	totalWritten += 4

	if err := binary.Write(w, binary.LittleEndian, index.FirstEntryOffset); err != nil {
		return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write first entry offset", pkgerrors.ValidationErrorContext{
			Field:    "FirstEntryOffset",
			Value:    index.FirstEntryOffset,
			Expected: "written successfully",
		})
	}
	totalWritten += 8

	for i, entry := range index.Entries {
		if err := binary.Write(w, binary.LittleEndian, entry); err != nil {
			return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, fmt.Sprintf("failed to write entry %d", i), pkgerrors.ValidationErrorContext{
				Field:    "Entries",
				Value:    i,
				Expected: "written successfully",
			})
		}
		totalWritten += fileformat.IndexEntrySize
	}

	return totalWritten, nil
}
