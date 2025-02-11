package workflow

import (
	"fmt"
	"maps"
	"strings"

	"github.com/expr-lang/expr"
	"gopkg.in/yaml.v3"
)

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
// h headers are preferred over global headers.
func (h *Header) compile(context *ExecutionContext) (map[string][]string, error) {
	gHeaders := context.globalOptions.headers
	vars := context.store.Map()

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
	vars := context.store.Map()

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
