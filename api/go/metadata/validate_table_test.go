// Shared table runner for Validate() tests across metadata types.

package metadata

import "testing"

type validatable interface {
	Validate() error
}

type validateCase struct {
	name    string
	subject validatable
	wantErr bool
}

func runValidateTable(t *testing.T, tests []validateCase) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.subject.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
