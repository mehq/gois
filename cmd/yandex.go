package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mzbaulhaque/gois/internal/util"
	"github.com/mzbaulhaque/gois/pkg/scraper/params"
	"github.com/mzbaulhaque/gois/pkg/scraper/services"
)

func newYandexCmd() *cobra.Command {
	yandexCmd := &cobra.Command{
		Use:   "yandex [flags] <query>",
		Short: "Search images using yandex",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			compact, _ := c.Flags().GetBool("compact")
			imageColor, _ := c.Flags().GetString("image-color")
			imageSize, _ := c.Flags().GetString("image-size")
			imageType, _ := c.Flags().GetString("image-type")
			orientation, _ := c.Flags().GetString("orientation")
			safeSearch, _ := c.Flags().GetString("safe-search")

			config := &services.YandexConfig{
				Compact:     compact,
				ImageColor:  imageColor,
				ImageSize:   imageSize,
				ImageType:   imageType,
				Orientation: orientation,
				Query:       args[0],
				SafeSearch:  safeSearch,
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
	yandexCmd.Flags().String(
		"image-color",
		"",
		makeFlagUsageMessage(
			"Find images in your preferred color",
			params.ParamAll,
			params.ColorFull,
			params.ColorBlackAndWhite,
			params.ColorRed,
			params.ColorOrange,
			params.ColorYellow,
			params.ColorCyan,
			params.ColorGreen,
			params.ColorBlue,
			params.ColorViolet,
			params.ColorWhite,
			params.ColorBlack,
		),
	)
	yandexCmd.Flags().String(
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
	yandexCmd.Flags().String(
		"image-type",
		"",
		makeFlagUsageMessage(
			"Limit the kind of images that you find",
			params.ParamAll,
			params.ImageTypePhoto,
			params.ImageTypeLineDrawing,
			params.ImageTypeFace,
			params.ImageTypeClipArt,
			params.ImageTypeDemotivational,
		),
	)
	yandexCmd.Flags().String(
		"orientation",
		"",
		makeFlagUsageMessage(
			"Specify the orientation of images",
			params.ParamAll,
			params.OrientationLandscape,
			params.OrientationPortrait,
			params.AspectRatioSquare,
		),
	)
	yandexCmd.Flags().String(
		"safe-search",
		"",
		makeFlagUsageMessage(
			"Tell SafeSearch whether to filter sexually explicit content",
			params.SafeSearchOn,
			params.SafeSearchOff,
			params.SafeSearchModerate,
		),
	)

	return yandexCmd
}
