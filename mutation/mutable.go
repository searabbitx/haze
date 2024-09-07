package mutation

import (
	"github.com/kamil-s-solecki/haze/http"
	"strings"
)

type Mutable func(http.Request, func(string) string) []http.Request

func urlEncodeSpecials(val string) string {
	val = strings.Replace(val, "%", "%25", -1)
	val = strings.Replace(val, "\\", "%5c", -1)
	val = strings.Replace(val, "\"", "%22", -1)
	return val
}

func Path(rq http.Request, trans func(string) string) []http.Request {
	val := urlEncodeSpecials(trans(rq.Path))
	return []http.Request{rq.WithPath(val)}
}

func Parameter(rq http.Request, trans func(string) string) []http.Request {
	if rq.Query == "" {
		return []http.Request{}
	}

	result := []http.Request{}
	for _, p := range strings.Split(rq.Query, "&") {
		key := strings.Split(p, "=")[0]
		val := strings.Split(p, "=")[1]
		keyNVal := key + "=" + urlEncodeSpecials(trans(val))

		q := strings.Replace(rq.Query, p, keyNVal, 1)
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
		key := strings.Split(p, "=")[0]
		val := strings.Split(p, "=")[1]
		keyNVal := key + "=" + urlEncodeSpecials(trans(val))

		q := strings.Replace(body, p, keyNVal, 1)
		result = append(result, rq.WithBody([]byte(q)))
	}
	return result
}

func Header(rq http.Request, trans func(string) string) []http.Request {
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

func Cookie(rq http.Request, trans func(string) string) []http.Request {
	result := []http.Request{}
	for key, val := range rq.Cookies {
		enc := urlEncodeSpecials(trans(val))
		result = append(result, rq.WithCookie(key, enc))
	}
	return result
}

func AllMutatables() []Mutable {
	return []Mutable{Path, Parameter, BodyParameter, Header, Cookie}
}
