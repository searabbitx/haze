package mutable

import (
	"github.com/kamil-s-solecki/haze/http"
)

type Mutable struct {
	Name  string
	Apply func(http.Request, func(string) string) []http.Request
}

func AllMutatables() []Mutable {
	return []Mutable{Path, Parameter, ParameterName, BodyParameter, BodyParameterName, MultipartFormParameter, Header, Cookie, JsonParameter, JsonParameterRaw, CookieJsonParameter}
}
