package workflow

import (
	"os"
	"path/filepath"

	"github.com/santhanuv/srotas/internal/log"
	"gopkg.in/yaml.v3"
)

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

	logger.Debug("Changed current working directory to %s for parsing", configDir)
	defer func() {
		if err := os.Chdir(wd); err == nil {
			logger.Debug("Changed current working directory back to %s", wd)
		} else {
			logger.Debug("Unable to reset current working directory to %s: Error: %v", wd, err)
		}
	}()

	var def Definition
	err = yaml.Unmarshal(cfg, &def)

	if err != nil {
		return nil, err
	}

	return &def, nil
}
