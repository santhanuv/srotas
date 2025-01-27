package cmd

import (
	"os"
	"path/filepath"

	"github.com/santhanuv/srotas/internal/log"
	"github.com/santhanuv/srotas/workflow"
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
		logger := log.New(os.Stderr, os.Stderr, os.Stderr)

		if configPath == "" {
			logger.Fatal("Config: Invalid configuration file")
		}

		configPath, err := filepath.Abs(configPath)

		if err != nil {
			logger.Fatal("Config: %v", err)
		}

		flowDef, err := workflow.ParseConfig(configPath)

		if err != nil {
			logger.Fatal("Parse error: %v", err)
		}

		execCtx, err := workflow.NewExecutionContext(
			workflow.WithGlobalOptions(flowDef.BaseUrl, flowDef.Headers),
			workflow.WithLogger(logger))

		if err != nil {
			logger.Fatal("Execution context error: %v", err)
		}

		err = workflow.Execute(flowDef, execCtx)

		if err != nil {
			logger.Fatal("Execution error: %v", err)
		}
	},
}
