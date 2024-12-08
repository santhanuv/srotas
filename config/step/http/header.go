package http

import (
	"log"
	"strings"

	"github.com/santhanuv/srotas/contract"
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

// build replaces the variable reference with the actual value from the context and also appends the global headers if any.
func (h *Header) build(context contract.ExecutionContext) map[string][]string {
	gHeaders := context.GlobalOptions().Headers
	store := context.Store()
	expandedHeader := make(map[string][]string, len(*h)+len(gHeaders))

	for key, values := range *h {
		expandedValues := make([]string, 0, len(values))

		for _, val := range values {
			if strings.HasPrefix(val, "$") {
				rawVarVal, ok := store.Get(val[1:])

				if !ok {
					log.Printf("variable:%s not found in store\n", val[1:])
					continue
				}

				var varVal string
				if varVal, ok = rawVarVal.(string); !ok {
					log.Printf("variable:%s does not have a string value\n", val[1:])
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
