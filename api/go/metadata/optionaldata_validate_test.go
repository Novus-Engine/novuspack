package metadata

import "testing"

type optionalDataValidatable struct {
	entry OptionalDataEntry
}

func (o optionalDataValidatable) Validate() error {
	return o.entry.validate()
}

// TestOptionalDataEntry_Validate tests the Validate method.
func TestOptionalDataEntry_Validate(t *testing.T) {
	tests := []validateCase{
		{
			name: "valid optional data entry",
			subject: optionalDataValidatable{entry: OptionalDataEntry{
				DataType:   OptionalDataTagsData,
				DataLength: 10,
				Data:       make([]byte, 10),
			}},
			wantErr: false,
		},
		{
			name: "empty data",
			subject: optionalDataValidatable{entry: OptionalDataEntry{
				DataType:   OptionalDataTagsData,
				DataLength: 0,
				Data:       []byte{},
			}},
			wantErr: true,
		},
		{
			name: "data length mismatch",
			subject: optionalDataValidatable{entry: OptionalDataEntry{
				DataType:   OptionalDataTagsData,
				DataLength: 10,
				Data:       make([]byte, 5),
			}},
			wantErr: true,
		},
		{
			name: "nil data",
			subject: optionalDataValidatable{entry: OptionalDataEntry{
				DataType:   OptionalDataTagsData,
				DataLength: 10,
				Data:       nil,
			}},
			wantErr: true,
		},
	}

	runValidateTable(t, tests)
}
