package mutable

import (
	"github.com/kamil-s-solecki/haze/http"
	"strings"
)

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
