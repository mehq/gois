package util

import (
	"fmt"
	"os"
	"text/template"
)


// PrintResults writes search result to stdout.
func PrintResults(items []interface{}, compact bool) {
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

	itemsLen := len(items)

	for i, item := range items {
		_ = tmpl.Execute(os.Stdout, item)

		if i != itemsLen-1 && !compact {
			fmt.Println("---")
		}
	}
}
