package step

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Step interface {
	Execute() error
}

type StepList []Step

func (s *StepList) UnmarshalYAML(value *yaml.Node) error {
	var rawSteps []map[string]any

	if err := value.Decode(&rawSteps); err != nil {
		return err
	}

	steps := make(StepList, 0, len(rawSteps))
	for _, rawStep := range rawSteps {
		var stepType string

		if st, ok := rawStep["type"].(string); !ok {
			return NewInvalidValueType("type", "string")
		} else {
			stepType = strings.ToLower(st)
		}

		if stepType == "request" {
			step, err := parseRequestStep(rawStep)

			if err != nil {
				return err
			}

			fmt.Printf("%#v\n\n", step.Variables)
			steps = append(steps, step)
		}
	}

	*s = steps
	return nil
}

func NewInvalidValueType(field, expectedType string) error {
	return fmt.Errorf("invalid type for field %s: expected a type of %s", field, expectedType)
}

func NewMissingRequiredField(field string) error {
	return fmt.Errorf("required field %s is missing", field)
}
