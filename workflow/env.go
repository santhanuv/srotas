package workflow

import (
	"fmt"
	"maps"

	"github.com/expr-lang/expr"
)

type Env struct {
	varExprs    map[string]string
	headerExprs map[string][]string
}

func (e *Env) AppendVars(varExprList ...map[string]string) error {
	for _, v := range varExprList {
		for name, val := range v {
			if _, ok := e.varExprs[name]; ok {
				return fmt.Errorf("variable '%s' is already defined", name)
			}
			e.varExprs[name] = val
		}
	}

	return nil
}

func (e *Env) AppendHeaders(headerExprList ...map[string][]string) {
	for _, headers := range headerExprList {
		for key, val := range headers {
			e.headerExprs[key] = append(e.headerExprs[key], val...)
		}
	}
}

func (e *Env) Compile(vars map[string]any) (map[string]any, map[string][]string, error) {
	var (
		cVars    map[string]any
		cHeaders map[string][]string
	)

	if vars != nil {
		cVars := make(map[string]any, len(e.varExprs)+len(vars))
		maps.Copy(cVars, vars)
	}

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

func NewEnv(varExprs map[string]string, headerExprs map[string][]string) *Env {
	if varExprs == nil {
		varExprs = make(map[string]string)
	}

	if headerExprs == nil {
		headerExprs = make(map[string][]string)
	}

	return &Env{
		varExprs:    varExprs,
		headerExprs: headerExprs,
	}
}
