package mutable

import (
	"bytes"
	"github.com/kamil-s-solecki/haze/http"
	"regexp"
)

var MultipartFormParameter = Mutable{"MultipartFormParameter", multipartFormParameter}

func multipartFormParameter(rq http.Request, trans func(string) string) []http.Request {
	boundary := extractBoundary(rq)
	result := []http.Request{}
	body := rq.Body

	var mut []byte
	for i := 0; i != -1; {
		mut, i = mutateNextValue(body, boundary, i, trans)
		if i == -1 {
			break
		}
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
	start := bytes.Index(body, twoRns) + len(twoRns)
	end := bytes.Index(body[start:], boundary)
	if end == -1 || start == -1 {
		return -1, -1
	}
	return start + from, end + start + from
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
