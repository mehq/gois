package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main().
func Execute() {
	c := newRootCmd()

	cobra.AddTemplateFunc("platform", func() string { return fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH) })
	c.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "version %s" .Version}} ({{platform}})
`)
	c.SetHelpCommand(&cobra.Command{Hidden: true})

	c.AddCommand(newBingCmd(), newFlickrCmd(), newGoogleCmd())

	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}

func buildFlagUsageMessage(args ...string) string {
	return fmt.Sprintf("%s [%s (default), %s]", args[0], args[1], strings.Join(args[2:], ", "))
}
