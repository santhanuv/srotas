package workflow

import (
	"fmt"

	"github.com/expr-lang/expr"
)

// Validator is a data structure that represents the validations for the HTTP response..
type Validator struct {
	Status_code *uint    // Expected status code for the http response.
	Asserts     []Assert // Assert expr expressions on the response body.
}

// Validate validates the http response.
// Returns an error if the validation is falied.
func (v *Validator) Validate(context *ExecutionContext, statusCode uint, rb *responseBody) error {
	if v.Status_code != nil && *v.Status_code != statusCode {
		return fmt.Errorf("status code: expected '%d' but got '%d'", *v.Status_code, statusCode)
	}

	vars := context.store.Map()
	vars["response"] = rb.body

	for _, assert := range v.Asserts {
		err := assert.Validate(vars, rb)

		if err != nil {
			return err
		}
	}

	return nil
}

// Assert represents assertions for the http response. It should be a valid expr expression.
// The http response and execution variables are available as environment for the expression evaluation.
type Assert string

// Validate runs the assertion expr expressions with the response and variables as the environment.
func (a *Assert) Validate(vars map[string]any, rb *responseBody) error {
	val, err := expr.Eval(string(*a), vars)

	if err != nil {
		return fmt.Errorf("invalid expression '%s' for assert: %v", *a, err)
	}

	isValid, ok := val.(bool)

	if !ok {
		return fmt.Errorf("evaluating expression '%s' should produce a boolean for assert", *a)
	}

	if !isValid {
		return fmt.Errorf("assertion '%s' failed", *a)
	}

	return nil
}
