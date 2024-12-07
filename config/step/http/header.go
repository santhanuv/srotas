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

// expandVariables replaces the variable reference with the actual value from the context store.
func (h *Header) expandVariables(context contract.ExecutionContext) map[string][]string {
	expandedHeader := make(map[string][]string, len(*h))
	store := context.Store()

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

	return expandedHeader
}
