package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/santhanuv/srotas/internal/http"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&httpCommand)
	httpCommand.Flags().StringSliceP("query", "q", []string{}, "Specify query parameters seperated by comma if any")
	httpCommand.Flags().StringSliceP("headers", "H", []string{}, "Specify headers seperated by comma if any")
	httpCommand.Flags().StringP("body", "B", "", "Optionally specify the data to be send with request as the request body. When specifying body Content-Type header should be set, otherwise it is send as text/plain. Currently only json is supported")
}

var httpCommand = cobra.Command{
	Use:   "http [METHOD] [URL]",
	Short: "Sends http METHOD request to the specified URL",
	Long: `Sends http METHOD request to the specified URL:
	METHOD can be any http request methods like GET, POST, PUT, DELETE,... .
	
	URL specifies the url to send request to.
	`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		method, rawURL := args[0], args[1]
		method = strings.ToUpper(method)

		rawQueryParams, err := cmd.Flags().GetStringSlice("query")

		if err != nil {
			log.Fatalf("query param: %v", err)
		}

		queryParams, err := parseQueryParams(rawQueryParams)

		if err != nil {
			log.Fatalf("query param: %v", err)
		}

		rawHeaders, err := cmd.Flags().GetStringSlice("headers")

		if err != nil {
			log.Fatalf("header: %v", err)
		}

		headers, err := parseHeaders(rawHeaders)

		if err != nil {
			log.Fatalf("header: %v", err)
		}

		rawRequestBody, err := cmd.Flags().GetString("body")

		if err != nil {
			log.Fatalf("body: %v", err)
		}

		c := http.NewClient(0, nil)
		req := &http.Request{
			Method:      method,
			Url:         rawURL,
			Headers:     headers,
			QueryParams: queryParams,
			Body:        []byte(rawRequestBody),
		}

		res, err := c.Do(req)

		if err != nil {
			log.Fatalf("%s", err)
		}

		responseJson, err := json.Marshal(*res)

		if err != nil {
			log.Fatalf("%s", err)
		}

		fmt.Println(string(responseJson))
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
