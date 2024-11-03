package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	Url        string
	Method     string
	Header     http.Header
	QueryParam url.Values
	Body       []byte
}

func NewRequest(method string, requestUrl string, body []byte) *Request {
	return &Request{
		Url:        requestUrl,
		Method:     method,
		Header:     http.Header{},
		QueryParam: url.Values{},
		Body:       body,
	}
}

func (r *Request) ToHttpRequest() (*http.Request, error) {
	contentType := r.Header.Get("Content-Type")

	if r.Body != nil && strings.HasPrefix(contentType, "application/json") {
		isValid := json.Valid(r.Body)

		if !isValid {
			return nil, fmt.Errorf("Invalid json: %s", string(r.Body))
		}

	} else if r.Body != nil {
		fmt.Println("Sending the request body as text/plain content type")
		r.Header.Set("Content-Type", "text/plain")
	}

	request, err := http.NewRequest(r.Method, r.Url, bytes.NewBuffer(r.Body))

	if err != nil {
		return nil, err
	}

	request.URL.RawQuery = r.QueryParam.Encode()
	request.Header = r.Header

	return request, nil
}

func (r *Request) SetQueryParams(queryParams map[string][]string) {
	for key, values := range queryParams {
		for _, value := range values {
			r.QueryParam.Add(key, value)
		}
	}
}

func (r *Request) SetHeaders(headers map[string][]string) {
	for key, values := range headers {
		for _, value := range values {
			r.Header.Add(key, value)
		}
	}
}
