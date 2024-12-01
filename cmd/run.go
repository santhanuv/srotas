package cmd

import (
	"log"
	"path/filepath"

	"github.com/santhanuv/srotas/internal/executor"
	"github.com/santhanuv/srotas/internal/parser"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&runCommand)
}

var runCommand = cobra.Command{
	Use:   "run [CONFIG]",
	Short: "Run the provided configuration.",
	Long:  "Runs the provided configuration file. The configuration can be provided as a yaml file.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configPath := args[0]

		if configPath == "" {
			log.Fatalf("Config: Invalid configuration file")
		}

		configPath, err := filepath.Abs(configPath)

		if err != nil {
			log.Fatalf("Config: %v", err)
		}

		flowDef, err := parser.ParseConfig(configPath)

		if err != nil {
			log.Fatalf("Parse error: %v", err)
		}

		err = executor.Execute(flowDef)

		if err != nil {
			log.Fatalf("Execution error: %v", err)
		}
	},
}
