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
		case "if":
			ifStep := &If{
				Type: rawStep.Type,
			}
			if err := rawStep.Step.Decode(ifStep); err != nil {
				errors = append(errors, err.Error())
				continue
			}
			steps = append(steps, ifStep)
		case "while":
			whileStep := &While{
				Type: rawStep.Type,
			}
			if err := rawStep.Step.Decode(whileStep); err != nil {
				errors = append(errors, err.Error())
				continue
			}
			steps = append(steps, whileStep)
		case "forEach":
			foreachStep := &ForEach{
				Type: rawStep.Type,
			}
			if err := rawStep.Step.Decode(foreachStep); err != nil {
				errors = append(errors, err.Error())
				continue
			}
			steps = append(steps, foreachStep)
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
}

func (r *rawStepNode) UnmarshalYAML(value *yaml.Node) error {
	r.Node = value
	return nil
}
