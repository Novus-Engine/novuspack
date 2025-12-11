//go:build bdd

// Package core provides BDD step definitions for NovusPack core domain testing.
//
// Domain: core
// Tags: @domain:core, @phase:1
package core

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cucumber/godog"
	novuspack "github.com/novus-engine/novuspack/api/go/v1"
)

// RegisterCorePropertiesSteps registers step definitions for package property patterns (package has/is).
//
// Domain: core
// Phase: 1
// Tags: @domain:core
func RegisterCorePropertiesSteps(ctx *godog.ScenarioContext) {
	// Specific "package * is" patterns - registered BEFORE consolidated pattern for precedence
	ctx.Step(`^PackageDataVersion is examined$`, packageDataVersionIsExamined)
	ctx.Step(`^PackageDataVersion has initial value$`, packageDataVersionHasInitialValue)
	ctx.Step(`^package comment has invalid UTF-8 bytes$`, packageCommentHasInvalidUTF8Bytes)
	// Note: "package compression type is specified in header flags (bits 15-8)" is also registered in file_format_steps.go
	// This registration ensures it matches before the consolidated pattern in this file
	ctx.Step(`^package compression type is specified in header flags \(bits 15-8\)$`, packageCompressionTypeIsSpecifiedInHeaderFlagsCore)

	// Specific "package has" patterns - registered BEFORE consolidated pattern for precedence
	ctx.Step(`^package has a comment$`, packageHasAComment)
	ctx.Step(`^package has no comment$`, packageHasNoComment)
	ctx.Step(`^package has CRC calculated$`, packageHasCRCCalculated)
	ctx.Step(`^package has PackageCRC calculated$`, packageHasPackageCRCCalculated)
	ctx.Step(`^package has digital signatures$`, packageHasDigitalSignatures)
	ctx.Step(`^package has files with per-file tags$`, packageHasFilesWithPerFileTags)
	// Specific patterns with captured values
	ctx.Step(`^package has format version (\d+)$`, packageHasFormatVersion)
	ctx.Step(`^package has FileCount of (\d+)$`, packageHasFileCount)
	ctx.Step(`^package has default header values with magic number (\d+)x(\d+)E(\d+)B$`, packageHasMagicNumber)

	// Consolidated "package has" patterns - Phase 5 (after specific patterns)
	// Capture the property description as a single argument
	ctx.Step(`^package has (.+)$`, packageHasProperty)
}

// getWorld and getWorldFileFormatFromContext are defined in package_lifecycle.go
// worldFileFormatCore interface is defined in generic_patterns.go

// Specific "package * is" step implementations

// packageDataVersionIsExamined examines PackageDataVersion field
func packageDataVersionIsExamined(ctx context.Context) error {
	wf := getWorldFileFormatFromContext(ctx)
	if wf == nil {
		return godog.ErrUndefined
	}
	header := wf.GetHeader()
	if header == nil {
		// Create a default header if none exists for examination
		header = &novuspack.PackageHeader{
			Magic:              novuspack.NPKMagic,
			FormatVersion:      novuspack.FormatVersion,
			Reserved:           0,
			PackageDataVersion: 1,
			MetadataVersion:    1,
		}
		wf.SetHeader(header)
	}
	// Just verify the field exists and is accessible
	_ = header.PackageDataVersion
	return nil
}

// packageDataVersionHasInitialValue verifies PackageDataVersion has initial value
func packageDataVersionHasInitialValue(ctx context.Context) error {
	wf := getWorldFileFormatFromContext(ctx)
	if wf == nil {
		return godog.ErrUndefined
	}
	header := wf.GetHeader()
	if header == nil {
		// Create header with initial value
		header = &novuspack.PackageHeader{
			Magic:              novuspack.NPKMagic,
			FormatVersion:      novuspack.FormatVersion,
			Reserved:           0,
			PackageDataVersion: 1,
			MetadataVersion:    1,
		}
		wf.SetHeader(header)
	}
	// Verify it has the initial value (1)
	if header.PackageDataVersion != 1 {
		return fmt.Errorf("PackageDataVersion is %d, expected 1", header.PackageDataVersion)
	}
	return nil
}

