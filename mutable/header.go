package mutable

import (
	"github.com/kamil-s-solecki/haze/http"
)

var Header = Mutable{"Header", header}

func header(rq http.Request, trans func(string) string) []http.Request {
	result := []http.Request{}
	for key, val := range rq.Headers {
		switch key {
		case "Content-Type", "Accept-Encoding", "Content-Encoding",
			"Connection", "Content-Length", "Host":
			continue
		}
		result = append(result, rq.WithHeader(key, trans(val)))
	}
	return result
}
