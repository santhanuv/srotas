package validation

import (
	"fmt"
	"strings"

	"github.com/santhanuv/srotas/contract"
	"github.com/santhanuv/srotas/internal/http"
	"github.com/tidwall/gjson"
)

type Assert struct {
	Value    string
	Selector string
}

func (a *Assert) Validate(context contract.ExecutionContext, response *http.Response) error {
	var expected any = a.Value

	if strings.HasPrefix(a.Value, "$") {
		var ok bool
		expected, ok = context.Store().Get(a.Value[1:])

		if !ok {
			return fmt.Errorf("Invalid variable name in assert")
		}
	}

	actual := gjson.GetBytes(response.Body, a.Selector).Value()

	if actual == nil {
		return fmt.Errorf("Assert failed: value not found in response body")
	}

	if expected != actual {
		return fmt.Errorf("Assert failed: expected '%s' but got '%s'", expected, actual)
	}

	return nil
}
