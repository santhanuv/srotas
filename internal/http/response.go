package http

import (
	"io"
	"net/http"
)

// Response represents the response from an http request
type Response struct {
	Status     string              // specifies the status of the response
	StatusCode uint                // specifies the http status code of response
	Headers    map[string][]string // headers set in the response
	Body       []byte              // the response body
}

// buildFromNative creates a Response instance from the native http.Response instance
func buildFromNative(response *http.Response) (*Response, error) {
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	return &Response{
		Status:     response.Status,
		StatusCode: uint(response.StatusCode),
		Headers:    response.Header,
		Body:       responseBody,
	}, nil
}
