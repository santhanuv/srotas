package config

import (
	"encoding/json"
	"io"
	"os"
)

// parseInput reads and parses input from in, returning extracted variables.
func parseInput(in *os.File) (map[string]any, error) {
	fileInfo, err := in.Stat()

	if err != nil {
		return nil, err
	}

	// Checks if the input is not connected to terminal, which means it is either piped or redirected
	if fileInfo.Mode()&os.ModeCharDevice != 0 {
		return nil, nil
	}

	data, err := io.ReadAll(os.Stdin)

	if err != nil {
		return nil, err
	}

	if string(data) == "" {
		return nil, nil
	}

	var output struct {
		Variables map[string]any
	}
	err = json.Unmarshal(data, &output)

	if err != nil {
		return nil, err
	}

	return output.Variables, nil
}
