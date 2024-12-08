package http

import (
	"net/http"
	"time"
)

// Client represents an http client
// It uses the native client from net/http package
type Client struct {
	httpClient http.Client // the underlying http client
}

// NewClient returns a new http client with the specified timeout in milliseconds
func NewClient(timeout uint) *Client {
	timeoutMS := time.Duration.Milliseconds(time.Duration(timeout))
	c := http.Client{
		Timeout: time.Duration(timeoutMS),
	}
	return &Client{
		httpClient: c,
	}
}

// Do sends an http request and returns an http resposne
func (hc *Client) Do(request *Request) (*Response, error) {
	req, err := request.buildNative()

	if err != nil {
		return nil, err
	}

	res, err := hc.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	response, err := buildFromNative(res)

	if err != nil {
		return nil, err
	}

	return response, nil
}
