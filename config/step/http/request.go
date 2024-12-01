package http

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/santhanuv/srotas/contract"
)

type Request struct {
	Type        string
	Name        string
	Description string
	Url         string
	Method      string
	Body        *RequestBody `yaml:"body"`
	Headers     Header
	QueryParams QueryParam `yaml:"query_params"`
	Store       map[string]string
	Timeout     uint
	Delay       uint
}

func (r *Request) Execute(context contract.ExecutionContext) error {
	var body io.ReadCloser

	if r.Body != nil {
		var err error
		if body, err = r.Body.transform(context); err != nil {
			return err
		}
	}

	eURL := r.expandURL(context.BaseUrl())
	req, err := http.NewRequest(r.Method, eURL, body)

	if err != nil {
		return err
	}

	for hk, hv := range r.Headers {
		for _, v := range hv {
			req.Header.Add(hk, v)
		}
	}

	urlValues := url.Values{}
	for qp, qv := range r.QueryParams {
		for _, v := range qv {
			urlValues.Add(qp, v)
		}
	}
	req.URL.RawQuery = urlValues.Encode()

	delayDuration := time.Duration(r.Delay) * time.Millisecond
	if delayDuration > 0 {
		log.Printf("Delaying request for %s", delayDuration)
		time.Sleep(delayDuration)
	}

	httpRes, err := context.HttpClient().Do(req)

	if err != nil {
		return err
	}

	resBody, err := io.ReadAll(httpRes.Body)
	defer httpRes.Body.Close()
	if err != nil {
		return err
	}

	log.Printf("Response: %v\n", string(resBody))

	// TODO: update for parsing response based on the type, assuming json here
	storeFromJsonResponse(resBody, r.Store, context)

	return nil
}

func (r *Request) expandURL(baseUrl string) string {
	if baseUrl == "" {
		return r.Url
	}

	if strings.Contains(r.Url, "://") {
		return r.Url
	}

	baseUrl = strings.TrimSuffix(baseUrl, "/")
	url := strings.TrimPrefix(r.Url, "/")

	return fmt.Sprintf("%s/%s", baseUrl, url)
}
