package conv

import "strconv"

// Atoi is equivalent to strconv.Atoi except it returns 0 on error.
func Atoi(rawValue string) int {
	value, err := strconv.Atoi(rawValue)

	if err != nil {
		return 0
	}

	return value
}

// Itoa is equivalent to strconv.Itoa.
func Itoa(rawValue int) string {
	return strconv.Itoa(rawValue)
}
