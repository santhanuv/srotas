package workflow

import (
	"fmt"

	"github.com/expr-lang/expr"
)

// responseBody represents the response body obtained after executing the Request step.
type responseBody struct {
	body any `expr:"response"` // The json response from executing HTTP request.
}

// store stores the new set of variables after evaluating the variable expressions in varExprs
func (rb *responseBody) store(varExprs map[string]string, context *ExecutionContext) error {
	context.logger.Debug("Storing variables from response")

	if varExprs == nil {
		return nil
	}

	newVars := make(map[string]any, len(varExprs))

	vars := context.store.Map()

	for vn, ve := range varExprs {
		val, err := expr.Eval(ve, vars)

		if err != nil {
			return fmt.Errorf("invalid expression '%s' for variable '%s': %v", ve, vn, err)
		}

		newVars[vn] = val
	}

	context.store.Add(newVars)

	return nil
}
