package metadata

import (
	"testing"

	"github.com/novus-engine/novuspack/api/go/fileformat"
)

// TestHashEntry_Validate tests the Validate method.
func TestHashEntry_Validate(t *testing.T) {
	tests := []struct {
		name    string
		entry   HashEntry
		wantErr bool
	}{
		{
			name: "valid hash entry",
			entry: HashEntry{
				HashType:    fileformat.HashTypeSHA256,
				HashPurpose: fileformat.HashPurposeContentVerification,
				HashLength:  32,
				HashData:    make([]byte, 32),
			},
			wantErr: false,
		},
		{
			name: "empty hash data",
			entry: HashEntry{
				HashType:    fileformat.HashTypeSHA256,
				HashPurpose: fileformat.HashPurposeContentVerification,
				HashLength:  0,
				HashData:    []byte{},
			},
			wantErr: true,
		},
		{
			name: "hash length mismatch",
			entry: HashEntry{
				HashType:    fileformat.HashTypeSHA256,
				HashPurpose: fileformat.HashPurposeContentVerification,
				HashLength:  32,
				HashData:    make([]byte, 16),
			},
			wantErr: true,
		},
		{
			name: "nil hash data",
			entry: HashEntry{
				HashType:    fileformat.HashTypeSHA256,
				HashPurpose: fileformat.HashPurposeContentVerification,
				HashLength:  32,
				HashData:    nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.entry.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
