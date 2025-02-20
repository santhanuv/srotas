package workflow

import (
	"errors"
	"testing"
)

func TestWhile_Validate(t *testing.T) {
	tests := []struct {
		name  string
		While While
		err   bool
	}{
		{
			name: "'name' is not provided",
			While: While{
				Type:      "while",
				Condition: "condition",
				Body:      StepList{},
			},
			err: true,
		},
		{
			name: "'condition' field is not provided",
			While: While{
				Type:     "while",
				StepName: "name",
				Body:     StepList{},
			},
			err: true,
		},
		{
			name: "'body' field is not provided",
			While: While{
				Type:      "while",
				StepName:  "name",
				Condition: "condition",
			},
			err: true,
		},
		{
			name: "All required fields are provided",
			While: While{
				Type:      "while",
				StepName:  "name",
				Condition: "condition",
				Body:      StepList{},
			},
			err: false,
		},
	}

	for _, tt := range tests {
		err := tt.While.Validate()

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
