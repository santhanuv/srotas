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

	requestUrl.RawQuery = r.QueryParam.Encode()

	if err != nil {
		return nil, err
	}

	return &http.Request{
		URL:    requestUrl,
		Method: r.Method,
		Header: r.Header,
	}, nil
}

func (r *Request) SetQueryParams(queryParams map[string][]string) {
	for key, values := range queryParams {
		for _, value := range values {
			r.QueryParam.Add(key, value)
		}
	}
}
