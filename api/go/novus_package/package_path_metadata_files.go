// This file implements special metadata file operations for path metadata.
// It contains methods for saving and loading path metadata from special
// metadata files (type 65001) and updating package header flags. This file
// should contain only special metadata file operations (SavePathMetadataFile,
// LoadPathMetadataFile, UpdateSpecialMetadataFlags) as specified in
// api_metadata.md Section 8.3.
//
// Specification: api_metadata.md: 8.3 Special Metadata File Management

package novus_package

import (
	"context"

	"github.com/goccy/go-yaml"
	"github.com/samber/lo"

	"github.com/novus-engine/novuspack/api/go/fileformat"
	"github.com/novus-engine/novuspack/api/go/generics"
	"github.com/novus-engine/novuspack/api/go/internal"
	"github.com/novus-engine/novuspack/api/go/metadata"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// SavePathMetadataFile saves path metadata to special metadata file (type 65001).
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//
// Returns:
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 8.3 Special Metadata File Management
func (p *filePackage) SavePathMetadataFile(ctx context.Context) error {
	if err := internal.CheckContext(ctx, "SavePathMetadataFile"); err != nil {
		return err
	}

	// Check if PathMetadataEntries is empty
	if len(p.PathMetadataEntries) == 0 {
		// Remove existing special file if it exists
		if _, exists := p.SpecialFiles[65001]; exists {
			delete(p.SpecialFiles, 65001)
			// Also remove from FileEntries
			for i, fe := range p.FileEntries {
				if fe.Type == 65001 {
					p.FileEntries = append(p.FileEntries[:i], p.FileEntries[i+1:]...)
					break
				}
			}
		}
		return nil
	}

	// Create YAML structure
	yamlData := struct {
		Paths []*metadata.PathMetadataEntry `yaml:"paths"`
	}{
		Paths: p.PathMetadataEntries,
	}

	// Marshal to YAML bytes
	yamlBytes, err := yaml.Marshal(&yamlData)
	if err != nil {
		return pkgerrors.NewPackageError(
			pkgerrors.ErrTypeIO,
			"failed to marshal path metadata to YAML",
			err,
			struct{}{},
		)
	}

	// Check if special file already exists
	specialFile, exists := p.SpecialFiles[65001]
	if exists {
		// Update existing file data
		specialFile.Data = yamlBytes
		specialFile.OriginalSize = uint64(len(yamlBytes))
		specialFile.StoredSize = uint64(len(yamlBytes))
		specialFile.IsDataLoaded = true
	} else {
		// Find next sequential FileID (max existing + 1)
		// TODO: Consider caching max FileID and updating incrementally for better performance
		// with large numbers of FileEntries
		fileIDs := lo.Map(p.FileEntries, func(fe *metadata.FileEntry, _ int) uint64 {
			return fe.FileID
		})
		var nextFileID uint64 = 1
		if len(fileIDs) > 0 {
			nextFileID = lo.Max(fileIDs) + 1
		}

		// Create new FileEntry for special file
		specialFile = metadata.NewFileEntry()
		specialFile.FileID = nextFileID // Sequential, unique FileID
		specialFile.Type = 65001
		specialFile.Paths = []generics.PathEntry{
			{PathLength: uint16(len("/__NVPK_PATH_65001__.nvpkpath")), Path: "/__NVPK_PATH_65001__.nvpkpath"},
		}
		specialFile.CompressionType = 0 // No compression (uncompressed)
		specialFile.EncryptionType = 0  // No encryption
		specialFile.OriginalSize = uint64(len(yamlBytes))
		specialFile.StoredSize = uint64(len(yamlBytes))
		specialFile.Data = yamlBytes
		specialFile.IsDataLoaded = true
		specialFile.PathCount = 1 // One path

		// Set required tags via OptionalData
		// Tags are stored as OptionalDataEntry with type OptionalDataTagsData
		// For now, just mark that tags should be set when OptionalData is fully implemented
		// TODO: Add tags via OptionalData once tag serialization is implemented

		// Add to SpecialFiles and FileEntries
		p.SpecialFiles[65001] = specialFile
		p.FileEntries = append(p.FileEntries, specialFile)
	}

	// Update package header flags
	return p.UpdateSpecialMetadataFlags(ctx)
}

// LoadPathMetadataFile loads path metadata from special metadata file (type 65001).
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//
// Returns:
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 8.3 Special Metadata File Management
func (p *filePackage) LoadPathMetadataFile(ctx context.Context) error {
	if err := internal.CheckContext(ctx, "LoadPathMetadataFile"); err != nil {
		return err
	}

	// Check if special file type 65001 exists
	specialFile, exists := p.SpecialFiles[65001]
	if !exists {
		// Path metadata is optional - return success if not present
		return nil
	}

	// Verify special file has paths
	if len(specialFile.Paths) == 0 {
		return pkgerrors.NewPackageError(
			pkgerrors.ErrTypeValidation,
			"path metadata special file has no paths",
			nil,
			pkgerrors.ValidationErrorContext{
				Field:    "Paths",
				Value:    nil,
				Expected: "at least one path",
			},
		)
	}

	// Read file data using ReadFile (automatically decompresses if LZ4)
	data, err := p.ReadFile(ctx, specialFile.Paths[0].Path)
	if err != nil {
		return pkgerrors.WrapErrorWithContext(
			err,
			pkgerrors.ErrTypeIO,
			"failed to read path metadata file",
			pkgerrors.ValidationErrorContext{
				Field: "Path",
				Value: specialFile.Paths[0].Path,
			},
		)
	}

	// Parse YAML into intermediate structure
	var yamlData struct {
		Paths []*metadata.PathMetadataEntry `yaml:"paths"`
	}
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		return pkgerrors.NewPackageError(
			pkgerrors.ErrTypeValidation,
			"failed to parse path metadata YAML",
			err,
			pkgerrors.ValidationErrorContext{
				Field:    "YAML",
				Value:    string(data),
				Expected: "valid YAML with 'paths' key",
			},
		)
	}

	// Validate each PathMetadataEntry
	for i, pme := range yamlData.Paths {
		if err := pme.Validate(); err != nil {
			return pkgerrors.WrapErrorWithContext(
				err,
				pkgerrors.ErrTypeValidation,
				"path metadata entry validation failed",
				pkgerrors.ValidationErrorContext{
					Field:    "PathMetadataEntry",
					Value:    i,
					Expected: "valid path metadata entry",
				},
			)
		}
	}

	// Assign to package
	p.PathMetadataEntries = yamlData.Paths

	return nil
}

