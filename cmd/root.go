package cmd

import (
	"io"
	"os"

	"github.com/santhanuv/srotas/internal/config"
	"github.com/santhanuv/srotas/internal/log"
	"github.com/spf13/cobra"
)

func NewRootCmd(logger *log.Logger, in *os.File, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "srotas",
		Short: "Srotas is a cli for testing api",
		Long:  "Srotas is a flexible cli tool for testing api with different flows",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}

			return nil
		},
	}

	cr := config.NewConfigRunner()
	cmd.AddCommand(newRunCommand(logger, in, out, cr))
	cmd.AddCommand(newHttpCommand(out))

	return cmd
}
