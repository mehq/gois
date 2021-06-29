package out

import (
	"fmt"
	"os"
	"text/template"
)

func getOutTemplate(compact bool) (*template.Template, error) {
	if compact {
		tmpl, err := template.New("scraperOut").Parse("{{.URL}}\n")

		if err != nil {
			return nil, fmt.Errorf("failed to parse text as output template (compact mode)")
		}

		return tmpl, nil
	}

	tmpl, err := template.New("scraperOut").Parse(`Title: {{.Title}}
Webpage: {{.ReferenceURL}}
Resolution: {{.Width}}x{{.Height}}
URL: {{.URL}}
Thumbnail: {{.ThumbnailURL}}
`)

	if err != nil {
		return nil, fmt.Errorf("failed to parse text as output template")
	}

	return tmpl, nil
}

// PrintImageInfo writes search result to stdout.
func PrintImageInfo(items []interface{}, compact bool) error {
	tmpl, err := getOutTemplate(compact)

	if err != nil {
		return fmt.Errorf("failed to get output template")
	}

	itemsLen := len(items)

	for i, item := range items {
		err = tmpl.Execute(os.Stdout, item)

		if err != nil {
			return fmt.Errorf("failed to execute template")
		}

		if i != itemsLen-1 && !compact {
			fmt.Println("---")
		}
	}

	return nil
}
