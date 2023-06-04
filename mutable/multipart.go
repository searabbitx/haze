package mutable

import (
	"bytes"
	"github.com/kamil-s-solecki/haze/http"
	"regexp"
)

var MultipartFormParameter = Mutable{"MultipartFormParameter", multipartFormParameter}

func multipartFormParameter(rq http.Request, trans func(string) string) []http.Request {
	boundary := extractBoundary(rq)
	next := func(from int) ([]byte, int) {
		return mutateNextValue(rq.Body, boundary, from, trans)
	}

	result := []http.Request{}
	for mut, i := next(0); i != -1; mut, i = next(i) {
		result = append(result, rq.WithBody(mut))
	}
	return result
}

func mutateNextValue(body, boundary []byte, from int, trans func(string) string) ([]byte, int) {
	start, end := findValueRange(body, boundary, from)
	if start == -1 || end == -1 {
		return []byte{}, -1
	}

	val := copySlice(body, start, end)
	val = []byte(trans(string(val)))

	result := copySlice(body, 0, start)
	result = append(result, val...)
	result = append(result, copySlice(body, end, len(body))...)

	stoppedAt := end + len(boundary)
	return result, stoppedAt
}

func findValueRange(body, boundary []byte, from int) (int, int) {
	body = body[from:]
	twoRns := []byte("\r\n\r\n")

	twoRnsIdx := bytes.Index(body, twoRns) + len(twoRns)
	boundaryIdx := bytes.Index(body[twoRnsIdx:], boundary)
	if twoRnsIdx == -1 || boundaryIdx == -1 {
		return -1, -1
	}

	return from + twoRnsIdx, from + twoRnsIdx + boundaryIdx
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
