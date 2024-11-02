package client

import (
	"fmt"
	"net/http"
	"net/url"
)

type ApiClient struct {
	httpClient *http.Client
	Request    *http.Request
	Response   *http.Response
}

func NewApiClient(method string, rawRequestUrl string) *ApiClient {
	requestUrl, err := url.Parse(rawRequestUrl)

	if err != nil {
		return nil
	}

	return &ApiClient{
		httpClient: &http.Client{},
		Request: &http.Request{
			URL:    requestUrl,
			Method: method,
		},
		Response: nil,
	}
}

func (ac *ApiClient) Do() {
	fmt.Println("Requesting {} {}", ac.Request.Method, ac.Request.URL)
	res, err := ac.httpClient.Do(ac.Request)

	if err != nil {
		fmt.Println(err)
	}

	ac.Response = res
}
