package workflow

import (
	"errors"
	"testing"
)

func TestHttp_Validate(t *testing.T) {
	tests := []struct {
		name  string
		Http  Request
		err   bool
		field string
	}{
		{
			name: "'name' is not provided",
			Http: Request{
				Type:   "http",
				Url:    "url",
				Method: "GET",
			},
			err:   true,
			field: "name",
		},
		{
			name: "'url' field is not provided",
			Http: Request{
				Type:     "http",
				StepName: "name",
				Method:   "GET",
			},
			err:   true,
			field: "url",
		},
		{
			name: "'method' field is not provided",
			Http: Request{
				Type:     "http",
				StepName: "name",
				Url:      "url",
			},
			err:   true,
			field: "method",
		},
		{
			name: "All required fields are provided",
			Http: Request{
				Type:     "http",
				StepName: "name",
				Url:      "url",
				Method:   "GET",
			},
			err: false,
		},
	}

	for _, tt := range tests {
		err := tt.Http.Validate()

		if !tt.err && err != nil {
			t.Errorf("in test %q; expected no error but got %q", tt.name, err)
			continue
		}

		if tt.err {
			if err == nil {
				t.Errorf("in test %q; expected error but got none", tt.name)
				continue
			}

			var vErr *ValidationError
			if !errors.As(err, &vErr) {
				t.Errorf("in test %q; got %q; want validation error", tt.name, err)
			}
		}
	}
}
