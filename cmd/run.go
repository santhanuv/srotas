package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	deflog "log"
	"os"
	"path/filepath"
	"strings"

	"github.com/santhanuv/srotas/internal/log"
	"github.com/santhanuv/srotas/internal/store"
	"github.com/santhanuv/srotas/workflow"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

func init() {
	rootCmd.AddCommand(&runCommand)
	runCommand.Flags().BoolP("verbose", "v", false, "Enable verbose mode to display detailed logs about the execution of the config.")
	runCommand.Flags().String("env", "", "Environment for the execution of config. It should be json as a string or a path to json file. Supports headers and variables.")
	runCommand.Flags().Bool("output", false, "Output variables in the env")
}

var runCommand = cobra.Command{
	Use:   "run [CONFIG]",
	Short: "Run the provided configuration.",
	Long:  "Runs the provided configuration file. The configuration can be provided as a yaml file.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Config setup
		configPath := args[0]
		if configPath == "" {
			deflog.Fatal("Config: Invalid configuration file")
		}

		configPath, err := filepath.Abs(configPath)
		if err != nil {
			deflog.Fatalf("Config: %v", err)
		}

		// Verbose flag setup
		isVerbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			deflog.Fatalf("verbose flag error: %v", err)
		}

		// Logger setup
		logger := log.Logger{}
		configureLogger(&logger, isVerbose, configPath)

		// Env setup
		env := workflow.NewEnv(nil, nil)

		envFlagVal, err := cmd.Flags().GetString("env")
		if err != nil {
			logger.Fatal("Env flag error: %s", err)
		}

		err = extractEnv(env, envFlagVal)
		if err != nil {
			logger.Fatal("Env error: %s", args[0], err)
		}

		// Parsing
		logger.Debug("Parsing %s", configPath)

		flowDef, err := workflow.ParseConfig(configPath, &logger)
		if err != nil {
			logger.Fatal("Parse error: %v", err)
		}

		logger.Debug("Successfully parsed %s", configPath)

		// Context Initialization
		logger.Debug("Initializing execution context")

		env.AppendHeaders(flowDef.Headers)

		var s *store.Store
		if env.Variables != nil {
			s = store.NewStore(env.Variables)
		}

		execCtx, err := workflow.NewExecutionContext(
			workflow.WithGlobalOptions(flowDef.BaseUrl, env.Headers),
			workflow.WithLogger(&logger),
			workflow.WithStore(s))

		if err != nil {
			logger.Fatal("Execution context error: %v", err)
		}

		logger.Debug("Successfully initialized execution context")

		// Execution
		logger.Debug("Executing configuration")

		err = workflow.Execute(flowDef, execCtx)
		if err != nil {
			logger.Fatal("Execution error: %v", err)
		}

		logger.Debug("Successfully executed configuration")

		// Output updated variables
		outputRequired, err := cmd.Flags().GetBool("output")
		if err != nil {
			logger.Fatal("Output flag error: %s", err)
		}

		if outputRequired {
			logger.Debug("Output is being send to stdout")
			updatedVariables := execCtx.Variables()
			outJE := json.NewEncoder(os.Stdout)

			outJE.SetIndent("", " ")

			finalOutput := struct {
				Variables map[string]any
			}{
				Variables: updatedVariables,
			}

			if err := outJE.Encode(finalOutput); err != nil {
				logger.Fatal("failed to encode output as json: %s", err)
			}
		}
	},
}

func extractEnv(env *workflow.Env, flagEnv string) error {
	err := extractEnvFromString(env, flagEnv)

	if err != nil {
		return err
	}

	// Get input from stdin if pipe operator is used for the command
	pipeEnv, err := getEnvFromPipe()

	if err != nil {
		return err
	}

	err = extractEnvFromString(env, string(pipeEnv))

	if err != nil {
		return err
	}

	return nil
}

func getEnvFromPipe() ([]byte, error) {
	fileInfo, err := os.Stdin.Stat()

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

	return data, nil
}

func extractEnvFromString(env *workflow.Env, source string) error {
	if source == "" {
		return nil
	}

	var rawEnv struct {
		Variables map[string]any
		Headers   map[string][]string
	}

	if gjson.Valid(source) {
		err := json.Unmarshal([]byte(source), &rawEnv)
		if err != nil {
			return err
		}

		env.AppendVars(rawEnv.Variables)
		env.AppendHeaders(rawEnv.Headers)

		return nil
	}

	if strings.HasPrefix(source, "{") {
		return fmt.Errorf("env flag: invalid json.")
	}

	data, err := os.ReadFile(source)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &rawEnv)

	if err != nil {
		return err
	}

	env.AppendVars(rawEnv.Variables)
	env.AppendHeaders(rawEnv.Headers)

	return nil
}

func configureLogger(logger *log.Logger, isVerbose bool, configFileName string) {
	logger.SetForConfigFile(configFileName)
	logger.SetInfoWriter(os.Stderr)
	logger.SetErrorWriter(os.Stderr)

	if isVerbose {
		logger.SetDebugWriter(os.Stderr)
	} else {
		logger.SetDebugWriter(io.Discard)
	}
}
