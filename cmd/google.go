package cmd

import (
	"github.com/mzbaulhaque/gois/internal/util"
	"github.com/spf13/cobra"

	"github.com/mzbaulhaque/gois/pkg/scraper"
)

func newGoogleCmd() *cobra.Command {
	googleCmd := &cobra.Command{
		Use:   "google [flags] <query>",
		Short: "Search images using google",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			aspectRatio, _ := c.Flags().GetString("aspect-ratio")
			compact, _ := c.Flags().GetBool("compact")
			fileType, _ := c.Flags().GetString("file-type")
			imageColor, _ := c.Flags().GetString("image-color")
			imageSize, _ := c.Flags().GetString("image-size")
			imageType, _ := c.Flags().GetString("image-type")
			license, _ := c.Flags().GetString("license")
			region, _ := c.Flags().GetString("region")
			safeSearch, _ := c.Flags().GetString("region")

			config := &scraper.GoogleConfig{
				AspectRatio: aspectRatio,
				Compact:  compact,
				FileType: fileType,
				ImageColor: imageColor,
				ImageSize: imageSize,
				ImageType: imageType,
				License: license,
				Query:    args[0],
				Region: region,
				SafeSearch: safeSearch,
			}
			gs := &scraper.GoogleScraper{Config: config}
			items, err := gs.Scrape()

			if err != nil {
				return err
			}

			util.PrintResults(items, config.Compact)

			return nil
		},
	}

	googleCmd.Flags().StringP("aspect-ratio", "A", "", "Specify the shape of images [any (default), tall, square, wide, panoramic]")
	googleCmd.PersistentFlags().BoolP("compact", "c", false, "Print original image link per line with no other information.")
	googleCmd.Flags().StringP("file-type", "F", "", "Find images in the format that you prefer [any (default), jpg, gif, png, bmp, svg, webp, ico, raw]")
	googleCmd.Flags().StringP("image-color", "C", "", "Find images in preferred color [any (default), full-color, black-white, transparent, red, orange, yellow, green, teal, blue, purple, pink, white, gray, black, brown]")
	googleCmd.Flags().StringP("image-size", "S", "", "Find images in specific size [any (default), large, medium, icon, qsvga, vga, svga, xga, 2mp, 4mp, 6mp, 8mp, 10mp, 12mp, 15mp, 20mp, 40mp, 70mp]")
	googleCmd.Flags().StringP("image-type", "T", "", "Limit the kind of images that you find [any (default), face, photo, clip-art, line-drawing, animated]")
	googleCmd.Flags().StringP("license", "L", "", "License preference [all (default), creative-commons, commercial]")
	googleCmd.Flags().StringP("region", "R", "", "Find images published in a particular region [ISO 3166-1 alpha-2 code of your preferred region e.g --region=US]")
	googleCmd.Flags().StringP("safe-search", "s", "on", "Tell SafeSearch whether to filter sexually explicit content [on (default), off]")

	return googleCmd
}
