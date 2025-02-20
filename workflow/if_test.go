package workflow

import (
	"errors"
	"testing"
)

func TestIf_Validate(t *testing.T) {
	tests := []struct {
		name string
		If   If
		err  bool
	}{
		{
			name: "'name' is not provided",
			If: If{
				Type:      "if",
				Condition: "condition",
				Then:      StepList{},
			},
			err: true,
		},
		{
			name: "'condition' field is not provided",
			If: If{
				Type:     "if",
				StepName: "name",
				Then:     StepList{},
			},
			err: true,
		},
		{
			name: "'then' field is not provided",
			If: If{
				Type:      "if",
				StepName:  "name",
				Condition: "condition",
			},
			err: true,
		},
		{
			name: "All required fields are provided",
			If: If{
				Type:      "if",
				StepName:  "name",
				Condition: "condition",
				Then:      StepList{},
			},
			err: false,
		},
	}

	for _, tt := range tests {
		err := tt.If.Validate()

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
