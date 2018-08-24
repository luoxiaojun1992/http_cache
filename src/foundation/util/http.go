package util

import "strings"

func IfCache(cacheControl string) bool {
	return !strings.Contains(cacheControl, "no-cache") &&
		!strings.Contains(cacheControl, "no-store") &&
		!strings.Contains(cacheControl, "max-age=0")
}
