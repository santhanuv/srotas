package workflow_test

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"text/template"

	"github.com/santhanuv/srotas/internal/http"
	"github.com/santhanuv/srotas/internal/log"
	"github.com/santhanuv/srotas/internal/store"
	"github.com/santhanuv/srotas/workflow"
)

func TestHttpRequest_Validate(t *testing.T) {
	tests := []struct {
		name  string
		Http  workflow.Request
		err   bool
		field string
	}{
		{
			name: "'name' is not provided",
			Http: workflow.Request{
				Type:   "http",
				Url:    "url",
				Method: "GET",
			},
			err:   true,
			field: "name",
		},
		{
			name: "'url' field is not provided",
			Http: workflow.Request{
				Type:     "http",
				StepName: "name",
				Method:   "GET",
			},
			err:   true,
			field: "url",
		},
		{
			name: "'method' field is not provided",
			Http: workflow.Request{
				Type:     "http",
				StepName: "name",
				Url:      "url",
			},
			err:   true,
			field: "method",
		},
		{
			name: "All required fields are provided",
			Http: workflow.Request{
				Type:     "http",
				StepName: "name",
				Url:      "url",
				Method:   "GET",
			},
			err: false,
		},
	}

	for _, tt := range tests {
		err := tt.Http.Validate()

		if !tt.err && err != nil {
			t.Errorf("in test %q; expected no error but got %q", tt.name, err)
			continue
		}

		if tt.err {
			if err == nil {
				t.Errorf("in test %q; expected error but got none", tt.name)
				continue
			}

			var vErr *workflow.ValidationError
			if !errors.As(err, &vErr) {
				t.Errorf("in test %q; got %q; want validation error", tt.name, err)
			}
		}
	}
}

type mockHttpClient struct {
	expectedRes *http.Response
	err         error
	validator   func(*http.Request) error
}

var ErrFailedHttpRequestValidation = errors.New("validation of http request failed while testing")

func (m *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	if err := m.validator(req); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedHttpRequestValidation, err)
	}

	if m.err != nil {
		return nil, m.err
	}

	return m.expectedRes, nil
}

func TestHttpRequest_Execute_PostWithRequestBody(t *testing.T) {
	baseUrl := "https://domain.com"
	method := "POST"
	endPoint := "test"
	requestBodyTemplate := `{"Name": "{{ .name }}"}`

	requestBodyTemplateData := map[string]string{
		"name": "'test name'",
	}
	headers := &workflow.Header{
		"Content-Type": []string{"'application/json'"},
	}

	reqTempl, err := template.New("http request").Parse(requestBodyTemplate)
	if err != nil {
		t.Fatalf("failed to setup test: failed to create request template")
	}

	reqBody := &workflow.RequestBody{
		Template: reqTempl,
		Data:     requestBodyTemplateData,
	}

	req := workflow.Request{
		Type:     "http",
		StepName: "Http request",
		Url:      endPoint,
		Method:   method,
		Body:     reqBody,
		Headers:  headers,
		Validations: &workflow.Validator{
			Status_code: 201,
		},
	}

	mockHttpClient := mockHttpClient{
		expectedRes: &http.Response{
			StatusCode: 201,
			Status:     "201 Created",
		},
		err: nil,
		validator: func(req *http.Request) error {
			if req.Method != method {
				return fmt.Errorf("expected http method to be '%s' but got '%s'", method, req.Method)
			}

			expectedURL := fmt.Sprintf("%s/%s", baseUrl, endPoint)
			if req.Url != expectedURL {
				return fmt.Errorf("expected url to be '%s' but got '%s'", expectedURL, req.Url)
			}

			exepctedBody := []byte(`{"Name": "test name"}`)
			if !bytes.Equal(req.Body, exepctedBody) {
				return fmt.Errorf("expected request body to be '%s' but got '%s'", string(exepctedBody), string(req.Body))
			}

			expectedHeaders := map[string][]string{"Content-Type": {"application/json"}}
			if !reflect.DeepEqual(req.Headers, expectedHeaders) {
				return fmt.Errorf("expected http request headers to be %v but got %v", expectedHeaders, req.Headers)
			}

			return nil
		},
	}

	mockStore := store.NewStore(nil)

	logBuf := bytes.NewBuffer(nil)
	logger := log.New(logBuf, logBuf, logBuf)

	execContext, err := workflow.NewExecutionContext(workflow.WithGlobalOptions(baseUrl, nil),
		workflow.WithHttpClient(&mockHttpClient),
		workflow.WithStore(mockStore),
		workflow.WithLogger(logger))
	if err != nil {
		t.Fatalf("failed to setup test: unable to create execution context")
	}

	err = req.Execute(execContext)
	if err != nil {
		t.Fatalf("expected no error but got %q\nLogs: %s", err, logBuf.String())
	}
}

func TestHttpRequest_Execute_GetWithResponseBody(t *testing.T) {
	baseUrl := "https://domain.com"
	method := "GET"
	endPoint := "test"

	req := workflow.Request{
		Type:     "http",
		StepName: "Http request",
		Url:      endPoint,
		Method:   method,
		Store:    map[string]string{"name": "response.Name"},
		Validations: &workflow.Validator{
			Status_code: 200,
			Asserts:     []workflow.Assert{`response.Id == "123"`},
		},
	}

	mockHttpClient := mockHttpClient{
		expectedRes: &http.Response{
			StatusCode: 200,
			Status:     "200 OK",
			Body:       []byte(`{"Id":"123","Name":"test name"}`),
			Headers:    map[string][]string{"Content-Type": {"application/json"}},
		},
		err: nil,
		validator: func(req *http.Request) error {
			if req.Method != method {
				return fmt.Errorf("expected http method to be '%s' but got '%s'", method, req.Method)
			}

			expectedURL := fmt.Sprintf("%s/%s", baseUrl, endPoint)
			if req.Url != expectedURL {
				return fmt.Errorf("expected url to be '%s' but got '%s'", expectedURL, req.Url)
			}

			return nil
		},
	}

	mockStore := store.NewStore(nil)

	logBuf := bytes.NewBuffer(nil)
	logger := log.New(logBuf, logBuf, logBuf)

	execContext, err := workflow.NewExecutionContext(workflow.WithGlobalOptions(baseUrl, nil),
		workflow.WithHttpClient(&mockHttpClient),
		workflow.WithStore(mockStore),
		workflow.WithLogger(logger))
	if err != nil {
		t.Fatalf("failed to setup test: unable to create execution context")
	}

	err = req.Execute(execContext)
	if err != nil {
		t.Fatalf("expected no error but got %q\nLogs: %s", err, logBuf.String())
	}

	name, ok := mockStore.Get("name")

	if !ok || name == nil {
		t.Fatalf("expected 'name' to be stored from response but got nil\nLogs: %s", logBuf.String())
	}

	if str, ok := name.(string); ok && str != "test name" {
		t.Fatalf("expected string %q; but got '%v'\nLogs: %s", "test name", name, logBuf.String())
	}
}
