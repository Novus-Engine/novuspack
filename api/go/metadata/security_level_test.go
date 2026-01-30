// Package metadata provides metadata domain structures for the NovusPack implementation.
//
// This file contains unit tests for SecurityLevel type.
package metadata

import (
	"testing"
)

// TestSecurityLevel tests SecurityLevel constants.
func TestSecurityLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    SecurityLevel
		expected SecurityLevel
	}{
		{"None", SecurityLevelNone, 0},
		{"Low", SecurityLevelLow, 1},
		{"Medium", SecurityLevelMedium, 2},
		{"High", SecurityLevelHigh, 3},
		{"Maximum", SecurityLevelMaximum, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.level != tt.expected {
				t.Errorf("SecurityLevel%s = %v, want %v", tt.name, tt.level, tt.expected)
			}
		})
	}
}

// TestSecurityLevel_Comparison tests SecurityLevel comparisons.
func TestSecurityLevel_Comparison(t *testing.T) {
	if SecurityLevelNone >= SecurityLevelLow {
		t.Errorf("SecurityLevelNone should be less than SecurityLevelLow")
	}
	if SecurityLevelLow >= SecurityLevelMedium {
		t.Errorf("SecurityLevelLow should be less than SecurityLevelMedium")
	}
	if SecurityLevelMedium >= SecurityLevelHigh {
		t.Errorf("SecurityLevelMedium should be less than SecurityLevelHigh")
	}
	if SecurityLevelHigh >= SecurityLevelMaximum {
		t.Errorf("SecurityLevelHigh should be less than SecurityLevelMaximum")
	}
}

// TestSecurityLevel_TypeConversion tests SecurityLevel type conversions.
func TestSecurityLevel_TypeConversion(t *testing.T) {
	// Test conversion to int
	levelInt := int(SecurityLevelHigh)
	if levelInt != 3 {
		t.Errorf("int(SecurityLevelHigh) = %v, want 3", levelInt)
	}

	// Test conversion from int
	levelFromInt := SecurityLevel(2)
	if levelFromInt != SecurityLevelMedium {
		t.Errorf("SecurityLevel(2) = %v, want SecurityLevelMedium", levelFromInt)
	}
}
