package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/santhanuv/srotas/internal/http"
	"github.com/santhanuv/srotas/internal/log"
	"github.com/santhanuv/srotas/internal/store"
	"github.com/santhanuv/srotas/workflow"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&runCommand)

	runCommand.Flags().BoolP("debug", "D", false,
		"Enables debug mode, providing detailed logs about the execution of the configuration file.")

	runCommand.Flags().StringP("env", "E", "",
		"Loads global headers and variables from a JSON string or file. The JSON may contain Variables and Headers fields, where values are expressions. At least one of these fields must be present, and duplicate variable names result in an error.")

	runCommand.Flags().StringArrayP("header", "H", nil,
		"Adds an additional global header in the format 'key:value'. Multiple headers can be specified, and values for the same key are combined with those in the config file. The value supports expressions, allowing dynamic header generation using defined or command-line variables.")

	runCommand.Flags().StringArrayP("var", "V", nil,
		"Defines a global variable in the format name=value, where the value is an expression. Variables must be unique; redefining an existing one results in an error.")
}

// output represents the output of the run command
type output struct {
	// Variables stores the values of the output field from the config.
	// If output_all is set to true, it contains all variables from the execution.
	Variables map[string]any
}

// runCommandEnv contains all the parsed values of the run command.
type runCommandEnv struct {
	config    string
	debugMode bool
	env       *workflow.PreExecEnv
	pipedVars map[string]any
}

var runCommand = cobra.Command{
	Use:   "run [CONFIG]",
	Short: "Run the provided configuration.",
	Long:  "Runs the provided configuration file. The configuration can be provided as a yaml file.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.New(os.Stderr, io.Discard, os.Stderr)

		runCmdEnv, err := parseCommand(cmd, args)
		if err != nil {
			logger.Fatal("%v", err)
		}

		// Logger setup
		logger.SetConfig(runCmdEnv.config)
		if runCmdEnv.debugMode {
			logger.SetDebugOutput(os.Stderr)
			logger.SetDebugMode(true)
		}

		// Parsing
		def, err := workflow.ParseConfig(runCmdEnv.config, logger)
		if err != nil {
			logger.Fatal("", err)
		}

		logger.Debug("successfully parsed config.")

		// Context Initialization
		if err := runCmdEnv.env.AddVars(def.Variables); err != nil {
			logger.Fatal("variable error: %v", err)
		}

		runCmdEnv.env.AddHeaders(def.Headers)

		variables, headers, err := runCmdEnv.env.Compile(runCmdEnv.pipedVars)

		if err != nil {
			logger.Fatal("failed to initialize config for execution: %v", err)
		}

		var s *store.Store = store.NewStore(runCmdEnv.pipedVars)

		if variables != nil {
			s.Add(variables)
		}

		httpClient := http.NewClient(1500)

		execCtx, err := workflow.NewExecutionContext(
			workflow.WithHttpClient(httpClient),
			workflow.WithGlobalOptions(def.BaseUrl, headers),
			workflow.WithLogger(logger),
			workflow.WithStore(s))

		if err != nil {
			logger.Fatal("failed to initialize config for execution: %v", err)
		}

		logger.Debug("successfully initialized config for execution.")

		// Execution
		err = workflow.Execute(def, execCtx)
		if err != nil {
			logger.Fatal("failed to execute config: %v", err)
		}

		logger.Debug("successfully executed configuration")

		// Output updated variables
		if def.OutputAll || def.Output != nil {
			logger.Debug("output is being send to stdout")
			outJson, err := compileOutput(def.Output, execCtx.Variables(), def.OutputAll)

			if err != nil {
				logger.Fatal("failed to encode output as json: %v", err)
			}

			_, err = os.Stdout.Write(outJson)

			if err != nil {
				logger.Fatal("failed to write outpus: %v")
			}
		}

		logger.Debug("completed execution.")
	},
}

