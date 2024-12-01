package contract

import "net/http"

type ExecutionContext interface {
	HttpClient() *http.Client
	Store() Store
}
