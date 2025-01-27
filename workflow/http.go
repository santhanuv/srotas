package workflow

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/santhanuv/srotas/internal/http"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"gopkg.in/yaml.v3"
)

// Request represents the HTTP request in the workflow step.
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
	Validations Validator
}

// Execute executes the request with the given context.
func (r *Request) Execute(context *executionContext) error {
	req, err := r.build(context)
	if err != nil {
		return err
	}

	context.logger.Info("Sending http request '%s': %s %s", r.Name, req.Method, req.Url)

	delayDuration := time.Duration(r.Delay) * time.Millisecond
	if delayDuration > 0 {
		context.logger.Info("Delaying request for %s", delayDuration)
		time.Sleep(delayDuration)
	}

	res, err := context.httpClient.Do(req)
	if err != nil {
		return err
	}

	context.logger.Info("Http request '%s' responded with status %d", r.Name, res.StatusCode)

	storeFromResponse(res.Body, r.Store, context)

	err = r.Validations.Validate(context, res)

	if err != nil {
		return err
	}

	return nil
}

// build returns a custom http request after expanding all variables.
func (r *Request) build(context *executionContext) (*http.Request, error) {
	body, err := r.Body.build(context)

	if err != nil {
		return nil, err
	}

	gopts := context.globalOptions

	eURL, err := r.buildURL(gopts.baseUrl, context)

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
func (r *Request) buildURL(baseUrl string, context *executionContext) (string, error) {
	if baseUrl == "" {
		return r.Url, nil
	}

	var abURL = ""
	if !strings.Contains(r.Url, "://") {
		baseUrl = strings.TrimSuffix(baseUrl, "/")
		url := strings.TrimPrefix(r.Url, "/")

		abURL = fmt.Sprintf("%s/%s", baseUrl, url)
	}

	store := context.store
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

func storeFromResponse(body []byte, query map[string]string, context *executionContext) error {
	if query == nil {
		return nil
	}

	selectors := make([]string, 0, len(query))
	variables := make([]string, 0, len(query))

	for variable, selector := range query {
		selectors = append(selectors, selector)
		variables = append(variables, variable)
	}

	if ok := gjson.Valid(string(body)); !ok {
		return fmt.Errorf("Error: Invalid json response")
	}

	queryVal := gjson.GetManyBytes(body, selectors...)

	store := context.store

	for idx, qv := range queryVal {
		val := qv.Value()

		if val == nil {
			context.logger.Info("Warning: Setting nil value for %s", variables[idx])
		}

		store.Set(variables[idx], val)
	}

	return nil
}

// CSVMap is a data structure that maps a single key to multiple values.
type CSVMap map[string][]string

func (c *CSVMap) UnmarshalYAML(value *yaml.Node) error {
	const seperator string = ","

	var rawValue map[string]string

	if err := value.Decode(&rawValue); err != nil {
		return err
	}

	parsedMap := make(map[string][]string)
	for key, values := range rawValue {
		valueList := strings.Split(values, seperator)
		parsedMap[key] = valueList
	}

	*c = parsedMap

	return nil
}

// Header represents the HTTP headers for the workflow step.
type Header CSVMap

func (h *Header) UnmarshalYAML(value *yaml.Node) error {
	header := CSVMap{}

	if err := value.Decode(&header); err != nil {
		return err
	}

	*h = Header(header)

	return nil
}

// build replaces the variable reference with the actual value from the context and also appends the global headers if any.
func (h *Header) build(context *executionContext) map[string][]string {
	gHeaders := context.globalOptions.header
	store := context.store
	expandedHeader := make(map[string][]string, len(*h)+len(gHeaders))

	for key, values := range *h {
		expandedValues := make([]string, 0, len(values))

		for _, val := range values {
			if strings.HasPrefix(val, "$") {
				rawVarVal, ok := store.Get(val[1:])

				if !ok {
					context.logger.Info("variable:%s not found in store", val[1:])
					continue
				}

				var varVal string
				if varVal, ok = rawVarVal.(string); !ok {
					context.logger.Info("variable:%s does not have a string value", val[1:])
					continue
				}

				expandedValues = append(expandedValues, varVal)
				continue
			}

			expandedValues = append(expandedValues, val)
		}

		expandedHeader[key] = expandedValues
	}

	for key, values := range gHeaders {
		for _, val := range values {
			expandedHeader[key] = append(expandedHeader[key], val)
		}
	}

	return expandedHeader
}

// QueryParam represents the HTTP query params for the workflow step.
type QueryParam CSVMap

func (q *QueryParam) UnmarshalYAML(value *yaml.Node) error {
	queryParam := CSVMap{}

	if err := value.Decode(&queryParam); err != nil {
		return err
	}

	*q = QueryParam(queryParam)

	return nil
}

// expandVariables replaces the variable references with the acutal value from the context.
func (h *QueryParam) expandVariables(context *executionContext) map[string][]string {
	expandedQP := make(map[string][]string, len(*h))
	store := context.store

	for key, values := range *h {
		expandedValues := make([]string, 0, len(values))

		for _, val := range values {
			if strings.HasPrefix(val, "$") {
				rawVarVal, ok := store.Get(val[1:])

				if !ok {
					context.logger.Info("variable:%s not found in store", val[1:])
					continue
				}

				var varVal string
				if varVal, ok = rawVarVal.(string); !ok {
					context.logger.Info("variable:%s does not have a string value", val[1:])
					continue
				}

				expandedValues = append(expandedValues, varVal)
				continue
			}

			expandedValues = append(expandedValues, val)
		}

		expandedQP[key] = expandedValues
	}

	return expandedQP
}

// RequestBody represents the body of the HTTP request in the workflow step.
type RequestBody struct {
	Content []byte
	Data    map[string]string
}

func (rb *RequestBody) UnmarshalYAML(value *yaml.Node) error {
	var rawRequestBody struct {
		File string
		Data map[string]string
	}

	if err := value.Decode(&rawRequestBody); err != nil {
		return err
	}

	*rb = RequestBody{
		Data: rawRequestBody.Data,
	}

	if rawRequestBody.File != "" {
		file, err := os.Open(rawRequestBody.File)
		defer file.Close()

		if err != nil {
			return err
		}

		content, err := io.ReadAll(file)

		if err != nil {
			return err
		}

		rb.Content = content
	}

	return nil
}

// build merges rb.Data with rb.Content and returns the result.
func (rb *RequestBody) build(context *executionContext) ([]byte, error) {
	store := context.store

	var (
		updatedContent []byte
		err            error
	)

	for field, variable := range rb.Data {
		value, ok := store.Get(variable)

		if !ok {
			context.logger.Info("varable:'%s' not found in store", variable)
			continue
		}

		updatedContent, err = sjson.SetBytes(rb.Content, field, value)

		if err != nil {
			return nil, err
		}
	}

	return updatedContent, nil
}

// Validator is a data structure that represents the validations to the HTTP request.
type Validator struct {
	Status_code *uint
	Asserts     []Assert
}

func (v *Validator) Validate(context *executionContext, response *http.Response) error {
	if v.Status_code != nil && *v.Status_code != response.StatusCode {
		return fmt.Errorf("Status code: Expected %d but got %d", *v.Status_code, response.StatusCode)
	}

	for _, assert := range v.Asserts {
		err := assert.Validate(context, response)

		if err != nil {
			return err
		}
	}

	return nil
}

// Assert represents an assertion where the Value represents the expected value and Selector represents the GJSON string that is used to extract the value from the JSON response.
type Assert struct {
	Value    string
	Selector string
}

func (a *Assert) Validate(context *executionContext, response *http.Response) error {
	var expected any = a.Value

	if strings.HasPrefix(a.Value, "$") {
		var ok bool
		expected, ok = context.store.Get(a.Value[1:])

		if !ok {
			return fmt.Errorf("Invalid variable name in assert")
		}
	}

	actual := gjson.GetBytes(response.Body, a.Selector).Value()

	if actual == nil {
		return fmt.Errorf("Assert failed: value not found in response body")
	}

	if expected != actual {
		return fmt.Errorf("Assert failed: expected '%s' but got '%s'", expected, actual)
	}

	return nil
}
