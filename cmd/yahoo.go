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
			imageColor, _ := c.Flags().GetString("image-color")
			imageSize, _ := c.Flags().GetString("image-size")
			imageType, _ := c.Flags().GetString("image-type")
			safeSearch, _ := c.Flags().GetString("safe-search")

			config := &services.YahooConfig{
				Compact:    compact,
				ImageColor: imageColor,
				ImageSize:  imageSize,
				ImageType:  imageType,
				Query:      args[0],
				SafeSearch: safeSearch,
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
		"image-color",
		"",
		makeFlagUsageMessage(
			"Find images in your preferred color",
			params.ParamAll,
			params.ColorBlackAndWhite,
			params.ColorRed,
			params.ColorOrange,
			params.ColorYellow,
			params.ColorGreen,
			params.ColorTeal,
			params.ColorBlue,
			params.ColorPurple,
			params.ColorPink,
			params.ColorWhite,
			params.ColorGray,
			params.ColorBlack,
			params.ColorBrown,
		),
	)
	yahooCmd.Flags().String(
		"image-size",
		"",
		makeFlagUsageMessage(
			"Find images in specific size",
			params.ParamAll,
			params.ImageSizeSmall,
			params.ImageSizeMedium,
			params.ImageSizeLarge,
		),
	)
	yahooCmd.Flags().String(
		"image-type",
		"",
		makeFlagUsageMessage(
			"Limit the kind of images that you find",
			params.ParamAll,
			params.ImageTypePhoto,
			params.ImageTypeGraphic,
			params.ImageTypeAnimated,
			params.ImageTypeFace,
			params.OrientationPortrait,
			params.ImageTypeNonPortrait,
			params.ImageTypeClipArt,
			params.ImageTypeLineDrawing,
		),
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
