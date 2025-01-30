package workflow

import (
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

type If struct {
	Type       string
	Name       string
	Condition  string
	cCondition *vm.Program
	Then       StepList
	Else       StepList
}

func (i *If) Execute(context *executionContext) error {
	variables := context.store.ToMap()
	if i.cCondition == nil {

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

	if executionSteps == nil {
		context.logger.Info("Skipping conditional %s", i.Name)
	}

	for _, step := range executionSteps {
		err := step.Execute(context)

		if err != nil {
			return err
		}
	}

	return nil
}
