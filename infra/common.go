package infra

import "strings"

// ContainsString tells whether the target string is in given slice.
func ContainsString(target string, li []string) bool {
	for _, elem := range li {
		if target == elem {
			return true
		}
	}
	return false
}

func TrimSpace(src string) string {
	return strings.Trim(src, " ")
}
