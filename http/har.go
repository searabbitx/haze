package http

import (
	"encoding/json"
)

func ParseHar(data []byte) []Request {
	har := unmarshalData(data)

	result := []Request{}
	forEachEntry(har, func(entry map[string]interface{}) {
		result = append(result, entryToRequest(entry))
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

func forEachEntry(har map[string]interface{}, do func(map[string]interface{})) {
	log := har["log"].(map[string]interface{})
	entries := log["entries"].([]interface{})
	for _, entry := range entries {
		do(entry.(map[string]interface{}))
	}
}

func entryToRequest(entry map[string]interface{}) Request {
	method := extractMethod(entry)
	return Request{Method: method}
}

func extractRequest(entry map[string]interface{}) map[string]interface{} {
	return entry["request"].(map[string]interface{})
}

func extractMethod(entry map[string]interface{}) string {
	return extractRequest(entry)["method"].(string)
}
