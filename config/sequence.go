package config

import (
	"github.com/santhanuv/srotas/config/step"
)

type Sequence struct {
	Name        string
	Description string
	Variables   map[string]any
	Steps       step.StepList
}
