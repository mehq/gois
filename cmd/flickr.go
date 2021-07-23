package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mzbaulhaque/gois/internal/util"
	"github.com/mzbaulhaque/gois/pkg/scraper/params"
	"github.com/mzbaulhaque/gois/pkg/scraper/services"
)

func newFlickrCmd() *cobra.Command {
	flickrCmd := &cobra.Command{
		Use:   "flickr [flags] <query>",
		Short: "Search images using flickr",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			compact, _ := c.Flags().GetBool("compact")
			imageColor, _ := c.Flags().GetString("image-color")
			imageSize, _ := c.Flags().GetString("image-size")
			imageType, _ := c.Flags().GetString("image-type")
			orientation, _ := c.Flags().GetString("orientation")
			safeSearch, _ := c.Flags().GetString("safe-search")

			config := &services.FlickrConfig{
				Compact:     compact,
				ImageColor:  imageColor,
				ImageSize:   imageSize,
				ImageType:   imageType,
				Orientation: orientation,
				Query:       args[0],
				SafeSearch:  safeSearch,
			}
			fs := &services.FlickrScraper{Config: config}
			items, pages, err := fs.Scrape()

			if err != nil {
				return fmt.Errorf("%v", err)
			}

			util.PrintResults(items, pages, config.Compact)

			return nil
		},
	}

	flickrCmd.PersistentFlags().BoolP(
		"compact",
		"c",
		false,
		"Print original image link per line with no other information.",
	)
	flickrCmd.Flags().String(
		"image-color",
		"",
		makeFlagUsageMessage(
			"Find images in your preferred color",
			params.ParamAll,
			params.ColorRed,
			params.ColorDarkOrange,
			params.ColorOrange,
			params.ColorPalePink,
			params.ColorLemonYellow,
			params.ColorSchoolBusYellow,
			params.ColorGreen, params.ColorDarkLimeGreen,
			params.ColorCyan,
			params.ColorBlue,
			params.ColorViolet,
			params.ColorPink,
			params.ColorWhite,
			params.ColorGray,
			params.ColorBlack,
		),
	)
	flickrCmd.Flags().String(
		"image-size",
		"",
		makeFlagUsageMessage(
			"Find images in specific size",
			params.ParamAll,
			params.ImageSizeLarge,
			params.ImageSizeMedium,
		),
	)
	flickrCmd.Flags().String(
		"image-type",
		"",
		makeFlagUsageMessage(
			"Limit the kind of images that you find",
			params.ParamAll,
			params.ColorBlackAndWhite,
			params.ImageTypeShallowDepthOfField,
			params.ImageTypeMinimal,
			params.ImageTypePatterns,
		),
	)
	flickrCmd.Flags().String(
		"orientation",
		"",
		makeFlagUsageMessage(
			"Specify the orientation of images",
			params.ParamAll,
			params.OrientationLandscape,
			params.OrientationPortrait,
			params.AspectRatioSquare,
			params.AspectRatioPanoramic,
		),
	)
	flickrCmd.Flags().String(
		"safe-search",
		"",
		makeFlagUsageMessage(
			"Tell SafeSearch whether to filter sexually explicit content",
			params.SafeSearchOn,
			params.SafeSearchOff,
			params.SafeSearchModerate,
		),
	)

	return flickrCmd
}
