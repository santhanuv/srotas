package workflow

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"os"
	"strings"
	"time"

	"github.com/expr-lang/expr"
	"github.com/santhanuv/srotas/internal/http"
	"github.com/tidwall/sjson"
	"gopkg.in/yaml.v3"
)

// Request represents the HTTP request in the workflow step.
type Request struct {
	Type        string
	Name        string
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
func (r *Request) Execute(context *ExecutionContext) error {
	context.logger.Debug("Executing http step '%s'", r.Name)

	req, err := r.build(context)
	if err != nil {
		return err
	}

	context.logger.Info("Sending http request '%s': %s %s", r.Name, req.Method, req.Url)
	context.logger.DebugJson(req.Body, "Http request body:")

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
	context.logger.DebugJson(res.Body, "Http response body:")

	var pb any
	err = json.Unmarshal(res.Body, &pb)

	resBody := responseBody{
		body: pb,
	}

	if err != nil {
		return fmt.Errorf("http response parse error: %v", err)
	}

	err = resBody.store(r.Store, context)

	if err != nil {
		return err
	}

	context.logger.Debug("Running validations on the http response")
	err = r.Validations.Validate(context, res.StatusCode, &resBody)

	if err != nil {
		fres := struct {
			StatusCode uint
			Body       any
		}{
			StatusCode: res.StatusCode,
			Body:       resBody.body,
		}
		jres, je := json.MarshalIndent(fres, "", " ")

		if je != nil {
			return fmt.Errorf("unable to output response: %v", je)
		}

		return fmt.Errorf("%v\nResponse: %s", err, string(jres))
	}

	context.logger.Debug("Completed validating the http response")

	return nil
}

// build returns a custom http request after evaluating all value expressions.
func (r *Request) build(context *ExecutionContext) (*http.Request, error) {
	body, err := r.Body.build(context)

	if err != nil {
		return nil, err
	}

	gopts := context.globalOptions

	eURL, err := r.buildURL(gopts.baseUrl, context)

	if err != nil {
		return nil, err
	}

	headers, err := r.Headers.compile(context)

	if err != nil {
		return nil, err
	}

	queryParams, err := r.QueryParams.compile(context)

	if err != nil {
		return nil, err
	}

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
func (r *Request) buildURL(baseUrl string, context *ExecutionContext) (string, error) {
	if r.Url == "" {
		return "", fmt.Errorf("invalid url '%s' for http request '%s'", r.Url, r.Name)
	}

	store := context.store
	for idx, urlParam := range strings.Split(r.Url, "/:") {
		if idx == 0 {
			continue
		}

		val, ok := store.Get(urlParam)
		if !ok {
			return "", fmt.Errorf("variable not found: url '%s': '%s'", urlParam, r.Url)
		}

		r.Url = strings.ReplaceAll(r.Url, fmt.Sprintf(":%s", urlParam), fmt.Sprintf("%v", val))
	}

	var abURL = r.Url

	if !strings.Contains(r.Url, "://") {
		baseUrl = strings.TrimSuffix(baseUrl, "/")
		url := strings.TrimPrefix(r.Url, "/")

		abURL = fmt.Sprintf("%s/%s", baseUrl, url)
	}

	return abURL, nil
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

// compile returns the compiled headers after evaluating the value expressions of headers and also appends the global headers if any.
func (h *Header) compile(context *ExecutionContext) (map[string][]string, error) {
	gHeaders := context.globalOptions.headers
	vars := context.store.ToMap()

	cHeaders := make(map[string][]string, len(*h)+len(gHeaders))
	maps.Copy(cHeaders, gHeaders)

	for key, ves := range *h {
		vals := make([]string, 0, len(ves))

		for _, ve := range ves {
			val, err := expr.Eval(ve, vars)

			if err != nil {
				e := fmt.Errorf("invalid expression '%s' for header '%s': %v", ve, key, err)
				return nil, e
			}

			v, ok := val.(string)

			if !ok {
				e := fmt.Errorf("expression '%s' for header '%s' should evaluate to string", ve, key)
				return nil, e
			}

			vals = append(vals, v)
		}

		cHeaders[key] = vals
	}

	return cHeaders, nil
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

// compile returns the compiled query parameters after evaluating the value expressions for each query parameter.
func (q *QueryParam) compile(context *ExecutionContext) (map[string][]string, error) {
	cqps := make(map[string][]string, len(*q))
	vars := context.store.ToMap()

	for key, ves := range *q {
		vals := make([]string, 0, len(ves))

		for _, ve := range ves {
			val, err := expr.Eval(ve, vars)

			if err != nil {
				e := fmt.Errorf("invalid expression '%s' for query param '%s': %v", ve, key, err)
				return nil, e
			}

			v, ok := val.(string)

			if !ok {
				e := fmt.Errorf("expression '%s' for query param '%s' should evalute to string", ve, key)
				return nil, e
			}

			vals = append(vals, v)
		}

		cqps[key] = vals
	}

	return cqps, nil
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

// build builds the request body with rb.Content as the base and updates the field values after evaluating expressions in rb.Data.
func (rb *RequestBody) build(context *ExecutionContext) ([]byte, error) {
	vars := context.store.ToMap()

	var updatedContent []byte = rb.Content

	for f, e := range rb.Data {
		val, err := expr.Eval(e, vars)

		if err != nil {
			return nil, fmt.Errorf("expression '%s' cannot be evaluated for variable '%s': %v", e, f, err)
		}

		updatedContent, err = sjson.SetBytes(rb.Content, f, val)

		if err != nil {
			return nil, err
		}
	}

	return updatedContent, nil
}

type responseBody struct {
	body any `expr:"response"`
}

// store stores the new set of variables after evaluating the variable expressions in varExprs
func (rb *responseBody) store(varExprs map[string]string, context *ExecutionContext) error {
	context.logger.Debug("Storing variables from response")

	if varExprs == nil {
		return nil
	}

	newVars := make(map[string]any, len(varExprs))

	vars := context.store.ToMap()
	vars["response"] = rb.body

	for vn, ve := range varExprs {
		val, err := expr.Eval(ve, vars)

		if err != nil {
			return fmt.Errorf("invalid expression '%s' for variable '%s': %v", ve, vn, err)
		}

		newVars[vn] = val
	}

	context.store.Add(newVars)

	return nil
}

// Validator is a data structure that represents the validations to the HTTP request.
type Validator struct {
	Status_code *uint
	Asserts     []Assert
}

func (v *Validator) Validate(context *ExecutionContext, statusCode uint, rb *responseBody) error {
	if v.Status_code != nil && *v.Status_code != statusCode {
		return fmt.Errorf("Status code: Expected %d but got %d", *v.Status_code, statusCode)
	}

	vars := context.store.ToMap()
	vars["response"] = rb.body

	for _, assert := range v.Asserts {
		err := assert.Validate(vars, rb)

		if err != nil {
			return err
		}
	}

	return nil
}

// Assert represents an assertion where the Value represents the expected value and Selector represents the GJSON string that is used to extract the value from the JSON response.
type Assert string

func (a *Assert) Validate(vars map[string]any, rb *responseBody) error {
	val, err := expr.Eval(string(*a), vars)

	if err != nil {
		return fmt.Errorf("invalid expression '%s' for assert: %v", *a, err)
	}

	isValid, ok := val.(bool)

	if !ok {
		return fmt.Errorf("evaluating expression '%s' should produce a boolean for assert", *a)
	}

	if !isValid {
		return fmt.Errorf("assertion '%s' failed", *a)
	}

	return nil
}
