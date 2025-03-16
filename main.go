package main

import (
	"os"

	"github.com/santhanuv/srotas/cmd"
	"github.com/santhanuv/srotas/internal/log"
)

func main() {
	logger := log.New(os.Stderr, os.Stderr, os.Stderr)

	rootCmd := cmd.NewRootCmd(logger, os.Stdin, os.Stdout)

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal("%v", err)
	}
}
