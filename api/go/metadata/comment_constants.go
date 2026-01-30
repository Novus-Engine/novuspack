// This file contains constants related to package comments including maximum
// comment length and validation constants. This file should contain only
// constant definitions used for comment validation and processing.
//
// Specification: api_metadata.md: 1. Comment Management

// Package novuspack provides metadata domain structures for the NovusPack implementation.
//
// This package contains structures and constants related to package metadata
// as specified in docs/tech_specs/package_file_format.md and api_metadata.md.
package metadata

// MaxCommentLength is the maximum allowed comment length (1MB - 1)
// Specification: package_file_format.md: 6.1.1 NewFileIndex Function
const MaxCommentLength = 1048575
