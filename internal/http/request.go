package http

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

// Request represents a http request.
type Request struct {
	// URL specifies URL to access.
	Url string
	// Method specifies the http method (GET, POST, PUT, etd.)
	// An empty string means GET
	Method string
	// Body is the  request body to be send with the request.
	// A nil body means the request has no body.
	Body []byte
	// Headers contains the headers to be sent with the request
	Headers map[string][]string
	// QueryParams specifies the query parameters to be added in the request URL
	QueryParams map[string][]string
}

// buildHttpRequest builds the native http.Request from the custom Request type.
func (hr *Request) buildHttpRequest() (*http.Request, error) {
	var body io.Reader

	if hr.Body != nil {
		body = bytes.NewBuffer(hr.Body)
	}

	req, err := http.NewRequest(hr.Method, hr.Url, body)

	if err != nil {
		return nil, err
	}

	req.Header = http.Header(hr.Headers)
	req.URL.RawQuery = url.Values(hr.QueryParams).Encode()

	return req, nil
}
