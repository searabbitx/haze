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

func BodyParameter(rq http.Request, trans func(string) string) []http.Request {
	if len(rq.Body) == 0 {
		return []http.Request{}
	}

	result := []http.Request{}
	body := string(rq.Body)
	for _, p := range strings.Split(body, "&") {
		q := strings.Replace(body, p, trans(p), 1)
		result = append(result, rq.WithBody([]byte(q)))
	}
	return result
}

func Header(rq http.Request, trans func(string) string) []http.Request {
	result := []http.Request{}
	for key, val := range rq.Headers {
		result = append(result, rq.WithHeader(key, trans(val)))
	}
	return result
}

func AllMutatables() []Mutable {
	return []Mutable{Path, Parameter, BodyParameter, Header}
}
