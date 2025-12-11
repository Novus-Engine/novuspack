// Package novuspack provides the core NovusPack file format implementation.
//
// This package implements the NovusPack (.npk) file format as specified in
// docs/tech_specs/package_file_format.md.
package signatures

// Signature type constants
// Specification: ../../docs/tech_specs/package_file_format.md Section 7.2.1 - SignatureType Field
const (
	SignatureTypeMLDSA  = 0x01 // ML-DSA (Module-Lattice Digital Signature Algorithm)
	SignatureTypeSLHDSA = 0x02 // SLH-DSA (Stateless Hash-based Digital Signature Algorithm)
	SignatureTypePGP    = 0x03 // PGP (Pretty Good Privacy)
	SignatureTypeX509   = 0x04 // X.509 Certificate-based signature
)