// UpdateSpecialMetadataFlags updates package header flags based on special files.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//
// Returns:
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 8.3 Special Metadata File Management
func (p *filePackage) UpdateSpecialMetadataFlags(ctx context.Context) error {
	if err := internal.CheckContext(ctx, "UpdateSpecialMetadataFlags"); err != nil {
		return err
	}

	// Check if any special files exist
	hasSpecialFiles := len(p.SpecialFiles) > 0

	// Update bit 6 (FlagHasSpecialMetadata)
	if hasSpecialFiles {
		p.header.Flags |= fileformat.FlagHasSpecialMetadata
	} else {
		p.header.Flags &^= fileformat.FlagHasSpecialMetadata
	}

	// Check if any PathMetadataEntry has properties (tags)
	hasPerFileTags := false
	for _, pme := range p.PathMetadataEntries {
		if len(pme.Properties) > 0 {
			hasPerFileTags = true
			break
		}
	}

	// Update bit 5 (FlagHasPerFileTags)
	if hasPerFileTags {
		p.header.Flags |= fileformat.FlagHasPerFileTags
	} else {
		p.header.Flags &^= fileformat.FlagHasPerFileTags
	}

	// Check if any PathMetadataEntry has extended attributes
	hasExtendedAttrs := false
	for _, pme := range p.PathMetadataEntries {
		if len(pme.FileSystem.ExtendedAttrs) > 0 {
			hasExtendedAttrs = true
			break
		}
	}

	// Update bit 3 (FlagHasExtendedAttrs)
	if hasExtendedAttrs {
		p.header.Flags |= fileformat.FlagHasExtendedAttrs
	} else {
		p.header.Flags &^= fileformat.FlagHasExtendedAttrs
	}

	// Update PackageInfo flags to match header
	p.Info.HasMetadataFiles = hasSpecialFiles
	p.Info.HasPerFileTags = hasPerFileTags
	p.Info.HasExtendedAttrs = hasExtendedAttrs

	return nil
}
