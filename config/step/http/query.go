package http

import (
	"strings"

	"gopkg.in/yaml.v3"
)

type QueryParam map[string][]string

func (q *QueryParam) UnmarshalYAML(value *yaml.Node) error {
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

	*q = parsedMap

	return nil
}
