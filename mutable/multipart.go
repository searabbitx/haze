package mutable

import (
	"bytes"
	"github.com/kamil-s-solecki/haze/http"
	"regexp"
)

var MultipartFormParameter = Mutable{"MultipartFormParameter", multipartFormParameter}

func multipartFormParameter(rq http.Request, trans func(string) string) []http.Request {
	boundary := extractBoundary(rq)
	twoRns := []byte("\r\n\r\n")
	start := bytes.Index(rq.Body, twoRns) + len(twoRns)
	end := bytes.Index(rq.Body[start:], boundary) + start

	mut := copySlice(rq.Body, 0, start)
	val := copySlice(rq.Body, start, end)
	mut = append(mut, []byte(trans(string(val)))...)
	mut = append(mut, copySlice(rq.Body, end, len(rq.Body))...)

	return []http.Request{rq.WithBody(mut)}
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
