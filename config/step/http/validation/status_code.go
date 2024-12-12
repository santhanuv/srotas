package validation

import (
	"fmt"

	"github.com/santhanuv/srotas/contract"
	"github.com/santhanuv/srotas/internal/http"
)

type StatusCode uint

func (s *StatusCode) Validate(context contract.ExecutionContext, response *http.Response) error {
	if uint(*s) != response.StatusCode {
		return fmt.Errorf("Status code: Expected %d but got %d", *s, response.StatusCode)
	}

	return nil
}
