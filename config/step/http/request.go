package http

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/santhanuv/srotas/contract"
	"gopkg.in/yaml.v3"
)

type Request struct {
	Type        string
	Name        string
	Description string
	Url         string
	Method      string
	Body        *RequestBody
	Headers     Header
	QueryParams QueryParam `yaml:"query_params"`
	Timeout     uint
}

type RequestBody struct {
	File string
	Data map[string]string
}

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

func (r *Request) Execute(context contract.ExecutionContext) error {
	req, err := http.NewRequest(r.Method, r.Url, nil)

	if err != nil {
		return err
	}

	for hk, hv := range r.Headers {
		for _, v := range hv {
			req.Header.Add(hk, v)
		}
	}

	urlValues := url.Values{}
	for qp, qv := range r.QueryParams {
		for _, v := range qv {
			urlValues.Add(qp, v)
		}
	}
	req.URL.RawQuery = urlValues.Encode()

	res, err := context.HttpClient().Do(req)

	if err != nil {
		return err
	}

	resBody, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return err
	}

	log.Printf("Response: %v", string(resBody))

	return nil
}
