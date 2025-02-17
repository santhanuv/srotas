package main

import (
	"io"
	"os"

	"github.com/santhanuv/srotas/cmd"
	"github.com/santhanuv/srotas/internal/log"
)

func main() {
	logger := log.New(os.Stderr, io.Discard, os.Stderr)

	rootCmd := cmd.NewRootCmd(logger, os.Stdin, os.Stdout)

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal("%v", err)
	}
}
