package cmd

import (
	"github.com/mzbaulhaque/gois/internal/util"
	"github.com/mzbaulhaque/gois/pkg/scraper"

	"github.com/spf13/cobra"
)

func newBingCmd() *cobra.Command {
	bingCmd := &cobra.Command{
		Use:   "bing [flags] <query>",
		Short: "Search images using bing",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			aspectRatio, _ := c.Flags().GetString("aspect-ratio")
			compact, _ := c.Flags().GetBool("compact")
			date, _ := c.Flags().GetString("date")
			imageColor, _ := c.Flags().GetString("image-color")
			imageSize, _ := c.Flags().GetString("image-size")
			imageType, _ := c.Flags().GetString("image-type")
			license, _ := c.Flags().GetString("license")
			people, _ := c.Flags().GetString("people")
			safeSearch, _ := c.Flags().GetString("region")

			config := &scraper.BingConfig{
				AspectRatio: aspectRatio,
				Compact:  compact,
				Date: date,
				ImageColor: imageColor,
				ImageSize: imageSize,
				ImageType: imageType,
				License: license,
				People: people,
				Query:    args[0],
				SafeSearch: safeSearch,
			}
			bs := &scraper.BingScraper{Config: config}
			items, err := bs.Scrape()

			if err != nil {
				return err
			}

			util.PrintResults(items, config.Compact)

			return nil
		},
	}

	bingCmd.Flags().StringP("aspect-ratio", "A", "", "Specify the shape of images [any (default), square, wide, tall]")
	bingCmd.PersistentFlags().BoolP("compact", "c", false, "Print original image link per line with no other information.")
	bingCmd.Flags().StringP("date", "D", "", "Specify date of images [any (default), past-day, past-week, past-month, past-year]")
	bingCmd.Flags().StringP("image-color", "C", "", "Find images in preferred color [any (default), full-color, black-white, red, orange, yellow, green, teal, blue, purple, pink, white, gray, black, brown]")
	bingCmd.Flags().StringP("image-size", "S", "", "Find images in specific size [any (default), small, medium, large, extra-large or specific size e.g. 300_300]")
	bingCmd.Flags().StringP("image-type", "T", "", "Limit the kind of images that you find [any (default), photo, clip-art, line-drawing, animated, transparent]")
	bingCmd.Flags().StringP("license", "L", "", "License preference [all (default), creative-commons, public-domain, free-share-use, free-share-use-commercially, free-modify-share-use, free-modify-share-use-commercially]")
	bingCmd.Flags().StringP("people", "P", "", "Apply people filter [any (default), face, head-shoulder]")
	bingCmd.Flags().StringP("safe-search", "s", "on", "Tell SafeSearch whether to filter sexually explicit content [on (default), off, moderate]")

	return bingCmd
}
