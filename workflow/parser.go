package workflow

import (
	"os"

	"gopkg.in/yaml.v3"
)

func ParseConfig(path string) (*Definition, error) {
	cfg, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var def Definition
	err = yaml.Unmarshal(cfg, &def)

	if err != nil {
		return nil, err
	}

	return &def, nil
}
