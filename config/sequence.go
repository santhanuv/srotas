package config

import (
	"github.com/santhanuv/srotas/config/step"
	"github.com/santhanuv/srotas/config/step/validation"
)

type Sequence struct {
	Name        string
	Description string
	Variables   map[string]any
	Steps       step.StepList
	Validations []validation.Validation
}
