package util

import (
	"fmt"
	"os"
	"text/template"
)

// PrintResults writes search result to stdout.
func PrintResults(items []interface{}, pages int, compact bool) {
	compactTmplText := "{{.URL}}\n"
	tmplText := `Title: {{.Title}}
Webpage: {{.ReferenceURL}}
Resolution: {{.Width}}x{{.Height}}
URL: {{.URL}}
Thumbnail: {{.ThumbnailURL}}
`
	var tmpl *template.Template

	if compact {
		tmpl, _ = template.New("scraperOut").Parse(compactTmplText)
	} else {
		tmpl, _ = template.New("scraperOut").Parse(tmplText)
	}

	for _, item := range items {
		_ = tmpl.Execute(os.Stdout, item)

		if !compact {
			fmt.Println("---")
		}
	}

	if !compact {
		fmt.Printf("Scraped %d pages, got %d items\n", pages, len(items))
	}
}
