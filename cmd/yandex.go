package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mzbaulhaque/gois/internal/util"
	"github.com/mzbaulhaque/gois/pkg/scraper/services"
)

func newYandexCmd() *cobra.Command {
	yandexCmd := &cobra.Command{
		Use:   "yandex [flags] <query>",
		Short: "Search images using yandex",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			compact, _ := c.Flags().GetBool("compact")

			config := &services.YandexConfig{
				Compact: compact,
				Query:   args[0],
			}
			fs := &services.YandexScraper{Config: config}
			items, pages, err := fs.Scrape()

			if err != nil {
				return fmt.Errorf("%v", err)
			}

			util.PrintResults(items, pages, config.Compact)

			return nil
		},
	}

	yandexCmd.PersistentFlags().BoolP(
		"compact",
		"c",
		false,
		"Print original image link per line with no other information.",
	)

	return yandexCmd
}
