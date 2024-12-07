package http

import (
	"fmt"
	"log"

	"github.com/santhanuv/srotas/contract"
	"github.com/tidwall/gjson"
)

func storeFromResponse(body []byte, query map[string]string, context contract.ExecutionContext) error {
	if query == nil {
		return nil
	}

	selectors := make([]string, 0, len(query))
	variables := make([]string, 0, len(query))

	for variable, selector := range query {
		selectors = append(selectors, selector)
		variables = append(variables, variable)
	}

	if ok := gjson.Valid(string(body)); !ok {
		return fmt.Errorf("Error: Invalid json response")
	}

	queryVal := gjson.GetManyBytes(body, selectors...)

	store := context.Store()

	for idx, qv := range queryVal {
		val := qv.Value()

		if val == nil {
			log.Printf("Warning: Setting nil value for %s", variables[idx])
		}

		store.Set(variables[idx], val)
	}

	return nil
}
