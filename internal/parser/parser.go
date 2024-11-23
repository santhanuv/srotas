package parser

import (
	"os"

	"github.com/santhanuv/srotas/config"
	"gopkg.in/yaml.v3"
)

type Parser struct {
	errors []error
}

func ParseConfig(path string) (*config.Definition, error) {
	cfg, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var def config.Definition
	err = yaml.Unmarshal(cfg, &def)

	if err != nil {
		return nil, err
	}

	return &def, nil
}
