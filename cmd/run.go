package cmd

import (
	"io"
	"os"
	"path/filepath"

	"github.com/santhanuv/srotas/internal/log"
	"github.com/santhanuv/srotas/workflow"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&runCommand)
	runCommand.Flags().BoolP("verbose", "v", false, "Enable verbose mode to display detailed logs about the execution of the config.")
}

var runCommand = cobra.Command{
	Use:   "run [CONFIG]",
	Short: "Run the provided configuration.",
	Long:  "Runs the provided configuration file. The configuration can be provided as a yaml file.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.Logger{}
		logger.SetInfoWriter(os.Stderr)
		logger.SetErrorWriter(os.Stderr)

		verboseMode, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			logger.Fatal("verbose: %v", err)
		}
		if verboseMode {
			logger.SetDebugWriter(os.Stderr)
		} else {
			logger.SetDebugWriter(io.Discard)
		}

		configPath := args[0]
		if configPath == "" {
			logger.Fatal("Config: Invalid configuration file")
		}

		configPath, err = filepath.Abs(configPath)
		if err != nil {
			logger.Fatal("Config: %v", err)
		}

		logger.Debug("Parsing %s", configPath)
		flowDef, err := workflow.ParseConfig(configPath, &logger)
		if err != nil {
			logger.Fatal("Parse error: %v", err)
		}
		logger.Debug("Successfully parsed %s", configPath)

		logger.Debug("Initializing execution context")
		execCtx, err := workflow.NewExecutionContext(
			workflow.WithGlobalOptions(flowDef.BaseUrl, flowDef.Headers),
			workflow.WithLogger(&logger))
		if err != nil {
			logger.Fatal("Execution context error: %v", err)
		}
		logger.Debug("Successfully initialized execution context")

		logger.Debug("Executing configuration")
		err = workflow.Execute(flowDef, execCtx)
		if err != nil {
			logger.Fatal("Execution error: %v", err)
		}
		logger.Debug("Successfully executed configuration")
	},
}
