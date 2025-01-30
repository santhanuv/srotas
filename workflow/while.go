package workflow

import (
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

type While struct {
	Type       string
	Name       string
	Init       map[string]any
	Condition  string
	Update     map[string]string
	cCondition *vm.Program
	cUpdation  map[string]*vm.Program
	Body       StepList
}

func (w *While) Execute(context *executionContext) error {
	variables := context.store.ToMap()

	for key, val := range w.Init {
		if _, ok := variables[key]; ok {
			return fmt.Errorf("While step: initialization error: key '%s' already exists in context", key)
		}

		variables[key] = val
	}

	if w.cCondition == nil {
		if w.Condition == "" {
			return fmt.Errorf("While step: condition is mandatory")
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
			context.logger.Error("While step: no loop updatation is set")
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
			context.logger.Debug("Exiting loop as condition is evaluated to false")
			break
		}

		context.logger.Debug("Executing while step '%s'", w.Name)
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
			context.logger.Debug("variable '%s' updated to '%v'", key, output)
		}
	}

	return nil
}
