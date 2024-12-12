package validation

import (
	"github.com/santhanuv/srotas/contract"
	"github.com/santhanuv/srotas/internal/http"
)

type Validator struct {
	Asserts    []Assert
	StatusCode StatusCode `yaml:"status_code"`
}

func (v *Validator) Validate(context contract.ExecutionContext, response *http.Response) error {
	if v.StatusCode > 0 {
		err := v.StatusCode.Validate(context, response)

		if err != nil {
			return err
		}
	}

	for _, assert := range v.Asserts {
		err := assert.Validate(context, response)

		if err != nil {
			return err
		}
	}

	return nil
}
