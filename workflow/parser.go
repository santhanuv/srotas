package workflow

import (
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
