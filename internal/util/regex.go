package util

import (
	"fmt"
	"regexp"
)

// SearchRegex performs regex search on given string using a pattern returning the first matching group. In case of
// failure it returns an error (depending on fatal).
func SearchRegex(expr string, target string, name string, fatal bool) (string, error) {
	r, err := regexp.Compile(expr)

	if err != nil {
		if fatal {
			return "", fmt.Errorf("regex compilation error: %v", err)
		}

		return "", nil
	}

	subMatch := r.FindStringSubmatch(target)

	if subMatch == nil {
		return "", fmt.Errorf("cannot parse %s", name)
	}

	return subMatch[1], nil
}
