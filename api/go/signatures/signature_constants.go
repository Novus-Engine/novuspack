// This file contains constants related to digital signatures including signature
// type constants (ML-DSA, SLH-DSA, PGP, X.509), signature size constants, and
// signature-related flags. This file should contain only constant definitions
// used throughout the signatures package.
//
// Specification: api_signatures.md: 2.1 Signature Type Constants

// Package novuspack provides the core NovusPack file format implementation.
//
// This package implements the NovusPack (.nvpk) file format as specified in
// package_file_format.md.
package signatures

// Signature type constants
// Specification: package_file_format.md: 8.2.1 SignatureType Field
const (
	SignatureTypeMLDSA  = 0x01 // ML-DSA (Module-Lattice Digital Signature Algorithm)
	SignatureTypeSLHDSA = 0x02 // SLH-DSA (Stateless Hash-based Digital Signature Algorithm)
	SignatureTypePGP    = 0x03 // PGP (Pretty Good Privacy)
	SignatureTypeX509   = 0x04 // X.509 Certificate-based signature
)
