package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/santhanuv/srotas/internal/client"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&httpCommand)
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
		c := client.NewClient()
		req := client.NewRequest(method, rawURL, nil)

		res, err := c.Do(*req)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
		}

		fmt.Println(*res)
	},
}