// packageCommentHasInvalidUTF8Bytes creates a package comment with invalid UTF-8 bytes
func packageCommentHasInvalidUTF8Bytes(ctx context.Context) error {
	wf := getWorldFileFormatFromContext(ctx)
	if wf == nil {
		return godog.ErrUndefined
	}

	// Create a comment with invalid UTF-8 bytes
	// Invalid UTF-8: 0xFF 0xFE is not valid UTF-8
	invalidUTF8 := []byte{0xFF, 0xFE, 't', 'e', 's', 't'}
	comment := &novuspack.PackageComment{
		CommentLength: uint32(len(invalidUTF8) + 1),
		Comment:       string(invalidUTF8) + "\x00",
		Reserved:      [3]uint8{0, 0, 0},
	}
	wf.SetComment(comment)
	// This should fail validation - store error for later verification
	err := comment.Validate()
	if err != nil {
		wf.SetError(err)
		// Don't return error here - let validation step check it
	}

	// Update header
	header := wf.GetHeader()
	if header == nil {
		header = &novuspack.PackageHeader{
			Magic:         novuspack.NPKMagic,
			FormatVersion: novuspack.FormatVersion,
			Reserved:      0,
		}
		wf.SetHeader(header)
	}
	header.SetFeature(novuspack.FlagHasPackageComment)
	header.CommentSize = uint32(len(invalidUTF8))

	return nil
}

// packageCompressionTypeIsSpecifiedInHeaderFlagsCore handles compression type in header flags
func packageCompressionTypeIsSpecifiedInHeaderFlagsCore(ctx context.Context) error {
	wf := getWorldFileFormatFromContext(ctx)
	if wf == nil {
		return godog.ErrUndefined
	}
	header := wf.GetHeader()
	if header == nil {
		header = &novuspack.PackageHeader{
			Magic:         novuspack.NPKMagic,
			FormatVersion: novuspack.FormatVersion,
			Reserved:      0,
		}
		wf.SetHeader(header)
	}
	// Set compression type in bits 15-8 of Flags
	// Compression type 1 (LZ4) as example
	compressionType := uint32(1)
	header.Flags = (header.Flags & 0xFF00FF) | (compressionType << 8)
	return nil
}

// Specific "package has" step implementations

// packageHasAComment creates a package with a comment
func packageHasAComment(ctx context.Context) error {
	world := getWorld(ctx)
	if world == nil {
		return godog.ErrUndefined
	}
	// Get world with file format methods
	wf := getWorldFileFormatFromContext(ctx)
	if wf == nil {
		return godog.ErrUndefined
	}

	// Create a valid package comment using SetComment
	commentText := "test comment"
	comment := novuspack.NewPackageComment()
	if err := comment.SetComment(commentText); err != nil {
		return fmt.Errorf("failed to set comment: %w", err)
	}

	// Validate comment
	if err := comment.Validate(); err != nil {
		return fmt.Errorf("invalid comment: %w", err)
	}

	// Store in world
	wf.SetComment(comment)

	// Update header if it exists, or create one
	header := wf.GetHeader()
	if header == nil {
		header = &novuspack.PackageHeader{
			Magic:         novuspack.NPKMagic,
			FormatVersion: novuspack.FormatVersion,
			Reserved:      0,
		}
		wf.SetHeader(header)
	}
	header.SetFeature(novuspack.FlagHasPackageComment)
	header.CommentSize = comment.CommentLength

	return nil
}

// packageHasNoComment creates a package without a comment
func packageHasNoComment(ctx context.Context) error {
	wf := getWorldFileFormatFromContext(ctx)
	if wf == nil {
		return godog.ErrUndefined
	}

	// Ensure no comment is set
	wf.SetComment(nil)

	// Update header
	header := wf.GetHeader()
	if header == nil {
		header = &novuspack.PackageHeader{
			Magic:         novuspack.NPKMagic,
			FormatVersion: novuspack.FormatVersion,
			Reserved:      0,
		}
		wf.SetHeader(header)
	}
	header.ClearFeature(novuspack.FlagHasPackageComment)
	header.CommentSize = 0
	header.CommentStart = 0

	return nil
}

// packageHasCRCCalculated indicates package has CRC calculated
func packageHasCRCCalculated(ctx context.Context) error {
	wf := getWorldFileFormatFromContext(ctx)
	if wf == nil {
		return godog.ErrUndefined
	}

	header := wf.GetHeader()
	if header == nil {
		header = &novuspack.PackageHeader{
			Magic:         novuspack.NPKMagic,
			FormatVersion: novuspack.FormatVersion,
			Reserved:      0,
		}
		wf.SetHeader(header)
	}
	// Set a non-zero CRC to indicate it's calculated
	header.PackageCRC = 0x12345678 // Example CRC value
	return nil
}

// packageHasPackageCRCCalculated same as packageHasCRCCalculated
func packageHasPackageCRCCalculated(ctx context.Context) error {
	return packageHasCRCCalculated(ctx)
}

