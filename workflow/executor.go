package workflow

import (
	"os"

	"github.com/santhanuv/srotas/internal/http"
	"github.com/santhanuv/srotas/internal/log"
	"github.com/santhanuv/srotas/internal/store"
)

// ExecutionContext holds contextual data for the execution of a config.
// It manages execution-specific state, including variables, headers, and other necessary metadata.
type ExecutionContext struct {
	httpClient    *http.Client   // Client used for the execution of http request.
	store         *store.Store   // Store used in the config execution.
	globalOptions *globalOptions // Global options for config execution.
	logger        *log.Logger    // Logger used in the config execution.
}

// globalOptions stores the base URL and global headers used during config execution.
type globalOptions struct {
	baseUrl string              // baseUrl to be used in the config execution.
	headers map[string][]string // global headers to be used in the config execution.
}

// ExecutionOption is the option for configuring [ExecutionContext].
type ExecutionOption func(context *ExecutionContext) error

// NewExecutionContext initializes and returns a new [ExecutionContext] with the specified options.
func NewExecutionContext(options ...ExecutionOption) (*ExecutionContext, error) {
	var context ExecutionContext

	for _, option := range options {
		err := option(&context)

		if err != nil {
			return nil, err
		}
	}

	if context.store == nil {
		context.store = store.NewStore(nil)
	}

	if context.logger == nil {
		context.logger = log.New(os.Stderr, os.Stderr, os.Stderr)
		context.logger.SetDebugMode(true)
	}

	if context.httpClient == nil {
		context.httpClient = http.NewClient(15000)
	}

	return &context, nil
}

// Execute executes the given definition with the specified context.
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

// WithGlobalOptions configures the [ExecutionContext] with the baseUrl and headers.
func WithGlobalOptions(baseUrl string, headers map[string][]string) ExecutionOption {
	return func(context *ExecutionContext) error {
		gOpts := globalOptions{
			baseUrl: baseUrl,
			headers: headers,
		}

		context.globalOptions = &gOpts

		return nil
	}
}

// WithHttpClient configures the [ExecutionContext] with the specified client.
func WithHttpClient(client *http.Client) ExecutionOption {
	return func(context *ExecutionContext) error {
		context.httpClient = client

		return nil
	}
}

// WithLogger configures the [ExecutionContext] with the specified logger.
func WithLogger(logger *log.Logger) ExecutionOption {
	return func(context *ExecutionContext) error {
		context.logger = logger

		return nil
	}
}

// WithStore configures the [ExecutionContext] with the specified store.
func WithStore(store *store.Store) ExecutionOption {
	return func(context *ExecutionContext) error {
		context.store = store

		return nil
	}
}

// Variables returns all the variables in the store as a map.
func (e *ExecutionContext) Variables() map[string]any {
	return e.store.Map()
}
