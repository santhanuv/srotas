package client

import (
	"net/http"
	"net/url"
)

type Request struct {
	Url        string
	Method     string
	Header     http.Header
	QueryParam url.Values
}

func NewRequest(method string, requestUrl string, body any) *Request {
	return &Request{
		Url:        requestUrl,
		Method:     method,
		Header:     http.Header{},
		QueryParam: url.Values{},
	}
}

func (r *Request) ToHttpRequest() (*http.Request, error) {
	requestUrl, err := url.Parse(r.Url)

	if err != nil {
		return nil, err
	}

	return &http.Request{
		URL:    requestUrl,
		Method: r.Method,
		Header: r.Header,
	}, nil
}
