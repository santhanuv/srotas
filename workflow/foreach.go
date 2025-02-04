package workflow

import (
	"fmt"
	"reflect"

	"github.com/expr-lang/expr"
)

type ForEach struct {
	Type string
	Name string
	List string
	As   string
	Body StepList
}

func (f *ForEach) Execute(context *ExecutionContext) error {
	context.logger.Debug("Executing forEach step")
	variables := context.store.ToMap()

	if val, ok := variables[f.As]; val != nil && ok {
		return fmt.Errorf("ForEach: variable '%s' already defined.", f.As)
	}

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
		context.logger.Debug("ForEach execution for %v", item)
		context.store.Set(f.As, item)

		for _, step := range f.Body {
			err := step.Execute(context)

			if err != nil {
				return err
			}
		}
	}

	context.store.Remove(f.As)

	context.logger.Debug("Completed the execution of forEach step")
	return nil
}
