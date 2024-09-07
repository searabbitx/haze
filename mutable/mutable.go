package mutable

import (
	"github.com/kamil-s-solecki/haze/http"
	"strings"
)

type Mutable struct {
	Name  string
	Apply func(http.Request, func(string) string) []http.Request
}

func urlEncodeSpecials(val string) string {
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

func AllMutatables() []Mutable {
	return []Mutable{Path, Parameter, ParameterName, BodyParameter, BodyParameterName, MultipartFormParameter, Header, Cookie, JsonParameter, JsonParameterRaw}
}
