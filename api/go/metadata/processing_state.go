// This file defines the ProcessingState type and constants representing the
// current state of file processing. It contains the ProcessingState type
// definition and state constants (Idle, Loading, Processing, Writing, Complete,
// Error). This file should contain only the ProcessingState type definition
// and constants as specified in api_file_mgmt_file_entry.md Section 1.1.
//
// Specification: api_file_mgmt_file_entry.md: 1.1 FileEntry Structure Definition

package metadata

// ProcessingState defines the current state of file data in SourceFile.
// This tracks what processing has been applied to the data to inform Write operations
// what additional processing is needed.
//
// Specification: api_file_mgmt_addition.md: 2. AddFile Operations
type ProcessingState uint8

const (
	// ProcessingStateRaw indicates data in SourceFile is raw (unprocessed).
	// Write operations must apply compression and/or encryption as specified.
	// Specification: api_file_mgmt_addition.md: 2. AddFile Operations
	ProcessingStateRaw ProcessingState = iota

	// ProcessingStateCompressed indicates data in SourceFile is compressed but not encrypted.
	// Write operations must apply encryption if specified, otherwise write directly.
	// Specification: api_file_mgmt_addition.md: 2. AddFile Operations
	ProcessingStateCompressed

	// ProcessingStateEncrypted indicates data in SourceFile is encrypted but not compressed.
	// This state is typically not used (compression should be applied before encryption),
	// but is included for completeness.
	// Specification: api_file_mgmt_addition.md: 2. AddFile Operations
	ProcessingStateEncrypted

	// ProcessingStateCompressedAndEncrypted indicates data in SourceFile is both compressed and encrypted.
	// Write operations must write data directly without additional processing.
	// Specification: api_file_mgmt_addition.md: 2. AddFile Operations
	ProcessingStateCompressedAndEncrypted

	// Legacy workflow states (deprecated, kept for backward compatibility)
	// These should not be used in new code.

	// ProcessingStateIdle is deprecated - use ProcessingStateRaw instead
	ProcessingStateIdle ProcessingState = 100

	// ProcessingStateLoading is deprecated - workflow state not needed
	ProcessingStateLoading = 101

	// ProcessingStateProcessing is deprecated - workflow state not needed
	ProcessingStateProcessing = 102

	// ProcessingStateWriting is deprecated - workflow state not needed
	ProcessingStateWriting = 103

	// ProcessingStateComplete is deprecated - workflow state not needed
	ProcessingStateComplete = 104

	// ProcessingStateError is deprecated - use error returns instead
	ProcessingStateError = 105
)
