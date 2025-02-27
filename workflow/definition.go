package workflow

import (
	"fmt"
)

// Definition represents the configuration structure that is unmarshalled from the config file.
type Definition struct {
	Version   string            // The configuration version used for execution.
	BaseUrl   string            `yaml:"base_url"` // The base URL applied to all HTTP requests.
	Timeout   uint              // The maximum time (in ms) allowed for HTTP requests.
	Variables map[string]string // Predefined variables available during execution.
	Headers   Header            // Global headers added to all HTTP requests.
	Steps     StepList          // The sequence of steps to be executed.
	Output    map[string]string // Defines variables to be included in the output.
	// If true, all variables in ExecutionContext are included in the output.
	OutputAll bool `yaml:"output_all"`
}

func (d *Definition) Validate() error {
	if d.Steps == nil {
		vErr := ValidationError{}
		vErr.Add(RequiredFieldError{Field: "steps"})

		return fmt.Errorf("config: %w", &vErr)
	}

	return nil
}
