package mutable

import (
	"bytes"
	"github.com/kamil-s-solecki/haze/http"
	"regexp"
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
	return val
}

var Parameter = Mutable{"Parameter", parameter}

func parameter(rq http.Request, trans func(string) string) []http.Request {
	result := []http.Request{}
	if rq.Query == "" {
		return result
	}
	do := func(key, val string) (string, string) {
		return key, urlEncodeSpecials(trans(val))
	}
	for _, q := range applyToEachParam(rq.Query, do) {
		result = append(result, rq.WithQuery(q))
	}
	return result
}

var ParameterName = Mutable{"ParameterName", parameterName}

func parameterName(rq http.Request, trans func(string) string) []http.Request {
	result := []http.Request{}
	if rq.Query == "" {
		return result
	}
	do := func(key, val string) (string, string) {
		return trans(key), val
	}
	for _, q := range applyToEachParam(rq.Query, do) {
		result = append(result, rq.WithQuery(q))
	}
	return result
}

var BodyParameter = Mutable{"BodyParameter", bodyParameter}

func bodyParameter(rq http.Request, trans func(string) string) []http.Request {
	result := []http.Request{}
	if len(rq.Body) == 0 || !rq.HasFormUrlEncodedBody() {
		return result
	}
	do := func(key, val string) (string, string) {
		return key, urlEncodeSpecials(trans(val))
	}
	for _, q := range applyToEachParam(string(rq.Body), do) {
		result = append(result, rq.WithBody([]byte(q)))
	}
	return result
}

var BodyParameterName = Mutable{"BodyParameterName", bodyParameterName}

func bodyParameterName(rq http.Request, trans func(string) string) []http.Request {
	result := []http.Request{}
	if len(rq.Body) == 0 || !rq.HasFormUrlEncodedBody() {
		return result
	}
	do := func(key, val string) (string, string) {
		return trans(key), val
	}
	for _, q := range applyToEachParam(string(rq.Body), do) {
		result = append(result, rq.WithBody([]byte(q)))
	}
	return result
}

func applyToEachParam(params string, do func(key, val string) (string, string)) []string {
	result := []string{}
	for _, p := range strings.Split(params, "&") {
		key := strings.Split(p, "=")[0]
		val := strings.Split(p, "=")[1]
		mutKey, mutVal := do(key, val)
		q := strings.Replace(params, p, mutKey+"="+mutVal, 1)
		result = append(result, q)
	}
	return result
}

var MultipartFormParameter = Mutable{"MultipartFormParameter", multipartFormParameter}

func multipartFormParameter(rq http.Request, trans func(string) string) []http.Request {
	boundary := extractBoundary(rq)
	twoRns := []byte("\r\n\r\n")
	start := bytes.Index(rq.Body, twoRns) + len(twoRns)
	end := bytes.Index(rq.Body[start:], boundary) + start

	mut := copySlice(rq.Body, 0, start)
	val := copySlice(rq.Body, start, end)
	mut = append(mut, []byte(trans(string(val)))...)
	mut = append(mut, copySlice(rq.Body, end, len(rq.Body))...)

	return []http.Request{rq.WithBody(mut)}
}

func extractBoundary(rq http.Request) []byte {
	r, _ := regexp.Compile("boundary=([^;]*)(;|$)")
	return []byte("\r\n--" + r.FindStringSubmatch(rq.Headers["Content-Type"])[1])
}

func copySlice(b []byte, start, end int) []byte {
	var result = make([]byte, end-start)
	copy(result, b[start:end])
	return result
}

func AllMutatables() []Mutable {
	return []Mutable{Path, Parameter, ParameterName, BodyParameter, BodyParameterName, MultipartFormParameter, Header, Cookie, JsonParameter}
}
