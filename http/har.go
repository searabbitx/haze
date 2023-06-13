package http

import (
	"encoding/json"
	"net/url"
	"strings"
)

func ParseHar(data []byte, target string) []Request {
	har := unmarshalData(data)

	result := []Request{}
	forEachEntry(har, func(entry map[string]interface{}) {
		entryUrl := extractHarUrl(entry).String()
		if strings.HasPrefix(entryUrl, target) {
			result = append(result, entryToRequest(entry))
		}
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
		Method:     extractHarMethod(entry),
		RequestUri: extractHarRequestUri(entry),
		Path:       extractHarPath(entry),
		Query:      extractHarQuery(entry),
		Cookies:    extractHarCookies(entry),
		Headers:    extractHarHeaders(entry),
		Body:       extractHarBody(entry),
	}
}

func extractHarMethod(entry map[string]interface{}) string {
	return entry["method"].(string)
}

func extractHarRequestUri(entry map[string]interface{}) string {
	url := extractHarUrl(entry)
	return url.Path + "?" + url.RawQuery
}

func extractHarPath(entry map[string]interface{}) string {
	url := extractHarUrl(entry)
	return url.Path
}

func extractHarQuery(entry map[string]interface{}) string {
	url := extractHarUrl(entry)
	return url.RawQuery
}

func extractHarUrl(entry map[string]interface{}) *url.URL {
	url, _ := url.Parse(entry["url"].(string))
	return url
}

func extractHarCookies(entry map[string]interface{}) map[string]string {
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

func extractHarHeaders(entry map[string]interface{}) map[string]string {
	result := map[string]string{}

	headers := entry["headers"].([]interface{})
	for _, header := range headers {
		header := header.(map[string]interface{})
		name := header["name"].(string)
		val := header["value"].(string)
		switch name {
		case "Connection", "Host", "Cookie":
			continue
		default:
			result[name] = val
		}
	}

	return result
}

func extractHarBody(entry map[string]interface{}) []byte {
	postData, ok := entry["postData"].(map[string]interface{})
	if !ok {
		return []byte{}
	}
	return []byte(postData["text"].(string))
}
