package util

import (
	"fmt"
	"regexp"
)

// SearchRegex performs regex search on given string using a pattern returning the first matching group. In case of
// failure it returns an error.
func SearchRegex(expr string, target string, fieldName string) (string, error) {
	r, err := regexp.Compile(expr)

	if err != nil {
		return "", fmt.Errorf("regex compilation error: %v", err)
	}

	subMatch := r.FindStringSubmatch(target)

	if subMatch == nil || len(subMatch) < 2 {
		return "", fmt.Errorf("cannot parse %s", fieldName)
	}

	return subMatch[1], nil
}
