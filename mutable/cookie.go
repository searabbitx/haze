package mutable

import (
	"github.com/kamil-s-solecki/haze/http"
)

var Cookie = Mutable{"Cookie", cookie}

func cookie(rq http.Request, trans func(string) string) []http.Request {
	result := []http.Request{}
	for key, val := range rq.Cookies {
		enc := urlEncodeSpecials(trans(val))
		result = append(result, rq.WithCookie(key, enc))
	}
	return result
}
