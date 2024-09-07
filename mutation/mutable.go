package mutation

import (
	"encoding/json"
	"fmt"
	"github.com/kamil-s-solecki/haze/http"
	"strings"
)

type Mutable struct {
	name  string
	apply func(http.Request, func(string) string) []http.Request
}

func urlEncodeSpecials(val string) string {
	val = strings.Replace(val, "%", "%25", -1)
	val = strings.Replace(val, "\\", "%5c", -1)
	val = strings.Replace(val, "\"", "%22", -1)
	return val
}

var Path = Mutable{"Path", path}

func path(rq http.Request, trans func(string) string) []http.Request {
	noLeadingSlash := rq.Path[1:]
	val := urlEncodeSpecials(trans(noLeadingSlash))
	return []http.Request{rq.WithPath("/" + val)}
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

var Cookie = Mutable{"Cookie", cookie}

func cookie(rq http.Request, trans func(string) string) []http.Request {
	result := []http.Request{}
	for key, val := range rq.Cookies {
		enc := urlEncodeSpecials(trans(val))
		result = append(result, rq.WithCookie(key, enc))
	}
	return result
}

var JsonParameter = Mutable{"JsonParameter", jsonParameter}

func jsonParameter(rq http.Request, trans func(string) string) []http.Request {
	if !rq.HasJsonBody() {
		return []http.Request{}
	}

	data := decodeJson(rq.Body)
	result := []http.Request{}

	for _, mutJson := range mutateJson(data, trans) {
		result = append(result, rq.WithBody(mutJson))
	}
	return result
}

type JsonMutation struct {
	Apply  func()
	Revert func()
}

func mutateJson(data interface{}, trans func(string) string) [][]byte {
	result := [][]byte{}

	var muts []JsonMutation
	switch data.(type) {
	case []any:
		muts = mutateJsonArray(data.([]interface{}), trans)
	case map[string]interface{}:
		muts = mutateJsonRecursive(data.(map[string]interface{}), trans)
	default:
		muts = mutateFlatJson(&data, trans)
	}

	for _, jsonMut := range muts {
		jsonMut.Apply()
		js, _ := json.Marshal(data)
		result = append(result, js)
		jsonMut.Revert()
	}
	return result
}

func mutateFlatJson(data *interface{}, trans func(string) string) []JsonMutation {
	orig := *data
	return []JsonMutation{JsonMutation{
		Apply: func() {
			*data = trans(fmt.Sprintf("%v", orig))
		},
		Revert: func() {
			*data = orig
		},
	}}
}

func mutateJsonRecursive(data map[string]interface{}, trans func(string) string) []JsonMutation {
	agg := []JsonMutation{}
	for key, val := range data {
		switch val.(type) {
		case map[string]interface{}:
			subs := mutateJsonRecursive(val.(map[string]interface{}), trans)
			agg = append(agg, subs...)
		case []any:
			arr := val.([]interface{})
			agg = append(agg, mutateJsonArray(arr, trans)...)
		default:
			agg = append(agg, mutateJsonLeaf(data, key, trans))
		}
	}
	return agg
}

func mutateJsonArray(arr []interface{}, trans func(string) string) []JsonMutation {
	res := []JsonMutation{}
	for i, v := range arr {
		i := i
		v := v
		switch v.(type) {
		case map[string]interface{}:
			muts := mutateJsonRecursive(v.(map[string]interface{}), trans)
			res = append(res, muts...)
		default:
			mut := JsonMutation{
				Apply: func() {
					arr[i] = trans(fmt.Sprintf("%v", v))
				},
				Revert: func() {
					arr[i] = v
				},
			}
			res = append(res, mut)
		}
	}
	return res
}

func mutateJsonLeaf(data map[string]interface{}, key string, trans func(string) string) JsonMutation {
	val := data[key]
	return JsonMutation{
		Apply: func() {
			data[key] = trans(fmt.Sprintf("%v", val))
		},
		Revert: func() {
			data[key] = val
		},
	}
}

func decodeJson(bs []byte) interface{} {
	var data interface{}
	json.Unmarshal(bs, &data)
	return data
}

func AllMutatables() []Mutable {
	return []Mutable{Path, Parameter, ParameterName, BodyParameter, BodyParameterName, Header, Cookie, JsonParameter}
}
