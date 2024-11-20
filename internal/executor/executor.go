package executor

import (
	"net/http"

	"github.com/santhanuv/srotas/config"
)

type ExecutionContext struct {
	httpClient *http.Client
}

func (e *ExecutionContext) HttpClient() *http.Client {
	return e.httpClient
}

func Execute(definition *config.Definition) error {
	steps := definition.Sequence.Steps

	context := &ExecutionContext{
		httpClient: http.DefaultClient,
	}
	for _, step := range steps {
		step.Execute(context)
	}

	return nil
}
