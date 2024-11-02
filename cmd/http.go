package cmd

import (
	"fmt"

	"github.com/santhanuv/srotas/client"
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

		client := client.NewApiClient(method, rawURL)

		client.Do()

		fmt.Println(client.Response)
	},
}
