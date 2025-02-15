package workflow

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Step represents an executable step within the configuration workflow.
// Each step must implement the Execute method, which performs the step's operation
// using the provided ExecutionContext.
type Step interface {
	Execute(execCtx *ExecutionContext) error
	Validate() error
}

// Represents a sequence of steps.
type StepList []Step

// UnmarshalYAML unmarshals the steps to specific step types.
func (s *StepList) UnmarshalYAML(value *yaml.Node) error {
	var rawSteps []struct {
		Type string
		Step rawStepNode
	}

	if err := value.Decode(&rawSteps); err != nil {
		return err
	}

	parser := newStepParser()
	steps := make(StepList, 0, len(rawSteps))
	var errors []string

	for _, rawStep := range rawSteps {
		step, err := parser.parse(rawStep.Type, rawStep.Step.Node)

		if err != nil {
			errors = append(errors, err.Error())
			continue
		}

		steps = append(steps, step)
	}

	if len(errors) > 0 {
		err := strings.Join(errors, "\n\n")
		return fmt.Errorf("%s", err)
	}

	*s = steps
	return nil
}

// rawStepNode allows to delay the parsing of actual step.
type rawStepNode struct {
	*yaml.Node
}

func (r *rawStepNode) UnmarshalYAML(value *yaml.Node) error {
	r.Node = value
	return nil
}
