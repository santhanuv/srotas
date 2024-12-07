package contract

import "github.com/santhanuv/srotas/internal/http"

type ExecutionContext interface {
	HttpClient() *http.Client
	Store() Store
	BaseUrl() string
}
