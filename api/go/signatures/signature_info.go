// This file implements the SignatureInfo and SignatureValidationResult structures
// providing detailed signature information and validation results. It contains
// type definitions and methods for working with signature metadata. This file
// should contain all code related to signature information structures as
// specified in api_signatures.md Section 2.2.
//
// Specification: api_signatures.md: 2.2 Signature Information Structure

// Package signatures provides signatures domain structures for the NovusPack implementation.
//
// This package contains structures and constants related to digital signatures
// as specified in docs/tech_specs/package_file_format.md and api_security.md.
package signatures

// SignatureInfo represents information about a digital signature in the package.
//
// SignatureInfo provides metadata about signatures in the package, including
// validation status, trust information, and signature details. This is different
// from the Signature struct which represents the actual signature data.
//
// TODO: Full implementation per specification.
//
// Specification: api_metadata.md: 7.2 SignatureInfo Structure
type SignatureInfo struct {
	// Index is the signature index in the package
	Index int

	// Type is the signature type (ML-DSA, SLH-DSA, PGP, X.509)
	Type uint32

	// Size is the size of signature data in bytes
	Size uint32

	// Offset is the offset to signature data from start of file
	Offset uint64

	// Flags contains signature-specific flags
	Flags uint32

	// Timestamp is the Unix timestamp when signature was created
	Timestamp uint32

	// Comment is the signature comment (if any)
	Comment string

	// Algorithm is the algorithm name/description
	Algorithm string

	// SecurityLevel is the security level (1-5)
	SecurityLevel int

	// Valid indicates whether signature is valid
	Valid bool

	// Trusted indicates whether signature is trusted
	Trusted bool

	// Error contains error message if validation failed
	Error string
}
