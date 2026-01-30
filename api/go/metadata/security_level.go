// This file defines the SecurityLevel type and constants representing security
// levels for packages and signatures. It contains the SecurityLevel type definition
// and security level constants (None, Low, Medium, High, Maximum). This file
// should contain only the SecurityLevel type definition and constants.
//
// Specification: api_metadata.md: 0 Overview

// Package metadata provides metadata domain structures for the NovusPack implementation.
//
// This package contains structures and constants related to package metadata
// as specified in docs/tech_specs/package_file_format.md and api_metadata.md.
package metadata

// SecurityLevel represents the security level of a package.
//
// TODO: Full implementation with all security levels per specification.
//
// Specification: api_security.md: 2 SecurityStatus Structure
type SecurityLevel int

const (
	// SecurityLevelNone indicates no security features
	SecurityLevelNone SecurityLevel = iota

	// SecurityLevelLow indicates basic security features
	SecurityLevelLow

	// SecurityLevelMedium indicates moderate security features
	SecurityLevelMedium

	// SecurityLevelHigh indicates high security features
	SecurityLevelHigh

	// SecurityLevelMaximum indicates maximum security features
	SecurityLevelMaximum
)
