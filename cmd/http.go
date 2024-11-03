package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/santhanuv/srotas/internal/client"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&httpCommand)
	httpCommand.Flags().StringSliceP("query", "q", []string{}, "Specify query parameters seperated by comma if any")
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
			fmt.Fprintf(os.Stderr, "%s", err)
		}

		queryParams, err := ParseQueryParams(rawQueryParams)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
		}

		c := client.NewClient()
		req := client.NewRequest(method, rawURL, nil)
		req.SetQueryParams(queryParams)

		res, err := c.Do(*req)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
		}

		responseJson, err := json.Marshal(*res)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
		}

		fmt.Println(string(responseJson))
	},
}

func ParseQueryParams(rawQueryParams []string) (map[string][]string, error) {
	queryParams := make(map[string][]string)

	for _, rqp := range rawQueryParams {
		pairs := strings.Split(rqp, "=")

		if len(pairs) < 2 {
			return nil, fmt.Errorf("Invalid query parameter: %s", rqp)
		}

		key, value := pairs[0], pairs[1]
		queryParams[key] = append(queryParams[key], value)
	}

	return queryParams, nil
}
