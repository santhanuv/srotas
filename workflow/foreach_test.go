package workflow

import (
	"errors"
	"testing"
)

func TestForEach_Validate(t *testing.T) {
	tests := []struct {
		name    string
		forEach ForEach
		err     bool
	}{
		{
			name: "'name' is not provided",
			forEach: ForEach{
				Type: "forEach",
				As:   "as",
				List: "list",
				Body: StepList{},
			},
			err: true,
		},
		{
			name: "'as' field is not provided",
			forEach: ForEach{
				Type:     "forEach",
				StepName: "name",
				List:     "list",
				Body:     StepList{},
			},
			err: true,
		},
		{
			name: "'list' field is not provided",
			forEach: ForEach{
				Type:     "forEach",
				StepName: "name",
				As:       "as",
				Body:     StepList{},
			},
			err: true,
		},
		{
			name: "'body' is not provided",
			forEach: ForEach{
				Type:     "forEach",
				StepName: "name",
				As:       "as",
				List:     "list",
			},
			err: true,
		},
		{
			name: "All required fields are provided",
			forEach: ForEach{
				Type:     "forEach",
				StepName: "name",
				As:       "as",
				List:     "list",
				Body:     StepList{},
			},
			err: false,
		},
	}

	for _, tt := range tests {
		err := tt.forEach.Validate()

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
