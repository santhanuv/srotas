package workflow

import (
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/santhanuv/srotas/internal"
)

// While represents a loop that executes the steps in Body while the Condition evaluates to true.
// Init defines the initial variables for the loop.
// Update specifies variable expressions that are updated after each iteration.
type While struct {
	Type       string                 // The type of the step.
	Name       string                 // Identifier for the step.
	Init       map[string]any         // Initial variables for the loop.
	Condition  string                 // Expr conditional expression for the loop.
	Update     map[string]string      // Variable expressions to update after each iteration.
	cCondition *vm.Program            // Compiled condition.
	cUpdation  map[string]*vm.Program // Compiled update expressions.
	Body       StepList               // Steps to execute in each iteration.
}

// Validate checks the fields of the [While] step and returns a list of validation errors, if any.
func (w *While) Validate() error {
	vErr := internal.ValidationError{}

	if w.Name == "" {
		vErr.Add(internal.RequiredFieldError{Field: "name"})
	}

	if w.Condition == "" {
		vErr.Add(internal.RequiredFieldError{Field: "condition"})
	}

	if w.Body == nil {
		vErr.Add(internal.RequiredFieldError{Field: "body"})
	}

	if vErr.HasError() {
		return fmt.Errorf("while step: %w", &vErr)
	}

	return nil
}

// Execute executes the step with the specified context.
func (w *While) Execute(context *ExecutionContext) error {
	variables := context.store.Map()

	for key, val := range w.Init {
		if _, ok := variables[key]; ok {
			return fmt.Errorf("while step '%s': variable '%s' already defined", w.Name, key)
		}

		variables[key] = val
	}

	defer func() {
		for name := range w.Init {
			context.store.Remove(name)
		}
	}()

	if w.cCondition == nil {
		if w.Condition == "" {
			return fmt.Errorf("while step '%s': condition is mandatory", w.Name)
		}

		program, err := expr.Compile(w.Condition, expr.Env(variables), expr.AsBool())

		if err != nil {
			return err
		}

		w.cCondition = program
	}

	if w.cUpdation == nil {
		w.cUpdation = make(map[string]*vm.Program, len(w.Update))

		if w.Update == nil {
			context.logger.Error("while step '%s': no loop updatation is set", w.Name)
		}

		for key, uExpr := range w.Update {
			program, err := expr.Compile(uExpr, expr.Env(variables))

			if err != nil {
				return err
			}

			w.cUpdation[key] = program
		}
	}

	for {
		output, err := expr.Run(w.cCondition, variables)

		if err != nil {
			return err
		}

		ok := output.(bool)

		if !ok {
			break
		}

		for _, step := range w.Body {
			err := step.Execute(context)

			if err != nil {
				return err
			}
		}

		for key, uExpr := range w.cUpdation {
			output, err := expr.Run(uExpr, variables)

			if err != nil {
				return err
			}

			variables[key] = output
		}
	}

	return nil
}
