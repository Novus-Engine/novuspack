package metadata

import (
	"testing"
)

// TestProcessingState tests ProcessingState constants.
func TestProcessingState(t *testing.T) {
	tests := []struct {
		name     string
		state    ProcessingState
		expected ProcessingState
	}{
		// Primary states (per spec)
		{"Raw", ProcessingStateRaw, 0},
		{"Compressed", ProcessingStateCompressed, 1},
		{"Encrypted", ProcessingStateEncrypted, 2},
		{"CompressedAndEncrypted", ProcessingStateCompressedAndEncrypted, 3},
		// Legacy states (deprecated, offset by 100)
		{"Idle", ProcessingStateIdle, 100},
		{"Loading", ProcessingStateLoading, 101},
		{"Processing", ProcessingStateProcessing, 102},
		{"Writing", ProcessingStateWriting, 103},
		{"Complete", ProcessingStateComplete, 104},
		{"Error", ProcessingStateError, 105},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.state != tt.expected {
				t.Errorf("ProcessingState%s = %v, want %v", tt.name, tt.state, tt.expected)
			}
		})
	}
}

// TestProcessingState_Comparison tests ProcessingState comparisons for primary states.
func TestProcessingState_Comparison(t *testing.T) {
	// Test primary states (per spec)
	if ProcessingStateRaw >= ProcessingStateCompressed {
		t.Errorf("ProcessingStateRaw should be less than ProcessingStateCompressed")
	}
	if ProcessingStateCompressed >= ProcessingStateEncrypted {
		t.Errorf("ProcessingStateCompressed should be less than ProcessingStateEncrypted")
	}
	if ProcessingStateEncrypted >= ProcessingStateCompressedAndEncrypted {
		t.Errorf("ProcessingStateEncrypted should be less than ProcessingStateCompressedAndEncrypted")
	}

	// Test legacy states (deprecated)
	if ProcessingStateIdle >= ProcessingStateLoading {
		t.Errorf("ProcessingStateIdle should be less than ProcessingStateLoading")
	}
	if ProcessingStateLoading >= ProcessingStateProcessing {
		t.Errorf("ProcessingStateLoading should be less than ProcessingStateProcessing")
	}
}

// TestProcessingState_TypeConversion tests ProcessingState type conversions.
func TestProcessingState_TypeConversion(t *testing.T) {
	// Test conversion to uint8 for primary states
	stateUint8 := uint8(ProcessingStateEncrypted)
	if stateUint8 != 2 {
		t.Errorf("uint8(ProcessingStateEncrypted) = %v, want 2", stateUint8)
	}

	// Test conversion from uint8
	stateFromUint8 := ProcessingState(1)
	if stateFromUint8 != ProcessingStateCompressed {
		t.Errorf("ProcessingState(1) = %v, want ProcessingStateCompressed", stateFromUint8)
	}
}
