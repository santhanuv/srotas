package http

import (
	"io"
	"log"
	"os"

	"github.com/santhanuv/srotas/contract"
	"github.com/tidwall/sjson"
	"gopkg.in/yaml.v3"
)

type RequestBody struct {
	Content []byte
	Data    map[string]string
}

func (rb *RequestBody) UnmarshalYAML(value *yaml.Node) error {
	var rawRequestBody struct {
		File string
		Data map[string]string
	}

	if err := value.Decode(&rawRequestBody); err != nil {
		return err
	}

	*rb = RequestBody{
		Data: rawRequestBody.Data,
	}

	if rawRequestBody.File != "" {
		file, err := os.Open(rawRequestBody.File)
		defer file.Close()

		if err != nil {
			return err
		}

		content, err := io.ReadAll(file)

		if err != nil {
			return err
		}

		rb.Content = content
	}

	return nil
}

// build merges rb.Data with rb.Content and returns the result.
func (rb *RequestBody) build(context contract.ExecutionContext) ([]byte, error) {
	store := context.Store()

	var (
		updatedContent []byte
		err            error
	)

	for field, variable := range rb.Data {
		value, ok := store.Get(variable)

		if !ok {
			log.Printf("varable:'%s' not found in store\n", variable)
			continue
		}

		updatedContent, err = sjson.SetBytes(rb.Content, field, value)

		if err != nil {
			return nil, err
		}
	}

	return updatedContent, nil
}
