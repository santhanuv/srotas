package http

import (
	"strings"

	"gopkg.in/yaml.v3"
)

type Header map[string][]string

func (h *Header) UnmarshalYAML(value *yaml.Node) error {
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

	*h = Header(parsedMap)

	return nil
}
