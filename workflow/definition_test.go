package workflow

import (
	"errors"
	"testing"
)

func TestDefinition_Validate(t *testing.T) {
	tests := []struct {
		name string
		def  *Definition
		err  bool
	}{
		{
			name: "Required field is missing",
			def: &Definition{
				Version: "1.0.0",
				BaseUrl: "http://example.com",
				Timeout: 5000,
				Variables: map[string]string{
					"varName1": "varValue1",
				},
				Headers: Header{
					"Accept": []string{"application/json"},
				},
				OutputAll: true,
			},
			err: true,
		},
		{
			name: "All required fields are provided",
			def: &Definition{
				Steps: StepList{},
			},
			err: false,
		},
	}

	for _, tt := range tests {
		err := tt.def.Validate()

		if !tt.err && err != nil {
			t.Errorf("in test %q; expected no error but got %q", tt.name, err)
			continue
		}

		if tt.err {
			if err == nil {
				t.Errorf("in test %q; expected error but got none", tt.name)
				continue
			}

			var target *ValidationError
			if !errors.As(err, &target) {
				t.Errorf("in test %q; got %q; want validation error", tt.name, err)
			}
		}

	}
}
