package step

import "github.com/santhanuv/srotas/config/step/validation"

type RequestBody struct {
	file string
	data map[string]string
}

type RequestStep struct {
	Type        string
	Name        string
	Description string
	Url         string
	Method      string
	Body        *RequestBody
	Headers     map[string]string
	QueryParams map[string]string
	Variables   map[string]any
	Validation  validation.Validation
	Timeout     uint
}

func (r *RequestStep) Execute() error {
	return nil
}

func parseRequestStep(step map[string]any) (*RequestStep, error) {
	errs := make([]error, 0)

	if step["url"] == nil {
		err := NewMissingRequiredField("url")
		errs = append(errs, err)
	}

	if step["method"] == nil {
		err := NewMissingRequiredField("method")
		errs = append(errs, err)
	}

	name, ok := step["name"].(string)
	if step["name"] != nil && !ok {
		err := NewInvalidValueType("name", "string")
		errs = append(errs, err)
	}

	description, ok := step["description"].(string)
	if step["description"] != nil && !ok {
		err := NewInvalidValueType("description", "string")
		errs = append(errs, err)
	}

	url, ok := step["url"].(string)
	if step["url"] != nil && !ok {
		err := NewInvalidValueType("url", "string")
		errs = append(errs, err)
	}

	method, ok := step["method"].(string)
	if step["method"] != nil && !ok {
		err := NewInvalidValueType("method", "string")
		errs = append(errs, err)
	}

	headers, ok := step["headers"].(map[string]string)
	if step["headers"] != nil && !ok {
		err := NewInvalidValueType("headers", "mappings")
		errs = append(errs, err)
	}

	queryParams, ok := step["queryParams"].(map[string]string)
	if step["queryParams"] != nil && !ok {
		err := NewInvalidValueType("queryParams", "mappings")
		errs = append(errs, err)
	}

	variables, ok := step["variables"].(map[string]any)
	if step["variables"] != nil && !ok {
		err := NewInvalidValueType("variables", "mappings")
		errs = append(errs, err)
	}

	timeout, ok := step["timeout"].(uint)
	if step["timeout"] != nil && !ok {
		err := NewInvalidValueType("timeout", "number in ms")
		errs = append(errs, err)
	}

	body, err := parseRequestBody(step["body"])
	if err != nil {
		return nil, err
	}

	return &RequestStep{
		Type:        "request",
		Name:        name,
		Description: description,
		Url:         url,
		Method:      method,
		Body:        body,
		Headers:     headers,
		QueryParams: queryParams,
		Variables:   variables,
		Validation:  nil,
		Timeout:     timeout,
	}, nil
}

func parseRequestBody(body any) (*RequestBody, error) {
	if body == nil {
		return nil, nil
	}

	pb, ok := body.(map[string]any)
	if !ok {
		return nil, NewInvalidValueType("body", "mapping")
	}

	file, ok := pb["file"].(string)
	if !ok && pb["file"] != nil {
		return nil, NewInvalidValueType("body:file", "string")
	}

	data, ok := pb["data"].(map[string]any)
	if !ok && pb["data"] != nil {
		return nil, NewInvalidValueType("body:data", "mapping")
	}

	pd := make(map[string]string, len(data))
	for k, rv := range data {
		val, ok := rv.(string)
		if !ok {
			return nil, NewInvalidValueType("body:data:"+k, "string")
		}
		pd[k] = val
	}

	return &RequestBody{
		file: file,
		data: pd,
	}, nil
}
