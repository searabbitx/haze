package mutable

import (
	"bytes"
	"github.com/kamil-s-solecki/haze/http"
	"regexp"
)

var MultipartFormParameter = Mutable{"MultipartFormParameter", multipartFormParameter}

func multipartFormParameter(rq http.Request, trans func(string) string) []http.Request {
	boundary := extractBoundary(rq)
	return []http.Request{rq.WithBody(mutateNextValue(rq.Body, boundary, trans))}
}

func mutateNextValue(body, boundary []byte, trans func(string) string) []byte {
	start, end := findValueRange(body, boundary)
	val := copySlice(body, start, end)
	val = []byte(trans(string(val)))

	result := copySlice(body, 0, start)
	result = append(result, val...)
	result = append(result, copySlice(body, end, len(body))...)
	return result
}

func findValueRange(body, boundary []byte) (int, int) {
	twoRns := []byte("\r\n\r\n")
	start := bytes.Index(body, twoRns) + len(twoRns)
	end := bytes.Index(body[start:], boundary) + start
	return start, end
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
