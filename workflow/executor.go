package workflow

import (
	"os"

	"github.com/santhanuv/srotas/internal/http"
	"github.com/santhanuv/srotas/internal/log"
	"github.com/santhanuv/srotas/internal/store"
)

type ExecutionContext struct {
	httpClient    *http.Client
	store         *store.Store
	globalOptions *globalOptions
	logger        *log.Logger
}

type globalOptions struct {
	baseUrl string
	header  map[string][]string
}

type ExecutionOption func(context *ExecutionContext) error

func NewExecutionContext(options ...ExecutionOption) (*ExecutionContext, error) {
	var context ExecutionContext

	for _, option := range options {
		err := option(&context)

		if err != nil {
			return nil, err
		}
	}

	if context.httpClient == nil {
		context.httpClient = http.NewClient(15000)
	}

	if context.store == nil {
		context.store = store.NewStore(nil)
	}

	if context.logger == nil {
		context.logger = log.New(os.Stderr, os.Stderr, os.Stderr)
	}

	return &context, nil
}

func Execute(definition *Definition, context *ExecutionContext) error {
	steps := definition.Steps

	for _, step := range steps {
		err := step.Execute(context)

		if err != nil {
			return err
		}
	}

	return nil
}

func WithGlobalOptions(baseUrl string, header map[string][]string) ExecutionOption {
	return func(context *ExecutionContext) error {
		gOpts := globalOptions{
			baseUrl: baseUrl,
			header:  header,
		}

		context.globalOptions = &gOpts

		return nil
	}
}

func WithHttpClient(client *http.Client) ExecutionOption {
	return func(context *ExecutionContext) error {
		context.httpClient = client

		return nil
	}
}

func WithLogger(logger *log.Logger) ExecutionOption {
	return func(context *ExecutionContext) error {
		context.logger = logger

		return nil
	}
}

func WithStore(store *store.Store) ExecutionOption {
	return func(context *ExecutionContext) error {
		context.store = store

		return nil
	}
}

func (e *ExecutionContext) Variables() map[string]any {
	return e.store.ToMap()
}
