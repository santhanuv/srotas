package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/expr-lang/expr"
	"github.com/santhanuv/srotas/internal/http"
	"github.com/santhanuv/srotas/internal/log"
	"github.com/santhanuv/srotas/internal/store"
	"github.com/santhanuv/srotas/workflow"
)

// ConfigRunner holds the settings used to run a configuration.
type ConfigRunner struct {
	CfgPath        string              // Path to the YAML config file.
	Debug          bool                // Debug mode is enabled if true, which outputs detailed logs.
	sVarExpr       map[string]string   // Static variable expressions.
	gHeaderExpr    map[string][]string // Global header expressions.
	InputVars      map[string]any      // Compiled input variables.
	DefHttpTimeout uint                // Default timeout (in ms) for HTTP request.
}

// Run runs the configuration.
func (cr ConfigRunner) Run(logger *log.Logger, in *os.File, out io.Writer) error {
	// Parsing
	logger.Debug("parsing configuration...")

	def, err := workflow.ParseConfig(cr.CfgPath, logger)
	if err != nil {
		return fmt.Errorf("error on parsing config: %v", err)
	}

	// Context Initialization
	logger.Debug("initializing context...")

	if err := cr.AddVars(def.Variables); err != nil {
		return fmt.Errorf("error initializing variable: %v", err)
	}

	if err := cr.AddHeaders(def.Headers); err != nil {
		return fmt.Errorf("error initializing header: %v", err)
	}

	variables, headers, err := cr.Compile()

	if err != nil {
		return fmt.Errorf("failed to initialize config for execution: %v", err)
	}

	inputVars, err := parseInput(in)

	if err != nil {
		return fmt.Errorf("failed to parse input: %v", err)
	}

	for name := range inputVars {
		if _, ok := cr.sVarExpr[name]; ok {
			return fmt.Errorf("input variable '%s' is already defined", name)
		}
	}

	cr.InputVars = inputVars

	var s *store.Store = store.NewStore(cr.InputVars)

	if variables != nil {
		s.Add(variables)
	}

	httpClient := http.NewClient(cr.DefHttpTimeout)

	execCtx, err := workflow.NewExecutionContext(
		workflow.WithHttpClient(httpClient),
		workflow.WithGlobalOptions(def.BaseUrl, headers),
		workflow.WithLogger(logger),
		workflow.WithStore(s))

	if err != nil {
		return fmt.Errorf("failed to initialize config for execution: %v", err)
	}

	// Execution
	logger.Debug("executing configuration...")

	err = workflow.Execute(def, execCtx)
	if err != nil {
		return fmt.Errorf("failed to execute config: %v", err)
	}

	// Output updated variables
	if def.OutputAll || def.Output != nil {
		logger.Debug("output is being send to stdout")
		outJson, err := compileOutput(def.Output, execCtx.Variables(), def.OutputAll)

		if err != nil {
			return fmt.Errorf("failed to encode output as json: %v", err)
		}

		_, err = out.Write(outJson)

		if err != nil {
			return fmt.Errorf("failed to write output: %v", err)
		}

		_, err = out.Write([]byte("\n"))

		if err != nil {
			return fmt.Errorf("failed to write output: %v", err)
		}
	}

	logger.Debug("config executed successfully.")

	return nil
}

// AddVars merges the given variables into the [ConfigRunner].
// Each key-value pair represents a variable name and its corresponding expr expression.
// Returns an error if a variable with the same name already exists.
func (cr ConfigRunner) AddVars(exprs ...map[string]string) error {
	for _, v := range exprs {
		for name, val := range v {
			if _, ok := cr.sVarExpr[name]; ok {
				return fmt.Errorf("variable '%s' is already defined", name)
			}
			cr.sVarExpr[name] = val
		}
	}

	return nil
}

// AddHeaders merges the given headers into the [ConfigRunner].
// Each key-value pair represents a header name and its corresponding expr expressions.
// Headers with the same name will have their values appended.
func (cr ConfigRunner) AddHeaders(exprs ...map[string][]string) error {
	for _, headers := range exprs {
		for key, val := range headers {
			if _, ok := cr.gHeaderExpr[key]; ok {
				return fmt.Errorf("header '%s' is already defined", key)
			}

			cr.gHeaderExpr[key] = val
		}
	}

	return nil
}

// Compile evaluates the variable and header expressions using the provided
// vars as the evaluation environment.
// Returns the evaluated variables and headers, with expressions resolved to their final values.
func (cr ConfigRunner) Compile() (map[string]any, map[string][]string, error) {
	var (
		cVars    map[string]any
		cHeaders map[string][]string
	)

	if cr.sVarExpr != nil {
		cVars = make(map[string]any, len(cr.sVarExpr))
		for vn, ve := range cr.sVarExpr {
			val, err := expr.Eval(ve, cr.InputVars)

			if err != nil {
				e := fmt.Errorf("variable '%s': %v", vn, err)
				return nil, nil, e
			}

			if _, ok := cVars[vn]; ok {
				return nil, nil, fmt.Errorf("variable '%s' is alread defined", vn)
			}

			cVars[vn] = val
		}
	}

	if cr.gHeaderExpr != nil {
		cHeaders = make(map[string][]string, len(cr.gHeaderExpr))
		for key, exprList := range cr.gHeaderExpr {
			for _, e := range exprList {
				v, err := expr.Eval(e, cVars)

				if err != nil {
					e := fmt.Errorf("header '%s': %v", key, err)
					return nil, nil, e
				}

				val, ok := v.(string)

				if !ok {
					err := fmt.Errorf("header '%s' should be a string: cannot compile %s", key, e)
					return nil, nil, err
				}

				cHeaders[key] = append(cHeaders[key], val)
			}
		}
	}

	return cVars, cHeaders, nil
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

	output := struct {
		Variables map[string]any
	}{
		Variables: oVars,
	}

	outJson, err := json.MarshalIndent(output, "", " ")

	if err != nil {
		return nil, err
	}

	return outJson, nil
}

func NewConfigRunner() *ConfigRunner {
	return &ConfigRunner{
		sVarExpr:       map[string]string{},
		gHeaderExpr:    map[string][]string{},
		InputVars:      map[string]any{},
		DefHttpTimeout: 15000,
	}
}
