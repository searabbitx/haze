package mutation

import (
	"github.com/kamil-s-solecki/haze/http"
	"strings"
)

type Mutable func(http.Request, func(string) string) []http.Request

func Path(rq http.Request, trans func(string) string) []http.Request {
	return []http.Request{rq.WithPath(trans(rq.Path))}
}

func Parameter(rq http.Request, trans func(string) string) []http.Request {
	if rq.Query == "" {
		return []http.Request{}
	}

	result := []http.Request{}
	for _, p := range strings.Split(rq.Query, "&") {
		q := strings.Replace(rq.Query, p, trans(p), 1)
		result = append(result, rq.WithQuery(q))
	}
	return result
}

func AllMutatables() []Mutable {
	return []Mutable{Path, Parameter}
}
