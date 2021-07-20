package util

import (
	"fmt"
	"regexp"
)

func SearchRegex(expr string, target string, name string, fatal bool) (string, error) {
	r, err := regexp.Compile(expr)

	if err != nil {
		if fatal {
			return "", err
		}

		return "", nil
	}

	subMatch := r.FindStringSubmatch(target)

	if subMatch == nil {
		return "", fmt.Errorf("cannot parse %s", name)
	}

	return subMatch[1], nil
}

func SearchRegexMultiple(expr string, target string, name string, fatal bool) ([]string, error) {
	r, err := regexp.Compile(expr)

	if err != nil {
		if fatal {
			return nil, err
		}

		return nil, nil
	}

	subMatches := r.FindAllStringSubmatch(target, -1)

	if subMatches == nil {
		if !fatal {
			return nil, nil
		}

		return nil, fmt.Errorf("cannot match regex for %s", name)
	}

	result := make([]string, len(subMatches))

	for i := 0; i < len(result); i += 1 {
		result[i] = subMatches[i][1]
	}

	return result, nil
}
