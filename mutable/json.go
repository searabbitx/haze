package mutable

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kamil-s-solecki/haze/http"
	"strings"
)

var JsonParameter = Mutable{"JsonParameter", jsonParameter}

func jsonParameter(rq http.Request, trans func(string) string) []http.Request {
	identity := func(b []byte) []byte {
		return b
	}
	return jsonParameterWithPostProcessing(rq, trans, identity)
}

var JsonParameterRaw = Mutable{"JsonParameterRaw", jsonParameterRaw}

func jsonParameterRaw(rq http.Request, trans func(string) string) []http.Request {
	ntrans := func(val string) string {
		return "QREMOVE" + strings.Replace(trans(val), `"`, "__QUOT__", -1) + "QREMOVE"
	}
	post := func(js []byte) []byte {
		js = bytes.Replace(js, []byte("\"QREMOVE"), []byte(""), -1)
		js = bytes.Replace(js, []byte("QREMOVE\""), []byte(""), -1)
		js = bytes.Replace(js, []byte("__QUOT__"), []byte(`"`), -1)
		return js
	}
	return jsonParameterWithPostProcessing(rq, ntrans, post)
}

func jsonParameterWithPostProcessing(rq http.Request, trans func(string) string, post func([]byte) []byte) []http.Request {
	if !rq.HasJsonBody() {
		return []http.Request{}
	}

	data := decodeJson(rq.Body)
	result := []http.Request{}

	for _, mutJson := range mutateJson(data, trans) {
		result = append(result, rq.WithBody(post(mutJson)))
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
	return []JsonMutation{{
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
