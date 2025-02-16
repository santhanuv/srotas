package workflow

import (
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

// If represents a conditional step that executes Then steps when Condition evaluates to true;
// otherwise, it executes Else steps if provided.
type If struct {
	Type       string      // The type of the step.
	StepName   string      `yaml:"name"` // Identifier for the step.
	Condition  string      // Expression that determines which branch to execute.
	cCondition *vm.Program // Precompiled condition expression.
	Then       StepList    // Steps to execute if Condition is true.
	Else       StepList    // Steps to execute if Condition is false.
}

// Validate checks the fields of the [If] step and returns a list of validation errors, if any.
func (i *If) Validate() error {
	vErr := ValidationError{}

	if i.StepName == "" {
		vErr.Add(RequiredFieldError{Field: "name"})
	}

	if i.Condition == "" {
		vErr.Add(RequiredFieldError{Field: "condition"})
	}

	if i.Then == nil {
		vErr.Add(RequiredFieldError{Field: "then"})
	}

	if vErr.HasError() {
		return fmt.Errorf("if step: %w", &vErr)
	}

	return nil
}

func (i *If) Name() string {
	return i.StepName
}

// Execute executes the step with the specified context.
func (i *If) Execute(context *ExecutionContext) error {
	variables := context.store.Map()

	if i.cCondition == nil {
		if i.Condition == "" {
			return fmt.Errorf("if step '%s': condition is mandatory", i.StepName)
		}

		program, err := expr.Compile(i.Condition, expr.Env(variables), expr.AsBool())

		if err != nil {
			return err
		}

		i.cCondition = program
	}

	output, err := expr.Run(i.cCondition, variables)

	if err != nil {
		return err
	}

	ok := output.(bool)

	var executionSteps StepList = nil

	if ok {
		executionSteps = i.Then
	} else {
		executionSteps = i.Else
	}

	for _, step := range executionSteps {
		err := step.Execute(context)

		if err != nil {
			return err
		}
	}

	context.logger.Debug("successfully completed the execution of if step '%s'.", i.StepName)

	return nil
}
