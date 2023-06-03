package mutable

import (
	"github.com/kamil-s-solecki/haze/http"
)

var Path = Mutable{"Path", path}

func path(rq http.Request, trans func(string) string) []http.Request {
	noLeadingSlash := rq.Path[1:]
	val := urlEncodeSpecials(trans(noLeadingSlash))
	return []http.Request{rq.WithPath("/" + val)}
}
