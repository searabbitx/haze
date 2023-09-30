package mutable

import (
	"github.com/kamil-s-solecki/haze/http"
	"github.com/kamil-s-solecki/haze/utils"
	"strings"
)

var CookieJsonParameter = Mutable{"CookieJsonParameter", cookieJsonParameter}

func cookieJsonParameter(rq http.Request, trans func(string) string) []http.Request {
	identity := func(b []byte) []byte {
		return b
	}
	return cookieJsonParameterWithPostProcessing(rq, trans, identity)
}

func urlDecode(val string) string {
	return strings.Replace(val, "%22", "\"", -1)
}

func cookieJsonParameterWithPostProcessing(rq http.Request, trans func(string) string, post func([]byte) []byte) []http.Request {
	result := []http.Request{}

	for key, val := range rq.Cookies {
		if !rq.HasJsonCookie(key) {
			continue
		}
		data, _ := decodeJson([]byte(urlDecode(val)))
		for _, mutJson := range mutateJson(data, trans) {
			result = append(result, rq.WithCookie(key, utils.UrlEncodeSpecials(string(post(mutJson)))))
		}
	}
	return result
}
