package utils

import (
	"strings"
)

func UrlEncodeSpecials(val string) string {
	val = strings.Replace(val, "%", "%25", -1)
	val = strings.Replace(val, "\\", "%5c", -1)
	val = strings.Replace(val, "\"", "%22", -1)
	val = strings.Replace(val, "\x00", "%00", -1)
	val = strings.Replace(val, " ", "%20", -1)
	val = strings.Replace(val, "\t", "%09", -1)
	val = strings.Replace(val, "\f", "%0c", -1)
	val = strings.Replace(val, "\r", "%0d", -1)
	val = strings.Replace(val, "\n", "%0a", -1)
	val = strings.Replace(val, ";", "%3b", -1)
	return val
}
