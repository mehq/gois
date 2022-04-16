package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mehq/gois/internal/util"
	"github.com/mehq/gois/pkg/scraper/params"
	"github.com/mehq/gois/pkg/scraper/services"
)

func newGoogleCmd() *cobra.Command {
	googleCmd := &cobra.Command{
		Use:   "google [flags] <query>",
		Short: "Search images using google",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			aspectRatio, _ := c.Flags().GetString("aspect-ratio")
			compact, _ := c.Flags().GetBool("compact")
			imageColor, _ := c.Flags().GetString("image-color")
			imageSize, _ := c.Flags().GetString("image-size")
			imageType, _ := c.Flags().GetString("image-type")
			safeSearch, _ := c.Flags().GetString("safe-search")

			config := &services.GoogleConfig{
				AspectRatio: aspectRatio,
				Compact:     compact,
				ImageColor:  imageColor,
				ImageSize:   imageSize,
				ImageType:   imageType,
				Query:       args[0],
				SafeSearch:  safeSearch,
			}
			gs := &services.GoogleScraper{Config: config}
			items, pages, err := gs.Scrape()

			if err != nil {
				return fmt.Errorf("%v", err)
			}

			util.PrintResults(items, pages, config.Compact)

			return nil
		},
	}

	googleCmd.Flags().String(
		"aspect-ratio",
		"",
		makeFlagUsageMessage(
			"Specify the shape of images",
			params.ParamAll,
			params.AspectRatioTall,
			params.AspectRatioSquare,
			params.AspectRationWide,
			params.AspectRatioPanoramic,
		),
	)
	googleCmd.PersistentFlags().BoolP(
		"compact",
		"c",
		false,
		"Print original image link per line with no other information.",
	)
	googleCmd.Flags().String(
		"image-color",
		"",
		makeFlagUsageMessage(
			"Find images in your preferred color",
			params.ParamAll,
			params.ColorFull,
			params.ColorBlackAndWhite,
			params.ImageTypeTransparent,
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
	googleCmd.Flags().String(
		"image-size",
		"",
		makeFlagUsageMessage(
			"Find images in specific size",
			params.ParamAll,
			params.ImageSizeLarge,
			params.ImageSizeMedium,
			params.ImageSizeIcon,
		),
	)
	googleCmd.Flags().String(
		"image-type",
		"",
		makeFlagUsageMessage(
			"Limit the kind of images that you find",
			params.ParamAll,
			params.ImageTypeFace,
			params.ImageTypePhoto,
			params.ImageTypeClipArt,
			params.ImageTypeLineDrawing,
			params.ImageTypeAnimated,
		),
	)
	googleCmd.Flags().String(
		"safe-search",
		"",
		makeFlagUsageMessage(
			"Tell SafeSearch whether to filter sexually explicit content",
			params.SafeSearchOn,
			params.SafeSearchOff,
		),
	)

	return googleCmd
}
