package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/santhanuv/srotas/internal/config"
	"github.com/santhanuv/srotas/internal/log"
	"github.com/spf13/cobra"
)

// newRunCommand creates a new instance of run command.
func newRunCommand(logger *log.Logger, in *os.File, out io.Writer, cr *config.ConfigRunner) *cobra.Command {
	runCommand := &cobra.Command{
		Use:   "run [CONFIG]",
		Short: "Run the provided configuration.",
		Long:  "Runs the provided configuration file. The configuration can be provided as a yaml file.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := parseCommand(cr, cmd, args); err != nil {
				return err
			}

			logger.SetDebugMode(cr.Debug)

			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			if err := cr.Run(logger, in, out); err != nil {
				return err
			}

			return nil
		},
	}

	runCommand.Flags().BoolP("debug", "D", false, `
		Enables debug mode, providing detailed logs about the execution of the
		configuration file.`)

	runCommand.Flags().StringP("env", "E", "", `
		Loads global headers and variables from a JSON string or file. The
		JSON may contain Variables and Headers fields, where values are expressions.
		At least one of these fields must be present, and duplicate variable names
		and header names will result in an error.

		Multiple headers with the same name cannot be defined globally
		(in --header, --env, or the config file). If a duplicate header is defined,
		an error is raised.
		`)

	runCommand.Flags().StringArrayP("header", "H", nil, `
		Adds an additional global header in the format 'key:value'. Multiple headers
		can be specified by using the flag multipe times.

		Multiple headers with the same name cannot be defined globally (in --header,
		--env, or the config file). If a duplicate header is defined, an error is raised.

		The value supports expressions, allowing dynamic header generation using defined
		or command-line variables.`)

	runCommand.Flags().StringArrayP("var", "V", nil, `
		Defines a global variable in the format name=value, where the value is an expression.
		Variables must be unique; redefining an existing one results in an error.`)

	return runCommand
}

// parseCommand extracts flags and arguments from the command line into [ConfigRunner] instance.
func parseCommand(cr *config.ConfigRunner, cmd *cobra.Command, args []string) error {
	// Config setup
	configPath := args[0]
	if configPath == "" {
		return fmt.Errorf("config file is required. Please provide a valid YAML config file. Use --help for more information")
	}

	configPath, err := filepath.Abs(configPath)
	if err != nil {
		return fmt.Errorf("invalid config: %v", err)
	}

	// Verbose flag setup
	debugMode, err := cmd.Flags().GetBool("debug")
	if err != nil {
		return fmt.Errorf("invalid value for 'debug': %v", err)
	}

	// Env setup
	efv, err := cmd.Flags().GetString("env")
	if err != nil {
		return fmt.Errorf("invalid value for 'env': %v", err)
	}

	envVars, envHeaders, err := extractEnvFromString(efv)
	if err != nil {
		return fmt.Errorf("invalid value for 'env': %v", err)
	}

	// Header flags
	fhs, err := cmd.Flags().GetStringArray("header")
	if err != nil {
		return fmt.Errorf("invalid value for 'header': %v", err)
	}

	fHeaders, err := parseStringHeaders(fhs)
	if err != nil {
		return fmt.Errorf("invalid value for 'header': %v", err)
	}

	// Vars flag
	fvs, err := cmd.Flags().GetStringArray("var")
	if err != nil {
		return fmt.Errorf("invalid value for 'var': %v", err)
	}

	fVars, err := parseStringVars(fvs)
	if err != nil {
		return fmt.Errorf("invalid value for 'var': %v", err)
	}

	cr.CfgPath = configPath
	cr.Debug = debugMode

	if err := cr.AddVars(fVars); err != nil {
		return err
	}

	if err := cr.AddVars(envVars); err != nil {
		return err
	}

	if err := cr.AddHeaders(fHeaders); err != nil {
		return err
	}

	if err := cr.AddHeaders(envHeaders); err != nil {
		return err
	}

	return nil
}

// extractEnvFromString parses the source string, extracts variables and headers, and appends them to the given env.
func extractEnvFromString(source string) (map[string]string, map[string][]string, error) {
	if source == "" {
		return map[string]string{}, map[string][]string{}, nil
	}

	var rawEnv struct {
		Variables map[string]string
		Headers   map[string][]string
	}

	if json.Valid([]byte(source)) {
		if err := json.Unmarshal([]byte(source), &rawEnv); err != nil {
			return nil, nil, err
		}

		return rawEnv.Variables, rawEnv.Headers, nil
	}

	if strings.HasPrefix(source, "{") {
		return nil, nil, fmt.Errorf("invalid json: %v", source)
	}

	data, err := os.ReadFile(source)
	if err != nil {
		return nil, nil, err
	}

	err = json.Unmarshal(data, &rawEnv)
	if err != nil {
		return nil, nil, err
	}

	return rawEnv.Variables, rawEnv.Headers, nil
}

// parseStringHeaders parses a slice of headers which are "key:value" formatted strings and returns a map containing multivalued key-pairs
func parseStringHeaders(headers []string) (map[string][]string, error) {
	hm := map[string][]string{}
	for _, header := range headers {
		kvp := strings.Split(header, ":")

		if len(kvp) != 2 {
			return nil, fmt.Errorf("header must be in the format 'key:value'")
		}

		k, v := strings.TrimSpace(kvp[0]), strings.TrimSpace(kvp[1])
		if _, ok := hm[k]; !ok {
			hm[k] = []string{}
		}

		hm[k] = append(hm[k], v)
	}

	return hm, nil
}

// parseStringVars parses a slice of variable which are "key=value" formatted strings and returns a map containing key-value pairs.
func parseStringVars(vars []string) (map[string]string, error) {
	vm := map[string]string{}
	for _, variable := range vars {
		kvp := strings.Split(variable, "=")

		if len(kvp) != 2 {
			return nil, fmt.Errorf("variable must be in the format 'key=value'")
		}

		k, v := kvp[0], kvp[1]
		vm[k] = v
	}

	return vm, nil
}
