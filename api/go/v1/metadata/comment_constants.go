// Package novuspack provides metadata domain structures for the NovusPack implementation.
//
// This package contains structures and constants related to package metadata
// as specified in docs/tech_specs/package_file_format.md and docs/tech_specs/api_metadata.md.
package metadata

// MaxCommentLength is the maximum allowed comment length (1MB - 1)
// Specification: ../../docs/tech_specs/package_file_format.md Section 6.1.1
const MaxCommentLength = 1048575
