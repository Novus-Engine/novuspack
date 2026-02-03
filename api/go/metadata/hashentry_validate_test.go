package metadata

import (
	"testing"

	"github.com/novus-engine/novuspack/api/go/fileformat"
)

type hashEntryValidatable struct {
	entry HashEntry
}

func (h hashEntryValidatable) Validate() error {
	return h.entry.validate()
}

// TestHashEntry_Validate tests the Validate method.
func TestHashEntry_Validate(t *testing.T) {
	tests := []validateCase{
		{
			name: "valid hash entry",
			subject: hashEntryValidatable{entry: HashEntry{
				HashType:    fileformat.HashTypeSHA256,
				HashPurpose: fileformat.HashPurposeContentVerification,
				HashLength:  32,
				HashData:    make([]byte, 32),
			}},
			wantErr: false,
		},
		{
			name: "empty hash data",
			subject: hashEntryValidatable{entry: HashEntry{
				HashType:    fileformat.HashTypeSHA256,
				HashPurpose: fileformat.HashPurposeContentVerification,
				HashLength:  0,
				HashData:    []byte{},
			}},
			wantErr: true,
		},
		{
			name: "hash length mismatch",
			subject: hashEntryValidatable{entry: HashEntry{
				HashType:    fileformat.HashTypeSHA256,
				HashPurpose: fileformat.HashPurposeContentVerification,
				HashLength:  32,
				HashData:    make([]byte, 16),
			}},
			wantErr: true,
		},
		{
			name: "nil hash data",
			subject: hashEntryValidatable{entry: HashEntry{
				HashType:    fileformat.HashTypeSHA256,
				HashPurpose: fileformat.HashPurposeContentVerification,
				HashLength:  32,
				HashData:    nil,
			}},
			wantErr: true,
		},
	}

	runValidateTable(t, tests)
}
