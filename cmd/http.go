package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/santhanuv/srotas/internal/http"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&httpCommand)

	httpCommand.Flags().StringArrayP("query", "Q", []string{},
		"Define query parameters as 'key=value'. Multiple parameters can be specified using commas.")

	httpCommand.Flags().StringArrayP("header", "H", []string{},
		"Add request headers in 'key:value' format. Multiple headers can be specified using commas.")

	httpCommand.Flags().StringP("body", "B", "",
		"Provide a request body. Only JSON is supported.")
}

var httpCommand = cobra.Command{

	Use:   "http [METHOD] [URL]",
	Short: "Send an HTTP request to a specified URL.",
	Long: `Send an HTTP request using the specified METHOD:

	- METHOD: The HTTP method to use (GET, POST, PUT, DELETE, etc.).
	- URL: The target URL for the request.
		
Optional flags allow you to add query parameters, headers, and a request body.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		method, rawURL := args[0], args[1]
		method = strings.ToUpper(method)

		rawQueryParams, err := cmd.Flags().GetStringArray("query")

		if err != nil {
			log.Fatalf("error on parsing query params: %v", err)
		}

		queryParams, err := parseQueryParams(rawQueryParams)

		if err != nil {
			log.Fatalf("error on parsing query params: %v", err)
		}

		rawHeaders, err := cmd.Flags().GetStringArray("header")

		if err != nil {
			log.Fatalf("error on parsing headers: %v", err)
		}

		headers, err := parseHeaders(rawHeaders)

		if _, ok := headers["Content-Type"]; !ok {
			headers["Content-Type"] = []string{"application/json"}
		}

		if err != nil {
			log.Fatalf("error on parsing header: %v", err)
		}

		rawRequestBody, err := cmd.Flags().GetString("body")

		if err != nil {
			log.Fatalf("error on parsing request body: %v", err)
		}

		c := http.NewClient(0)

		req := &http.Request{
			Method:      method,
			Url:         rawURL,
			Headers:     headers,
			QueryParams: queryParams,
			Body:        []byte(rawRequestBody),
		}

		res, err := c.Do(req)

		if err != nil {
			log.Fatalf("sending http request failed: %v", err)
		}

		var responseJson bytes.Buffer
		json.Indent(&responseJson, res.Body, "", " ")
		if err != nil {
			log.Fatalf("failed to parse response: %s", err)
		}

		log.Printf("Response:\n%s", responseJson.String())
	},
}

func parseQueryParams(rawQueryParams []string) (map[string][]string, error) {
	queryParams := make(map[string][]string)

	for _, rqp := range rawQueryParams {
		rqp = strings.TrimSpace(rqp)
		pairs := strings.Split(rqp, "=")

		if len(pairs) < 2 {
			return nil, fmt.Errorf("Invalid query parameter: %s", rqp)
		}

		key, value := pairs[0], pairs[1]
		queryParams[key] = append(queryParams[key], value)
	}

	return queryParams, nil
}

func parseHeaders(rawHeaders []string) (map[string][]string, error) {
	headers := make(map[string][]string)

	for _, rh := range rawHeaders {
		rh = strings.TrimSpace(rh)
		pairs := strings.Split(rh, ":")

		if len(pairs) < 2 {
			return nil, fmt.Errorf("Invalid header: %s", rh)
		}

		key, value := pairs[0], pairs[1]
		headers[key] = append(headers[key], value)
	}

	return headers, nil
}
