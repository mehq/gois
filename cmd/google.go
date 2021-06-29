package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mzbaulhaque/gois/internal/util/out"
	"github.com/mzbaulhaque/gois/pkg/scraper/google"
)

var googleCmd = &cobra.Command{
	Use:   "google",
	Short: "Search images using google",
	Args:  cobra.ExactArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		config := &google.Config{
			Compact:  getBoolValue(c, "compact"),
			Explicit: getBoolValue(c, "explicit"),
			GIF:      getBoolValue(c, "gif"),
			Gray:     getBoolValue(c, "gray"),
			Query:    args[0],
			Height:   getIntValue(c, "height"),
			Width:    getIntValue(c, "width"),
		}
		scraper := &google.Google{Config: config}
		items, err := scraper.Scrape()

		if err != nil {
			return fmt.Errorf("failed to scrape google")
		}

		err = out.PrintImageInfo(items, config.Compact)

		if err != nil {
			return fmt.Errorf("failed to print google image info")
		}

		return nil
	},
}
