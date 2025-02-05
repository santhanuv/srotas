package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	deflog "log"
	"os"
	"path/filepath"
	"strings"

	"github.com/expr-lang/expr"
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

type output struct {
	Variables map[string]any
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
		isVerbose, err := cmd.Flags().GetBool("debug")
		if err != nil {
			deflog.Fatalf("verbose flag error: %v", err)
		}

		// Logger setup
		logger := log.Logger{}
		configureLogger(&logger, isVerbose, configPath)

		// Piped Input
		pVars, err := parsePipedOutput()

		if err != nil {
			logger.Fatal("parsing input error: %s", err)
		}

		// Env setup
		env := workflow.NewEnv(nil, nil)

		envFlagVal, err := cmd.Flags().GetString("env")
		if err != nil {
			logger.Fatal("Env flag error: %s", err)
		}

		err = extractEnvFromString(env, envFlagVal)
		if err != nil {
			logger.Fatal("Env error: %s", args[0], err)
		}

		// Header flags
		rfh, err := cmd.Flags().GetStringArray("header")
		if err != nil {
			logger.Fatal("Header flag error: %s", err)
		}

		fheaders, err := parseStringHeader(rfh)

		if err != nil {
			logger.Fatal("header flag error: %v", err)
		}

		// Vars flag
		rv, err := cmd.Flags().GetStringArray("var")
		if err != nil {
			logger.Fatal("Var flag error: %s", err)
		}

		fVars := parseStringVars(rv)

		// Parsing
		logger.Debug("Parsing %s", configPath)

		flowDef, err := workflow.ParseConfig(configPath, &logger)
		if err != nil {
			logger.Fatal("Parse error: %v", err)
		}

		logger.Debug("Successfully parsed %s", configPath)

		// Context Initialization
		logger.Debug("Initializing execution context")

		if err := env.AppendVars(flowDef.Variables); err != nil {
			logger.Fatal("config variable error: %v", err)
		}

		if err := env.AppendVars(fVars); err != nil {
			logger.Fatal("flag variable error: %v", err)
		}

		env.AppendHeaders(flowDef.Headers)
		env.AppendHeaders(fheaders)

		variables, headers, err := env.Compile(pVars)

		if err != nil {
			logger.Fatal("env compile error: %s", err)
		}

		var s *store.Store

		if variables != nil {
			s = store.NewStore(variables)

			if pVars != nil {
				s.Add(pVars)
			}
		}

		execCtx, err := workflow.NewExecutionContext(
			workflow.WithGlobalOptions(flowDef.BaseUrl, headers),
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
		if flowDef.OutputAll || flowDef.Output != nil {
			logger.Debug("Output is being send to stdout")
			outJson, err := writeOutput(flowDef.Output, execCtx, flowDef.OutputAll)

			if err != nil {
				logger.Fatal("failed to encode output as json: %s", err)
			}

			_, err = os.Stdout.Write(outJson)

			if err != nil {
				logger.Fatal("output write error: %s")
			}

			logger.DebugJson(outJson, "Output:")
		}
	},
}

func parsePipedOutput() (map[string]any, error) {
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

func extractEnvFromString(env *workflow.Env, source string) error {
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

		err = env.AppendVars(rawEnv.Variables)

		if err != nil {
			return err
		}

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

	err = env.AppendVars(rawEnv.Variables)

	if err != nil {
		return err
	}

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

func writeOutput(ves map[string]string, ec *workflow.ExecutionContext, outputAll bool) ([]byte, error) {
	if ves == nil && !outputAll {
		return nil, fmt.Errorf("output error: please ensure output field exists")
	}

	vars := ec.Variables()

	var oVars map[string]any

	if outputAll {
		oVars = vars
	} else {
		oVars = make(map[string]any, len(ves))
		for vn, ve := range ves {
			val, err := expr.Eval(ve, vars)

			if err != nil {
				return nil, err
			}

			oVars[vn] = val
		}
	}

	out := output{
		Variables: oVars,
	}

	outJson, err := json.MarshalIndent(out, "", " ")

	if err != nil {
		return nil, err
	}

	return outJson, nil
}

func parseStringHeader(headers []string) (map[string][]string, error) {
	hm := map[string][]string{}
	for _, header := range headers {
		kvp := strings.Split(header, ":")

		if len(kvp) != 2 {
			return nil, fmt.Errorf("header should be in the form 'key:val'")
		}

		k, v := strings.TrimSpace(kvp[0]), strings.TrimSpace(kvp[1])
		if _, ok := hm[k]; !ok {
			hm[k] = []string{}
		}

		hm[k] = append(hm[k], v)
	}

	return hm, nil
}

func parseStringVars(vars []string) map[string]string {
	vm := map[string]string{}
	for _, variable := range vars {
		kvp := strings.Split(variable, "=")
		k, v := kvp[0], kvp[1]
		vm[k] = v
	}

	return vm
}
