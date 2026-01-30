// This file implements file-path association operations for path metadata.
// It contains methods for associating and disassociating files with path
// metadata entries, and managing file-path relationships. This file should
// contain only file-path association operations (AssociateFileWithPath,
// DisassociateFileFromPath, UpdateFilePathAssociations, GetFilePathAssociations)
// as specified in api_metadata.md Section 8.2.
//
// Specification: api_metadata.md: 8.2 PathMetadata Management Methods

package novus_package

import (
	"context"

	"github.com/novus-engine/novuspack/api/go/internal"
	"github.com/novus-engine/novuspack/api/go/metadata"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// AssociateFileWithPath associates a file with a path metadata entry.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - filePath: File path string
//   - path: Path metadata path string
//
// Returns:
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 8.2 PathMetadata Management Methods
func (p *filePackage) AssociateFileWithPath(ctx context.Context, filePath string, path string) error {
	if err := internal.CheckContext(ctx, "AssociateFileWithPath"); err != nil {
		return err
	}

	// Find FileEntry by filePath
	foundFE, err := p.findFileEntryByPath(filePath)
	if err != nil {
		return pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeValidation,
			"file not found in AssociateFileWithPath", struct{}{})
	}

	// Find PathMetadataEntry by path
	foundPME, err := p.findPathMetadataByPath(path)
	if err != nil {
		return pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeValidation,
			"path not found in AssociateFileWithPath", struct{}{})
	}

	// Associate FileEntry with PathMetadataEntry
	if err := foundFE.AssociateWithPathMetadata(foundPME); err != nil {
		return pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeValidation,
			"failed to associate file with path in AssociateFileWithPath", struct{}{})
	}

	// Set parent path association for hierarchy traversal (always succeeds)
	p.setParentPathAssociation(foundPME)

	return nil
}

// DisassociateFileFromPath disassociates a file from a path metadata entry.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - filePath: File path string
//
// Returns:
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 8.2 PathMetadata Management Methods
func (p *filePackage) DisassociateFileFromPath(ctx context.Context, filePath string) error {
	if err := internal.CheckContext(ctx, "DisassociateFileFromPath"); err != nil {
		return err
	}

	// Find FileEntry by filePath
	foundFE, err := p.findFileEntryByPath(filePath)
	if err != nil {
		return err
	}

	// Remove from FileEntry.PathMetadataEntries map
	if foundFE.PathMetadataEntries != nil {
		delete(foundFE.PathMetadataEntries, filePath)
	}

	// Find corresponding PathMetadataEntry and remove FileEntry from its AssociatedFileEntries
	for _, pme := range p.PathMetadataEntries {
		if pme == nil || pme.Path.Path != filePath {
			continue
		}

		// Remove FileEntry from AssociatedFileEntries slice
		if pme.AssociatedFileEntries != nil {
			newAssociations := make([]*metadata.FileEntry, 0, len(pme.AssociatedFileEntries))
			for _, fe := range pme.AssociatedFileEntries {
				if fe != foundFE {
					newAssociations = append(newAssociations, fe)
				}
			}
			pme.AssociatedFileEntries = newAssociations
		}
	}

	return nil
}

// UpdateFilePathAssociations rebuilds all file-path associations.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//
// Returns:
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 8.2 PathMetadata Management Methods
func (p *filePackage) UpdateFilePathAssociations(ctx context.Context) error {
	if err := internal.CheckContext(ctx, "UpdateFilePathAssociations"); err != nil {
		return err
	}

	// Clear existing associations to make rebuild idempotent
	// Clear FileEntry.PathMetadataEntries
	for _, fe := range p.FileEntries {
		if fe != nil {
			fe.PathMetadataEntries = nil
		}
	}

	// Clear PathMetadataEntry.AssociatedFileEntries and ParentPath
	for _, pme := range p.PathMetadataEntries {
		if pme != nil {
			pme.AssociatedFileEntries = nil
			pme.ParentPath = nil
		}
	}

	// Build path metadata map for efficient lookup
	pathMap := make(map[string]*metadata.PathMetadataEntry)
	for _, pme := range p.PathMetadataEntries {
		if pme != nil {
			pathMap[pme.Path.Path] = pme
		}
	}

	// Rebuild associations by matching each FileEntry.Paths[].Path to PathMetadataEntry.Path.Path
	for _, fe := range p.FileEntries {
		if fe == nil {
			continue
		}
		// For each path in FileEntry.Paths
		for _, pathEntry := range fe.Paths {
			// Look up corresponding PathMetadataEntry
			if pme, exists := pathMap[pathEntry.Path]; exists {
				// Associate (errors are non-fatal - skip if association fails)
				_ = fe.AssociateWithPathMetadata(pme)
			}
		}
	}

	// Rebuild ParentPath pointers after associations
	for _, pme := range p.PathMetadataEntries {
		if pme == nil {
			continue
		}
		// Set parent path (always succeeds - missing parents are valid)
		p.setParentPathAssociation(pme)
	}

	return nil
}

// GetFilePathAssociations returns all file-path associations.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//
// Returns:
//   - map[string]*PathMetadataEntry: Map of file path -> PathMetadataEntry
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 8.2 PathMetadata Management Methods
func (p *filePackage) GetFilePathAssociations(ctx context.Context) (map[string]*metadata.PathMetadataEntry, error) {
	if err := internal.CheckContext(ctx, "GetFilePathAssociations"); err != nil {
		return nil, err
	}

	// Build association map from all FileEntries
	associations := make(map[string]*metadata.PathMetadataEntry)

	for _, fe := range p.FileEntries {
		if fe.PathMetadataEntries != nil {
			for path, pme := range fe.PathMetadataEntries {
				if pme != nil {
					associations[path] = pme
				}
			}
		}
	}

	return associations, nil
}
