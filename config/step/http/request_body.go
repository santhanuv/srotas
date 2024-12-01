package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/santhanuv/srotas/contract"
	"gopkg.in/yaml.v3"
)

type ContentType uint

const (
	jsonContent ContentType = iota
	textContent
)

type RequestBody struct {
	Type ContentType
	File []byte
	Data map[string]string
}

func (r *RequestBody) UnmarshalYAML(value *yaml.Node) error {
	var rawRequestBody struct {
		File string
		Data map[string]string
		Type string
	}

	if err := value.Decode(&rawRequestBody); err != nil {
		return err
	}

	var content []byte
	if rawRequestBody.File != "" {
		file, err := os.Open(rawRequestBody.File)
		defer file.Close()

		if err != nil {
			return err
		}

		content, err = io.ReadAll(file)

		if err != nil {
			return err
		}

		r.File = content
	}

	ct, err := parseContentType(rawRequestBody.Type)

	if err != nil {
		return nil
	}

	*r = RequestBody{
		Type: ct,
		File: content,
		Data: rawRequestBody.Data,
	}

	return nil
}

func parseContentType(contentType string) (ContentType, error) {
	switch contentType {
	case "json":
		return jsonContent, nil
	case "":
		return textContent, nil
	default:
		return textContent, fmt.Errorf("%s not supported", contentType)
	}
}

func (rb *RequestBody) transform(context contract.ExecutionContext) (io.ReadCloser, error) {
	switch rb.Type {
	case jsonContent:
		fileReader, err := rb.transformAsJson(context)
		if err != nil {
			return nil, err
		}
		return fileReader, nil
	case textContent:
		return nil, nil
	default:
		err := fmt.Errorf("Content type not supported")
		return nil, err
	}
}

func (rb *RequestBody) transformAsJson(context contract.ExecutionContext) (io.ReadCloser, error) {
	var fileData map[string]any
	if rb.File != nil {
		if err := json.Unmarshal(rb.File, &fileData); err != nil {
			return nil, err
		}
	}

	store := context.Store()
	if rb.Data != nil {
		for field, rawValue := range rb.Data {
			var value any = rawValue

			if strings.HasPrefix(rawValue, "$") {
				var ok bool

				variable := rawValue[1:]
				value, ok = store.Get(variable)

				if !ok {
					log.Printf("%s not found in store\n", variable)
					continue
				}
			}

			fileData[field] = value
		}
	}

	dataBytes, err := json.Marshal(fileData)

	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(dataBytes)
	return io.NopCloser(reader), nil

}