// packageHasDigitalSignatures indicates package has digital signatures
func packageHasDigitalSignatures(ctx context.Context) error {
	wf := getWorldFileFormatFromContext(ctx)
	if wf == nil {
		return godog.ErrUndefined
	}

	header := wf.GetHeader()
	if header == nil {
		header = &novuspack.PackageHeader{
			Magic:         novuspack.NPKMagic,
			FormatVersion: novuspack.FormatVersion,
			Reserved:      0,
		}
		wf.SetHeader(header)
	}
	header.SetFeature(novuspack.FlagHasSignatures)
	header.SignatureOffset = 1000 // Example offset
	return nil
}

// packageHasFilesWithPerFileTags indicates package has files with per-file tags
func packageHasFilesWithPerFileTags(ctx context.Context) error {
	wf := getWorldFileFormatFromContext(ctx)
	if wf == nil {
		return godog.ErrUndefined
	}

	header := wf.GetHeader()
	if header == nil {
		header = &novuspack.PackageHeader{
			Magic:         novuspack.NPKMagic,
			FormatVersion: novuspack.FormatVersion,
			Reserved:      0,
		}
		wf.SetHeader(header)
	}
	header.SetFeature(novuspack.FlagHasPerFileTags)
	return nil
}

// Specific "package has" handlers with captured values

// packageHasFormatVersion handles "package has format version X"
func packageHasFormatVersion(ctx context.Context, formatVersion string) error {
	wf := getWorldFileFormatFromContext(ctx)
	if wf == nil {
		return godog.ErrUndefined
	}
	header := wf.GetHeader()
	if header == nil {
		header = &novuspack.PackageHeader{
			Magic:         novuspack.NPKMagic,
			FormatVersion: novuspack.FormatVersion,
			Reserved:      0,
		}
		wf.SetHeader(header)
	}
	version, err := strconv.ParseUint(formatVersion, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid format version: %s", formatVersion)
	}
	header.FormatVersion = uint32(version)
	return nil
}

// packageHasFileCount handles "package has FileCount of X"
func packageHasFileCount(ctx context.Context, fileCount string) error {
	wf := getWorldFileFormatFromContext(ctx)
	if wf == nil {
		return godog.ErrUndefined
	}
	header := wf.GetHeader()
	if header == nil {
		header = &novuspack.PackageHeader{
			Magic:         novuspack.NPKMagic,
			FormatVersion: novuspack.FormatVersion,
			Reserved:      0,
		}
		wf.SetHeader(header)
	}
	count, err := strconv.ParseUint(fileCount, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid file count: %s", fileCount)
	}
	// FileCount is stored in the index, not directly in header
	// For BDD testing purposes, we'll just validate the count is parseable
	_ = count
	// TODO: Store file count in world for later validation
	return nil
}

// packageHasMagicNumber handles "package has default header values with magic number XxYEZ"
func packageHasMagicNumber(ctx context.Context, magic1, magic2, magic3 string) error {
	wf := getWorldFileFormatFromContext(ctx)
	if wf == nil {
		return godog.ErrUndefined
	}
	header := wf.GetHeader()
	if header == nil {
		header = &novuspack.PackageHeader{
			Magic:         novuspack.NPKMagic,
			FormatVersion: novuspack.FormatVersion,
			Reserved:      0,
		}
		wf.SetHeader(header)
	}
	// Parse magic number components (format: XxYEZ where X, Y, Z are digits)
	// This is a simplified parser - adjust based on actual format
	m1, err1 := strconv.ParseUint(magic1, 10, 32)
	m2, err2 := strconv.ParseUint(magic2, 10, 32)
	m3, err3 := strconv.ParseUint(magic3, 10, 32)
	if err1 != nil || err2 != nil || err3 != nil {
		return fmt.Errorf("invalid magic number format: %sx%sE%s", magic1, magic2, magic3)
	}
	// Construct magic number from components (adjust based on actual format)
	_ = m1
	_ = m2
	_ = m3
	// TODO: Set magic number based on parsed components
	return nil
}

// Consolidated "package has" pattern implementation - Phase 5

// packageHasProperty handles "package has ..." patterns (simple cases without captured values)
func packageHasProperty(ctx context.Context, property string) error {
	// Route to specific handlers for known properties
	switch property {
	case "a comment":
		return packageHasAComment(ctx)
	case "no comment":
		return packageHasNoComment(ctx)
	case "CRC calculated", "PackageCRC calculated":
		return packageHasCRCCalculated(ctx)
	case "digital signatures":
		return packageHasDigitalSignatures(ctx)
	case "files with per-file tags":
		return packageHasFilesWithPerFileTags(ctx)
	}
	// TODO: Handle other simple cases
	return godog.ErrPending
}
