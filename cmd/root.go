package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "srotas",
	Short: "Srotas is a cli for testing api",
	Long:  "Srotas is a flexible cli tool for testing api with different flows",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("%s", err)
	}
}
