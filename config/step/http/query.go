package http

import (
	"log"
	"strings"

	"github.com/santhanuv/srotas/contract"
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

// expandVariables replaces the variable references with the acutal value from the context.
func (h *QueryParam) expandVariables(context contract.ExecutionContext) map[string][]string {
	expandedQP := make(map[string][]string, len(*h))
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

		expandedQP[key] = expandedValues
	}

	return expandedQP
}
