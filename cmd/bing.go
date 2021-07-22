package cmd

import (
	"fmt"

	"github.com/mzbaulhaque/gois/internal/util"
	"github.com/mzbaulhaque/gois/pkg/scraper/params"
	"github.com/mzbaulhaque/gois/pkg/scraper/services"

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
			imageColor, _ := c.Flags().GetString("image-color")
			imageSize, _ := c.Flags().GetString("image-size")
			imageType, _ := c.Flags().GetString("image-type")
			peopleFilter, _ := c.Flags().GetString("people-filter")
			safeSearch, _ := c.Flags().GetString("safe-search")

			config := &services.BingConfig{
				AspectRatio:  aspectRatio,
				Compact:      compact,
				ImageColor:   imageColor,
				ImageSize:    imageSize,
				ImageType:    imageType,
				PeopleFilter: peopleFilter,
				Query:        args[0],
				SafeSearch:   safeSearch,
			}
			bs := &services.BingScraper{Config: config}
			items, pages, err := bs.Scrape()

			if err != nil {
				return fmt.Errorf("%v", err)
			}

			util.PrintResults(items, pages, config.Compact)

			return nil
		},
	}

	bingCmd.Flags().String(
		"aspect-ratio",
		"",
		buildFlagUsageMessage(
			"Specify the shape of images",
			"all",
			params.AspectRatioTall,
			params.AspectRatioSquare,
			params.AspectRationWide,
		),
	)
	bingCmd.PersistentFlags().BoolP(
		"compact",
		"c",
		false,
		"Print original image link per line with no other information.",
	)
	bingCmd.Flags().String(
		"image-color",
		"",
		buildFlagUsageMessage(
			"Find images in your preferred color",
			"all",
			params.ColorFull,
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
	bingCmd.Flags().String(
		"image-size",
		"",
		buildFlagUsageMessage(
			"Find images in specific size",
			"all",
			params.ImageSizeLarge,
			params.ImageSizeMedium,
			params.ImageSizeSmall,
			params.ImageSizeExtraLarge,
		),
	)
	bingCmd.Flags().String(
		"image-type",
		"",
		buildFlagUsageMessage(
			"Limit the kind of images that you find",
			"all",
			params.ImageTypePhoto,
			params.ImageTypeClipArt,
			params.ImageTypeLineDrawing,
			params.ImageTypeAnimated,
			params.ImageTypeTransparent,
		),
	)
	bingCmd.Flags().String(
		"people-filter",
		"",
		buildFlagUsageMessage(
			"Apply people filter",
			"all",
			params.ImageTypeFace,
			params.OrientationPortrait,
		),
	)
	bingCmd.Flags().String(
		"safe-search",
		"on",
		buildFlagUsageMessage(
			"Tell SafeSearch whether to filter sexually explicit content",
			params.SafeSearchOn,
			params.SafeSearchOff,
			params.SafeSearchModerate,
		),
	)

	return bingCmd
}
