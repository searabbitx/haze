package http

import (
	"encoding/json"
	"net/url"
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
	for _, e := range entries {
		erq := e.(map[string]interface{})["request"]
		do(erq.(map[string]interface{}))
	}
}

func entryToRequest(entry map[string]interface{}) Request {
	return Request{
		Method:     extractMethod(entry),
		RequestUri: extractRequestUri(entry),
		Path:       extractPath(entry),
		Query:      extractQuery(entry),
		Cookies:    extractCookies(entry),
	}
}

func extractMethod(entry map[string]interface{}) string {
	return entry["method"].(string)
}

func extractRequestUri(entry map[string]interface{}) string {
	url := extractUrl(entry)
	return url.Path + "?" + url.RawQuery
}

func extractPath(entry map[string]interface{}) string {
	url := extractUrl(entry)
	return url.Path
}

func extractQuery(entry map[string]interface{}) string {
	url := extractUrl(entry)
	return url.RawQuery
}

func extractUrl(entry map[string]interface{}) *url.URL {
	url, _ := url.Parse(entry["url"].(string))
	return url
}

func extractCookies(entry map[string]interface{}) map[string]string {
	result := map[string]string{}

	cookies := entry["cookies"].([]interface{})
	for _, cookie := range cookies {
		cookie := cookie.(map[string]interface{})
		name := cookie["name"].(string)
		val := cookie["value"].(string)
		result[name] = val
	}

	return result
}
