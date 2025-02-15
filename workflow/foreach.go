package workflow

import (
	"fmt"
	"reflect"

	"github.com/expr-lang/expr"
	"github.com/santhanuv/srotas/internal"
)

// ForEach represents a loop step that executes the steps in Body for each item in List.
// The As field defines the variable name that stores each item during execution.
type ForEach struct {
	Type string   // The type of the step.
	Name string   // Identifier for the step.
	List string   // The list of items to iterate over.
	As   string   // The variable name to store the current item in each iteration.
	Body StepList // The sequence of steps executed for each item.
}

// Validate checks the fields of the [ForEach] step and returns a list of validation errors, if any.
func (f *ForEach) Validate() error {
	vErr := internal.ValidationError{}

	if f.Name == "" {
		vErr.Add(internal.RequiredFieldError{Field: "name"})
	}

	if f.List == "" {
		vErr.Add(internal.RequiredFieldError{Field: "list"})
	}

	if f.As == "" {
		vErr.Add(internal.RequiredFieldError{Field: "as"})
	}

	if f.Body == nil {
		vErr.Add(internal.RequiredFieldError{Field: "body"})
	}

	if vErr.HasError() {
		return fmt.Errorf("foreach step: %w", &vErr)
	}

	return nil
}

// Execute executes the step with the specified context.
func (f *ForEach) Execute(context *ExecutionContext) error {
	variables := context.store.Map()

	if val, ok := variables[f.As]; val != nil && ok {
		return fmt.Errorf("foreach step '%s': variable '%s' already defined.", f.Name, f.As)
	}

	defer context.store.Remove(f.As)

	program, err := expr.Compile(f.List, expr.Env(variables), expr.AsKind(reflect.Slice))

	if err != nil {
		return err
	}

	output, err := expr.Run(program, variables)

	if err != nil {
		return err
	}

	items := output.([]any)

	for _, item := range items {
		context.store.Set(f.As, item)

		for _, step := range f.Body {
			err := step.Execute(context)

			if err != nil {
				return err
			}
		}
	}

	context.logger.Debug("successfully executed foreach step '%s'", f.Name)
	return nil
}
