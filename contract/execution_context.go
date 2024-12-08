package contract

import "github.com/santhanuv/srotas/internal/http"

type ExecutionContext interface {
	HttpClient() *http.Client
	Store() Store
	GlobalOptions() *Options
}

type Options struct {
	BaseUrl string
	Headers map[string][]string
}
