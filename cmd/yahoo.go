package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mzbaulhaque/gois/internal/util"
	"github.com/mzbaulhaque/gois/pkg/scraper/params"
	"github.com/mzbaulhaque/gois/pkg/scraper/services"
)

func newYahooCmd() *cobra.Command {
	yahooCmd := &cobra.Command{
		Use:   "yahoo [flags] <query>",
		Short: "Search images using yahoo",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			compact, _ := c.Flags().GetBool("compact")
			safeSearch, _ := c.Flags().GetString("safe-search")

			config := &services.YahooConfig{
				Compact:     compact,
				Query:       args[0],
				SafeSearch:  safeSearch,
			}
			fs := &services.YahooScraper{Config: config}
			items, pages, err := fs.Scrape()

			if err != nil {
				return fmt.Errorf("%v", err)
			}

			util.PrintResults(items, pages, config.Compact)

			return nil
		},
	}

	yahooCmd.PersistentFlags().BoolP(
		"compact",
		"c",
		false,
		"Print original image link per line with no other information.",
	)
	yahooCmd.Flags().String(
		"safe-search",
		"",
		makeFlagUsageMessage(
			"Tell SafeSearch whether to filter sexually explicit content",
			params.SafeSearchOn,
			params.SafeSearchOff,
		),
	)

	return yahooCmd
}
