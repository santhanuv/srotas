package http

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/santhanuv/srotas/config/step/http/validation"
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
	Validations validation.Validator
}

// Execute executes the request with the given context.
func (r *Request) Execute(context contract.ExecutionContext) error {
	req, err := r.build(context)
	if err != nil {
		return err
	}

	log.Printf("Sending http request '%s': %s %s\n", r.Name, req.Method, req.Url)

	delayDuration := time.Duration(r.Delay) * time.Millisecond
	if delayDuration > 0 {
		log.Printf("Delaying request for %s\n", delayDuration)
		time.Sleep(delayDuration)
	}

	res, err := context.HttpClient().Do(req)
	if err != nil {
		return err
	}

	log.Printf("Http request '%s' responded with status %d", r.Name, res.StatusCode)

	storeFromResponse(res.Body, r.Store, context)

	err = r.Validations.Validate(context, res)

	if err != nil {
		return err
	}

	return nil
}

// build returns a custom http request after expanding all variables.
func (r *Request) build(context contract.ExecutionContext) (*http.Request, error) {
	body, err := r.Body.build(context)

	if err != nil {
		return nil, err
	}

	gopts := context.GlobalOptions()

	eURL, err := r.buildURL(gopts.BaseUrl, context)

	if err != nil {
		return nil, err
	}

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

// buildURL combines baseUrl with r.Url and expands any URL parameters. If r.Url is already a fully qualified URL, it is returned as-is, just expanding url parameters.
func (r *Request) buildURL(baseUrl string, context contract.ExecutionContext) (string, error) {
	if baseUrl == "" {
		return r.Url, nil
	}

	var abURL = ""
	if !strings.Contains(r.Url, "://") {
		baseUrl = strings.TrimSuffix(baseUrl, "/")
		url := strings.TrimPrefix(r.Url, "/")

		abURL = fmt.Sprintf("%s/%s", baseUrl, url)
	}

	store := context.Store()
	for idx, urlParam := range strings.Split(abURL, "/:") {
		if idx == 0 {
			continue
		}

		val, ok := store.Get(urlParam)
		if !ok {
			return "", fmt.Errorf("Url prameter not found in variable store: '%s'", urlParam)
		}

		abURL = strings.ReplaceAll(abURL, fmt.Sprintf(":%s", urlParam), fmt.Sprintf("%v", val))
	}

	return abURL, nil
}
