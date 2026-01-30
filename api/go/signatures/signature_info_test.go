// Package signatures provides signatures domain structures for the NovusPack implementation.
//
// This file contains unit tests for SignatureInfo structure.
package signatures

import (
	"testing"
)

// TestSignatureInfo tests the SignatureInfo struct.
func TestSignatureInfo(t *testing.T) {
	info := SignatureInfo{
		Index:         0,
		Type:          1,
		Size:          256,
		Offset:        1024,
		Flags:         0x01,
		Timestamp:     1234567890,
		Comment:       "Test signature",
		Algorithm:     "ML-DSA",
		SecurityLevel: 5,
		Valid:         true,
		Trusted:       true,
		Error:         "",
	}

	if info.Index != 0 {
		t.Errorf("SignatureInfo.Index = %v, want 0", info.Index)
	}
	if info.Type != 1 {
		t.Errorf("SignatureInfo.Type = %v, want 1", info.Type)
	}
	if info.Size != 256 {
		t.Errorf("SignatureInfo.Size = %v, want 256", info.Size)
	}
	if info.Offset != 1024 {
		t.Errorf("SignatureInfo.Offset = %v, want 1024", info.Offset)
	}
	if info.Flags != 0x01 {
		t.Errorf("SignatureInfo.Flags = %v, want 0x01", info.Flags)
	}
	if info.Timestamp != 1234567890 {
		t.Errorf("SignatureInfo.Timestamp = %v, want 1234567890", info.Timestamp)
	}
	if info.Comment != "Test signature" {
		t.Errorf("SignatureInfo.Comment = %v, want Test signature", info.Comment)
	}
	if info.Algorithm != "ML-DSA" {
		t.Errorf("SignatureInfo.Algorithm = %v, want ML-DSA", info.Algorithm)
	}
	if info.SecurityLevel != 5 {
		t.Errorf("SignatureInfo.SecurityLevel = %v, want 5", info.SecurityLevel)
	}
	if !info.Valid {
		t.Errorf("SignatureInfo.Valid = %v, want true", info.Valid)
	}
	if !info.Trusted {
		t.Errorf("SignatureInfo.Trusted = %v, want true", info.Trusted)
	}
	if info.Error != "" {
		t.Errorf("SignatureInfo.Error = %v, want empty string", info.Error)
	}
}

// TestSignatureInfo_ZeroValue tests the zero value of SignatureInfo.
func TestSignatureInfo_ZeroValue(t *testing.T) {
	info := SignatureInfo{}

	if info.Index != 0 {
		t.Errorf("SignatureInfo zero value Index = %v, want 0", info.Index)
	}
	if info.Type != 0 {
		t.Errorf("SignatureInfo zero value Type = %v, want 0", info.Type)
	}
	if info.Size != 0 {
		t.Errorf("SignatureInfo zero value Size = %v, want 0", info.Size)
	}
	if info.Offset != 0 {
		t.Errorf("SignatureInfo zero value Offset = %v, want 0", info.Offset)
	}
	if info.Flags != 0 {
		t.Errorf("SignatureInfo zero value Flags = %v, want 0", info.Flags)
	}
	if info.Timestamp != 0 {
		t.Errorf("SignatureInfo zero value Timestamp = %v, want 0", info.Timestamp)
	}
	if info.Comment != "" {
		t.Errorf("SignatureInfo zero value Comment = %v, want empty string", info.Comment)
	}
	if info.Algorithm != "" {
		t.Errorf("SignatureInfo zero value Algorithm = %v, want empty string", info.Algorithm)
	}
	if info.SecurityLevel != 0 {
		t.Errorf("SignatureInfo zero value SecurityLevel = %v, want 0", info.SecurityLevel)
	}
	if info.Valid {
		t.Errorf("SignatureInfo zero value Valid = %v, want false", info.Valid)
	}
	if info.Trusted {
		t.Errorf("SignatureInfo zero value Trusted = %v, want false", info.Trusted)
	}
	if info.Error != "" {
		t.Errorf("SignatureInfo zero value Error = %v, want empty string", info.Error)
	}
}

// TestSignatureInfo_InvalidSignature tests SignatureInfo with invalid signature.
func TestSignatureInfo_InvalidSignature(t *testing.T) {
	info := SignatureInfo{
		Index:         1,
		Type:          2,
		Size:          128,
		Offset:        2048,
		Flags:         0x00,
		Timestamp:     987654321,
		Comment:       "Invalid signature",
		Algorithm:     "SLH-DSA",
		SecurityLevel: 3,
		Valid:         false,
		Trusted:       false,
		Error:         "signature verification failed",
	}

	if info.Valid {
		t.Errorf("SignatureInfo.Valid = %v, want false", info.Valid)
	}
	if info.Trusted {
		t.Errorf("SignatureInfo.Trusted = %v, want false", info.Trusted)
	}
	if info.Error == "" {
		t.Errorf("SignatureInfo.Error = %v, want non-empty error message", info.Error)
	}
}

// TestSignatureInfo_MultipleSignatures tests multiple SignatureInfo instances.
func TestSignatureInfo_MultipleSignatures(t *testing.T) {
	signatures := []SignatureInfo{
		{
			Index:   0,
			Type:    1,
			Valid:   true,
			Trusted: true,
		},
		{
			Index:   1,
			Type:    2,
			Valid:   true,
			Trusted: false,
		},
		{
			Index:   2,
			Type:    3,
			Valid:   false,
			Trusted: false,
		},
	}

	if len(signatures) != 3 {
		t.Errorf("signatures length = %v, want 3", len(signatures))
	}

	for i, sig := range signatures {
		if sig.Index != i {
			t.Errorf("signatures[%d].Index = %v, want %d", i, sig.Index, i)
		}
	}
}
