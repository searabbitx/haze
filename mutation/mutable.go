package mutation

import (
	"encoding/json"
	"fmt"
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

func JsonParameter(rq http.Request, trans func(string) string) []http.Request {
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

func mutateJson(data map[string]interface{}, trans func(string) string) [][]byte {
	result := [][]byte{}
	for _, jsonMut := range mutateJsonRecursive(data, trans, []JsonMutation{}) {
		jsonMut.Apply()
		js, _ := json.Marshal(data)
		result = append(result, js)
		jsonMut.Revert()
	}
	return result
}

func mutateJsonRecursive(data map[string]interface{}, trans func(string) string, agg []JsonMutation) []JsonMutation {
	for key, val := range data {
		switch val.(type) {
		case map[string]interface{}:
			subs := mutateJsonRecursive(val.(map[string]interface{}), trans, agg)
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
			muts := mutateJsonRecursive(v.(map[string]interface{}), trans, []JsonMutation{})
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

func decodeJson(bs []byte) map[string]interface{} {
	var data map[string]interface{}
	json.Unmarshal(bs, &data)
	return data
}

func AllMutatables() []Mutable {
	return []Mutable{Path, Parameter, BodyParameter, Header, Cookie, JsonParameter}
}
