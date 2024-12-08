package http

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/santhanuv/srotas/contract"
	"github.com/santhanuv/srotas/internal/http"
)

type Request struct {
	Type        string
	Name        string
	Description string
	Url         string
	Method      string
	Body        RequestBody `yaml:"body"`
	Headers     Header
	QueryParams QueryParam `yaml:"query_params"`
	Store       map[string]string
	Timeout     uint
	Delay       uint
}

// Execute executes the request with the given context.
func (r *Request) Execute(context contract.ExecutionContext) error {
	req, err := r.build(context)
	if err != nil {
		return err
	}

	log.Printf("Executing '%s' with request: %v\n", r.Name, req)

	delayDuration := time.Duration(r.Delay) * time.Millisecond
	if delayDuration > 0 {
		log.Printf("Delaying request for %s\n", delayDuration)
		time.Sleep(delayDuration)
	}

	res, err := context.HttpClient().Do(req)
	if err != nil {
		return err
	}

	log.Printf("Response for '%s': %v\n", r.Name, string(res.Body))

	storeFromResponse(res.Body, r.Store, context)

	return nil
}

// build returns a custom http request after expanding all variables.
func (r *Request) build(context contract.ExecutionContext) (*http.Request, error) {
	body, err := r.Body.build(context)

	if err != nil {
		return nil, err
	}

	gopts := context.GlobalOptions()

	eURL := r.expandURL(gopts.BaseUrl)
	headers := r.Headers.build(context)
	queryParams := r.QueryParams.expandVariables(context)

	req := http.Request{
		Method:      r.Method,
		Url:         eURL,
		Body:        body,
		Headers:     headers,
		QueryParams: queryParams,
	}

	return &req, nil
}

// expandURL combines baseUrl with r.Url. If r.Url is the full url, it is returned.
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
