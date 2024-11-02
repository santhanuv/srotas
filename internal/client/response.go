package client

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type Response struct {
	Status     string
	StatusCode int
	Headers    http.Header
	Body       any
}

func ParseHttpResponse(res *http.Response) (*Response, error) {
	response := &Response{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Headers:    res.Header,
	}

	defer res.Body.Close()
	contentType := res.Header.Get("Content-Type")

	var err error
	if strings.HasPrefix(contentType, "application/json") {
		var parsedBody any
		err = json.NewDecoder(res.Body).Decode(&parsedBody)

		response.Body = parsedBody
	} else {
		var parsedBody []byte
		parsedBody, err = io.ReadAll(res.Body)

		response.Body = string(parsedBody)
	}

	if err != nil {
		return nil, err
	}

	return response, nil
}
