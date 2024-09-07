package mutable

import (
	"encoding/json"
	"fmt"
	"github.com/kamil-s-solecki/haze/http"
)

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
