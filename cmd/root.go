package cmd

import (
	"github.com/mzbaulhaque/gois/internal/build"
	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     build.ProjectName,
		Short:   "Command line program to search images",
		Args:    cobra.NoArgs,
		Version: build.Version,
		Run:     func(c *cobra.Command, args []string) {
			_ = c.Usage()
		},
		SilenceUsage: true,
	}

	return rootCmd
}
