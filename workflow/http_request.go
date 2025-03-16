package workflow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/expr-lang/expr"
	"github.com/santhanuv/srotas/internal/http"
	"gopkg.in/yaml.v3"
)

// Request represents an HTTP request step in the execution flow.
type Request struct {
	Type        string            // The type of the step.
	StepName    string            `yaml:"name"` // Identifier for the step.
	Url         string            // The target URL for the request.
	Method      string            // The HTTP method (e.g., GET, POST).
	Body        *RequestBody      `yaml:"body"` // Request payload.
	Headers     *Header           // Custom headers for the request.
	QueryParams *QueryParam       `yaml:"query_params"` // Query parameters to append to the URL.
	Store       map[string]string // Variables mapped to expressions evaluated using the response.
	Delay       uint              // Wait time (milliseconds) before executing the request.
	Validations *Validator        // Validation rules for the response.
}

// Validate checks the fields of the [Request] step and returns a list of validation errors, if any.
func (r *Request) Validate() error {
	vErr := ValidationError{}
	if r.StepName == "" {
		err := RequiredFieldError{Field: "name"}
		vErr.Add(err)
	}

	if r.Url == "" {
		err := RequiredFieldError{Field: "url"}
		vErr.Add(err)
	}

	if r.Method == "" {
		err := RequiredFieldError{Field: "method"}
		vErr.Add(err)
	}

	if vErr.HasError() {
		return fmt.Errorf("http request step: %w", &vErr)
	}

	return nil
}

func (r *Request) Name() string {
	return r.StepName
}

// Execute executes the step with the specified context.
func (r *Request) Execute(context *ExecutionContext) error {
	req, err := r.build(context)
	if err != nil {
		return fmt.Errorf("failed executing http request '%s': %v", r.StepName, err)
	}

	context.logger.Info("sending http request '%s': %s %s", r.StepName, req.Method, req.Url)
	context.logger.DebugJson(req.Body, "http request: ")

	delayDuration := time.Duration(r.Delay) * time.Millisecond
	if delayDuration > 0 {
		context.logger.Info("Delaying request for %s", delayDuration)
		time.Sleep(delayDuration)
	}

	res, err := context.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed executing http request '%s': %w", r.StepName, err)
	}

	context.logger.Info("http request '%s' responded with status %d", r.StepName, res.StatusCode)
	context.logger.DebugJson(res.Body, "http response: ")

	var pb any

	if res.Body != nil {
		err = json.Unmarshal(res.Body, &pb)
		if err != nil {
			return fmt.Errorf("failed to parse response body for Http request %q: %v", r.StepName, err)
		}
	}

	resBody := responseBody{
		body: pb,
	}

	context.store.Set("response", resBody.body)
	defer context.store.Remove("response")

	err = resBody.store(r.Store, context)
	if err != nil {
		return fmt.Errorf("failed executing http request '%s': %v", r.StepName, err)
	}

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
			return fmt.Errorf("http request '%s': unable to output response: %v", r.StepName, je)
		}

		return fmt.Errorf("%v\nresponse: %s", err, string(jres))
	}

	context.logger.Debug("http response validation has been completed successfully.")

	return nil
}

// build returns a custom http request after evaluating all value expressions.
// Builds full URL using the base_url and the r.Url.
// r.Headers and r.QueryParams are also compiled.
func (r *Request) build(context *ExecutionContext) (*http.Request, error) {
	gopts := context.globalOptions
	eURL, err := r.buildURL(gopts.baseUrl, context)
	if err != nil {
		return nil, err
	}

	req := http.Request{
		Method: r.Method,
		Url:    eURL,
	}

	if r.Body != nil {
		body, err := r.Body.build(context)
		if err != nil {
			return nil, err
		}

		req.Body = body
	}

	headers, err := r.Headers.compile(context)
	if err != nil {
		return nil, err
	}

	req.Headers = headers

	queryParams, err := r.QueryParams.compile(context)
	if err != nil {
		return nil, err
	}

	req.QueryParams = queryParams

	return &req, nil
}

// buildURL combines baseUrl with r.Url and expands any URL parameters. If r.Url is already a fully qualified URL, it is returned as-is, just expanding url parameters.
func (r *Request) buildURL(baseUrl string, context *ExecutionContext) (string, error) {
	if r.Url == "" {
		return "", fmt.Errorf("invalid url '%s'", r.Url)
	}

	store := context.store
	for idx, urlParam := range strings.Split(r.Url, "/:") {
		if idx == 0 {
			continue
		}

		val, ok := store.Get(urlParam)
		if !ok {
			return "", fmt.Errorf("variable '%s' not found for url '%s'", urlParam, r.Url)
		}

		r.Url = strings.ReplaceAll(r.Url, fmt.Sprintf(":%s", urlParam), fmt.Sprintf("%v", val))
	}

	abURL := r.Url

	if !strings.Contains(r.Url, "://") {
		baseUrl = strings.TrimSuffix(baseUrl, "/")
		url := strings.TrimPrefix(r.Url, "/")

		abURL = fmt.Sprintf("%s/%s", baseUrl, url)
	}

	return abURL, nil
}

// RequestBody represents the payload for an HTTP request step.
//   - Data is a map where keys represent JSON fields to update or add,
//     and values are expressions evaluated at runtime before being inserted into the Content.
type RequestBody struct {
	Template *template.Template // Raw JSON payload.
	Data     map[string]string  // Dynamic fields evaluated and added/updated in Content.
}

func (rb *RequestBody) UnmarshalYAML(value *yaml.Node) error {
	var rawRb struct {
		Template string
		File     string
		Data     map[string]string
	}

	if err := value.Decode(&rawRb); err != nil {
		return err
	}

	*rb = RequestBody{
		Template: nil,
		Data:     rawRb.Data,
	}

	if rawRb.Template == "" && rawRb.File == "" {
		return fmt.Errorf("template or file should be provided for request body")
	}

	if rawRb.Template != "" {
		tmpl, err := template.New("request").Parse(rawRb.Template)
		if err != nil {
			return fmt.Errorf("request template error: %v", err)
		}

		rb.Template = tmpl

		return nil
	}

	if rawRb.File != "" {
		tmpl, err := template.New(rawRb.File).ParseFiles(rawRb.File)
		if err != nil {
			return fmt.Errorf("request template error: %v", err)
		}

		rb.Template = tmpl

		return nil
	}

	return fmt.Errorf("no template provided for request body")
}

// build builds the request body with rb.Content as the base and updates the field values after evaluating expressions in rb.Data.
func (rb *RequestBody) build(context *ExecutionContext) ([]byte, error) {
	if rb == nil {
		return nil, nil
	}

	vars := context.store.Map()

	tvars := map[string]any{}

	for v, e := range rb.Data {
		val, err := expr.Eval(e, vars)
		if err != nil {
			return nil, fmt.Errorf("expression '%s' cannot be evaluated for variable '%s': %v", e, v, err)
		}

		tvars[v] = val
	}

	var buf bytes.Buffer
	err := rb.Template.Execute(&buf, tvars)
	if err != nil {
		return nil, fmt.Errorf("error executing template: %v", err)
	}

	return buf.Bytes(), nil
}
