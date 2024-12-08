package executor

import (
	"github.com/santhanuv/srotas/config"
	"github.com/santhanuv/srotas/contract"
	"github.com/santhanuv/srotas/internal/http"
	"github.com/santhanuv/srotas/internal/store"
)

type ExecutionContext struct {
	httpClient    *http.Client
	localStore    *store.Store
	globalOptions *contract.Options
}

func Execute(definition *config.Definition) error {
	steps := definition.Sequence.Steps
	gopts := contract.Options{
		BaseUrl: definition.BaseUrl,
		Headers: definition.Headers,
	}

	context := &ExecutionContext{
		httpClient:    http.NewClient(definition.Timeout),
		localStore:    store.NewStore(nil),
		globalOptions: &gopts,
	}
	for _, step := range steps {
		err := step.Execute(context)

		if err != nil {
			return err
		}
	}

	return nil
}

func (e *ExecutionContext) HttpClient() *http.Client {
	return e.httpClient
}

func (e *ExecutionContext) Store() contract.Store {
	return e.localStore
}

func (e *ExecutionContext) GlobalOptions() *contract.Options {
	return e.globalOptions
}
