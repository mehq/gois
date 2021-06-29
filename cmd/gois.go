package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mzbaulhaque/gois/internal/build"
)

var rootCmd = &cobra.Command{
	Use:     build.ProjectName,
	Short:   "Command line program to search images",
	Args:    cobra.NoArgs,
	Version: build.Version,
	Run:     func(c *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main().
func Execute() error {
	if rootCmd.Execute() != nil {
		return fmt.Errorf("failed to execute root command")
	}

	return nil
}

func init() {
	cobra.AddTemplateFunc("platform", func() string { return fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH) })
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "version %s" .Version}} ({{platform}})
`)
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	err := viper.BindPFlags(rootCmd.PersistentFlags())

	if err != nil {
		panic(err)
	}

	bingCmd.Flags().BoolP("compact", "c", false, "Print original image link per line with no other information.")
	bingCmd.Flags().IntP("height", "H", 0, "Download images with given height (width must be provided as well).")
	bingCmd.Flags().IntP("width", "w", 0, "Download images with given width (height must be provided as well).")
	bingCmd.Flags().BoolP("explicit", "x", false, "Turn safe search off.")
	bingCmd.Flags().BoolP("gif", "g", false, "Download only gif images.")
	bingCmd.Flags().BoolP("gray", "B", false, "Download only black and white images.")
	rootCmd.AddCommand(bingCmd)

	googleCmd.Flags().BoolP("compact", "c", false, "Print original image link per line with no other information.")
	googleCmd.Flags().IntP("height", "H", 0, "Download images with given height (width must be provided as well).")
	googleCmd.Flags().IntP("width", "w", 0, "Download images with given width (height must be provided as well).")
	googleCmd.Flags().BoolP("explicit", "x", false, "Turn safe search off.")
	googleCmd.Flags().BoolP("gif", "g", false, "Download only gif images.")
	googleCmd.Flags().BoolP("gray", "B", false, "Download only black and white images.")
	rootCmd.AddCommand(googleCmd)
}

func getBoolValue(c *cobra.Command, name string) bool {
	val, err := c.Flags().GetBool(name)

	if err != nil {
		return false
	}

	return val
}

func getIntValue(c *cobra.Command, name string) int {
	val, err := c.Flags().GetInt(name)

	if err != nil {
		return 0
	}

	return val
}
