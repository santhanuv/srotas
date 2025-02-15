package workflow

import (
	"fmt"

	"github.com/expr-lang/expr"
)

type PreExecEnv struct {
	varExprs    map[string]string
	headerExprs map[string][]string
}

// AddVars merges the given variables into the PreExecutionEnv.
// Each key-value pair represents a variable name and its corresponding expr expression.
// Returns an error if a variable with the same name already exists.
func (e *PreExecEnv) AddVars(exprs ...map[string]string) error {
	for _, v := range exprs {
		for name, val := range v {
			if _, ok := e.varExprs[name]; ok {
				return fmt.Errorf("variable '%s' is already defined", name)
			}
			e.varExprs[name] = val
		}
	}

	return nil
}

// AddHeaders merges the given headers into the PreExecutionEnv.
// Each key-value pair represents a header name and its corresponding expr expressions.
// Headers with the same name will have their values appended.
func (e *PreExecEnv) AddHeaders(exprs ...map[string][]string) error {
	for _, headers := range exprs {
		for key, val := range headers {
			if _, ok := e.headerExprs[key]; ok {
				return fmt.Errorf("header '%s' is already defined", key)
			}

			e.headerExprs[key] = val
		}
	}

	return nil
}

// Compile evaluates the variable and header expressions using the provided vars as the evaluation environment.
// Returns the evaluated variables and headers, with expressions resolved to their final values.
func (e *PreExecEnv) Compile(vars map[string]any) (map[string]any, map[string][]string, error) {
	var (
		cVars    map[string]any
		cHeaders map[string][]string
	)

	if e.varExprs != nil {
		cVars = make(map[string]any, len(e.varExprs))
		for vn, ve := range e.varExprs {
			val, err := expr.Eval(ve, vars)

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

	if e.headerExprs != nil {
		cHeaders = make(map[string][]string, len(e.headerExprs))
		for key, exprList := range e.headerExprs {
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

// NewPreExecEnv initializes and returns a new [PreExecEnv] with the given variable and header expressions.
func NewPreExecEnv(vexprs map[string]string, hexprs map[string][]string) *PreExecEnv {
	if vexprs == nil {
		vexprs = make(map[string]string)
	}

	if hexprs == nil {
		hexprs = make(map[string][]string)
	}

	return &PreExecEnv{
		varExprs:    vexprs,
		headerExprs: hexprs,
	}
}
