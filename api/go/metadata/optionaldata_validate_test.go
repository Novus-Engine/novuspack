package metadata

import (
	"testing"
)

// TestOptionalDataEntry_Validate tests the Validate method.
func TestOptionalDataEntry_Validate(t *testing.T) {
	tests := []struct {
		name    string
		entry   OptionalDataEntry
		wantErr bool
	}{
		{
			name: "valid optional data entry",
			entry: OptionalDataEntry{
				DataType:   OptionalDataTagsData,
				DataLength: 10,
				Data:       make([]byte, 10),
			},
			wantErr: false,
		},
		{
			name: "empty data",
			entry: OptionalDataEntry{
				DataType:   OptionalDataTagsData,
				DataLength: 0,
				Data:       []byte{},
			},
			wantErr: true,
		},
		{
			name: "data length mismatch",
			entry: OptionalDataEntry{
				DataType:   OptionalDataTagsData,
				DataLength: 10,
				Data:       make([]byte, 5),
			},
			wantErr: true,
		},
		{
			name: "nil data",
			entry: OptionalDataEntry{
				DataType:   OptionalDataTagsData,
				DataLength: 10,
				Data:       nil,
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