// parseCommand extracts flags and arguments from the command line and returns a runCommandEnv instance with the parsed details.
func parseCommand(cmd *cobra.Command, args []string) (*runCommandEnv, error) {
	// Config setup
	configPath := args[0]
	if configPath == "" {
		return nil, fmt.Errorf("config file is required. Please provide a valid YAML config file. Use --help for more information.")
	}

	configPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("invalid config: %v", err)
	}

	// Verbose flag setup
	debugMode, err := cmd.Flags().GetBool("debug")
	if err != nil {
		return nil, fmt.Errorf("invalid value for 'debug': %v", err)
	}

	// Piped Input from stdin
	pVars, err := parsePipedInput()

	if err != nil {
		return nil, fmt.Errorf("failed to process input: %v", err)
	}

	// Env setup
	env := workflow.NewPreExecEnv(nil, nil)

	efv, err := cmd.Flags().GetString("env")
	if err != nil {
		return nil, fmt.Errorf("invalid value for 'env': %v", err)
	}

	err = extractEnvFromString(env, efv)
	if err != nil {
		return nil, fmt.Errorf("invalid value for 'env': %v", err)
	}

	// Header flags
	fhs, err := cmd.Flags().GetStringArray("header")
	if err != nil {
		return nil, fmt.Errorf("invalid value for 'header': %v", err)
	}

	fHeaders, err := parseStringHeaders(fhs)
	if err != nil {
		return nil, fmt.Errorf("invalid value for 'header': %v", err)
	}

	// Vars flag
	fvs, err := cmd.Flags().GetStringArray("var")
	if err != nil {
		return nil, fmt.Errorf("invalid value for 'var': %v", err)
	}

	fVars, err := parseStringVars(fvs)
	if err != nil {
		return nil, fmt.Errorf("invalid value for 'var': %v", err)
	}

	env.AddVars(fVars)
	env.AddHeaders(fHeaders)

	return &runCommandEnv{
		debugMode: debugMode,
		config:    configPath,
		env:       env,
		pipedVars: pVars,
	}, nil
}

// parsePipedInput reads and parses piped stdin input, returning extracted variables.
func parsePipedInput() (map[string]any, error) {
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

	if string(data) == "" {
		return nil, nil
	}

	var output output
	err = json.Unmarshal(data, &output)

	if err != nil {
		return nil, err
	}

	return output.Variables, nil
}

// extractEnvFromString parses the source string, extracts variables and headers, and appends them to the given env.
func extractEnvFromString(env *workflow.PreExecEnv, source string) error {
	if source == "" {
		return nil
	}

	var rawEnv struct {
		Variables map[string]string
		Headers   map[string][]string
	}

	if json.Valid([]byte(source)) {
		err := json.Unmarshal([]byte(source), &rawEnv)
		if err != nil {
			return err
		}

		err = env.AddVars(rawEnv.Variables)

		if err != nil {
			return err
		}

		env.AddHeaders(rawEnv.Headers)

		return nil
	}

	if strings.HasPrefix(source, "{") {
		return fmt.Errorf("invalid json: %v", source)
	}

	data, err := os.ReadFile(source)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &rawEnv)

	if err != nil {
		return err
	}

	err = env.AddVars(rawEnv.Variables)

	if err != nil {
		return err
	}

	env.AddHeaders(rawEnv.Headers)

	return nil
}

// compileOutput evaluates expressions in out using vars as the environment and returns a JSON representation of the output.
// If outputAll is true, all variables will be included in the output.
func compileOutput(out map[string]string, vars map[string]any, outputAll bool) ([]byte, error) {
	if out == nil && !outputAll {
		return nil, fmt.Errorf("output error: please ensure output field exists")
	}

	var oVars map[string]any

	if outputAll {
		oVars = vars
	} else {
		oVars = make(map[string]any, len(out))
		for vn, ve := range out {
			val, err := expr.Eval(ve, vars)

			if err != nil {
				return nil, err
			}

			oVars[vn] = val
		}
	}

	output := output{
		Variables: oVars,
	}

	outJson, err := json.MarshalIndent(output, "", " ")

	if err != nil {
		return nil, err
	}

	return outJson, nil
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
