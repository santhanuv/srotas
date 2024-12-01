package http

import (
	"io"
	"log"
	"net/http"
	"net/url"

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
}

func (r *Request) Execute(context contract.ExecutionContext) error {
	var body io.ReadCloser

	if r.Body != nil {
		var err error
		if body, err = r.Body.transform(context); err != nil {
			return err
		}
	}

	req, err := http.NewRequest(r.Method, r.Url, body)

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
