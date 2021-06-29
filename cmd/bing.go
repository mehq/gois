package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mzbaulhaque/gois/internal/util/out"
	"github.com/mzbaulhaque/gois/pkg/scraper/bing"
)

var bingCmd = &cobra.Command{
	Use:   "bing",
	Short: "Search images using bing",
	Args:  cobra.ExactArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		config := &bing.Config{
			Compact:  getBoolValue(c, "compact"),
			Explicit: getBoolValue(c, "explicit"),
			GIF:      getBoolValue(c, "gif"),
			Gray:     getBoolValue(c, "gray"),
			Query:    args[0],
			Height:   getIntValue(c, "height"),
			Width:    getIntValue(c, "width"),
		}
		scraper := &bing.Bing{Config: config}
		items, err := scraper.Scrape()

		if err != nil {
			return fmt.Errorf("failed to scrape bing")
		}

		err = out.PrintImageInfo(items, config.Compact)

		if err != nil {
			return fmt.Errorf("failed to print bing image info")
		}

		return nil
	},
}
