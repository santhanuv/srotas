package client

import (
	"net/http"
)

type Response struct {
	Status     string
	StatusCode int
	Headers    http.Header
}

func FromHttpResponse(res *http.Response) *Response {
	return &Response{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Headers:    res.Header,
	}
}
