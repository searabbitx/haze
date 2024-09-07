package mutation

import "github.com/kamil-s-solecki/haze/http"

type Mutable func(http.Request, func(string) string) http.Request

func Path(rq http.Request, trans func(string) string) http.Request {
	return rq.WithPath(trans(rq.Path))
}

func Parameter(rq http.Request, trans func(string) string) http.Request {
	return rq.WithQuery(trans(rq.Query))
}

func AllMutatables() []Mutable {
	return []Mutable{Path, Parameter}
}
