package workflow

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/santhanuv/srotas/internal/log"
	"gopkg.in/yaml.v3"
)

// ParseConfig reads the configuration file from the given path and returns a [Definition] representing its contents.
func ParseConfig(path string, logger *log.Logger) (*Definition, error) {
	cfg, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	configDir := filepath.Dir(path)
	err = os.Chdir(configDir)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := os.Chdir(wd); err != nil {
			logger.Error("Unable to reset current working directory to '%s': Error: %v", wd, err)
		}
	}()

	var def Definition
	err = yaml.Unmarshal(cfg, &def)

	if err != nil {
		return nil, err
	}

	return &def, nil
}

// stepParserFunc is a function type that parses a YAML node into a Step.
// It returns the parsed Step and an error if parsing fails.
type stepParserFunc func(node *yaml.Node) (Step, error)

// stepParser is a mapping of step types to their corresponding parsing functions.
// Each step type is associated with a stepParserFunc that knows how to parse that specific step.
type stepParser map[string]stepParserFunc

// parse parses the node based on the stepType and returns the corresponding step.
// It returns the parsed Step and an error if parsing fails.
func (sp stepParser) parse(stepType string, node *yaml.Node) (Step, error) {
	parser, ok := sp[stepType]

	if !ok {
		return nil, fmt.Errorf("unsupported type %s for step", stepType)
	}

	step, err := parser(node)

	if err != nil {
		return nil, err
	}

	return step, err
}

// newStepParser initializes and returns a new [stepParser].
func newStepParser() stepParser {
	return map[string]stepParserFunc{
		"http": func(node *yaml.Node) (Step, error) {
			step := &Request{
				Type: "http",
			}

			if err := parseStep(step, node); err != nil {
				return nil, err
			}

			return step, nil
		},
		"if": func(node *yaml.Node) (Step, error) {
			step := &If{
				Type: "if",
			}

			if err := parseStep(step, node); err != nil {
				return nil, err
			}

			return step, nil
		},
		"forEach": func(node *yaml.Node) (Step, error) {
			step := &ForEach{
				Type: "if",
			}

			if err := parseStep(step, node); err != nil {
				return nil, err
			}

			return step, nil
		},
		"while": func(node *yaml.Node) (Step, error) {
			step := &While{
				Type: "if",
			}

			if err := parseStep(step, node); err != nil {
				return nil, err
			}

			return step, nil
		},
	}
}

func parseStep(step Step, node *yaml.Node) error {
	if err := node.Decode(step); err != nil {
		return err
	}

	if err := step.Validate(); err != nil {
		return err
	}

	return nil
}
