package client

import (
	"net/http"
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

func (c *Client) Do(req Request) (*Response, error) {
	httpRequest, err := req.ToHttpRequest()

	if err != nil {
		return nil, err
	}

	httpResponse, err := c.httpClient.Do(httpRequest)

	if err != nil {
		return nil, err
	}

	return (*Response)(httpResponse), nil
}
