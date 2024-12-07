package executor

import (
	"github.com/santhanuv/srotas/config"
	"github.com/santhanuv/srotas/contract"
	"github.com/santhanuv/srotas/internal/http"
	"github.com/santhanuv/srotas/internal/store"
)

type ExecutionContext struct {
	httpClient *http.Client
	localStore *store.Store
	baseUrl    string
}

func (e *ExecutionContext) HttpClient() *http.Client {
	return e.httpClient
}

func (e *ExecutionContext) Store() contract.Store {
	return e.localStore
}

func (e *ExecutionContext) BaseUrl() string {
	return e.baseUrl
}

func Execute(definition *config.Definition) error {
	steps := definition.Sequence.Steps

	context := &ExecutionContext{
		httpClient: http.NewClient(),
		localStore: store.NewStore(nil),
		baseUrl:    definition.BaseUrl,
	}
	for _, step := range steps {
		err := step.Execute(context)

		if err != nil {
			return err
		}
	}

	return nil
}
