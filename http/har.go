package http

import (
	"encoding/json"
)

func ParseHar(data []byte) []Request {
	har := unmarshalData(data)

	result := []Request{}
	forEachEntry(har, func() {
		result = append(result, Request{})
	})
	return result
}

func unmarshalData(data []byte) map[string]interface{} {
	var har map[string]interface{}
	err := json.Unmarshal(data, &har)
	if err != nil {
		panic(err)
	}
	return har
}

func forEachEntry(har map[string]interface{}, do func()) {
	log := har["log"].(map[string]interface{})
	entries := log["entries"].([]interface{})
	for _ = range entries {
		do()
	}
}
