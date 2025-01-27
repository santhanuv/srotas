package workflow

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Step interface {
	Execute(execCtx *executionContext) error
}

type StepList []Step

func (s *StepList) UnmarshalYAML(value *yaml.Node) error {
	var rawSteps []struct {
		Type string
		Step rawStepNode
	}

	if err := value.Decode(&rawSteps); err != nil {
		return err
	}

	steps := make(StepList, 0, len(rawSteps))
	var errors []string
	for _, rawStep := range rawSteps {
		switch rawStep.Type {
		case "http":
			hs := &Request{
				Type: rawStep.Type,
			}
			if err := rawStep.Step.Decode(hs); err != nil {
				errors = append(errors, err.Error())
				continue
			}
			steps = append(steps, hs)
		default:
			return fmt.Errorf("unsupported type %s for step", rawStep.Type)
		}
	}

	if len(errors) > 0 {
		err := strings.Join(errors, "\n ")
		return fmt.Errorf("Steps:\n %s", err)
	}

	*s = steps
	return nil
}

// rawStepNode allows to delay the parsing of actual step.
type rawStepNode struct {
	*yaml.Node
	Type string
}

func (r *rawStepNode) UnmarshalYAML(value *yaml.Node) error {
	r.Node = value
	return nil
}
