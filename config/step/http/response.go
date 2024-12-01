package http

import (
	"fmt"
	"log"

	"github.com/santhanuv/srotas/contract"
	"github.com/tidwall/gjson"
)

func storeFromJsonResponse(body []byte, variableSelectorMap map[string]string, context contract.ExecutionContext) error {
	if variableSelectorMap == nil {
		return nil
	}

	selectors := make([]string, 0, len(variableSelectorMap))
	varNames := make([]string, 0, len(variableSelectorMap))

	for varName, selector := range variableSelectorMap {
		selectors = append(selectors, selector)
		varNames = append(varNames, varName)
	}

	if ok := gjson.Valid(string(body)); !ok {
		return fmt.Errorf("Error: Invalid json response")
	}

	queryVal := gjson.GetManyBytes(body, selectors...)

	store := context.Store()

	for idx, qv := range queryVal {
		val := qv.Value()

		if val == nil {
			log.Printf("Warning: Setting nil value for %s", varNames[idx])
		}

		store.Set(varNames[idx], val)
	}

	return nil
}
