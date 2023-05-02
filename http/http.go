package http

import (
	"bytes"
	"strings"
)

type Request struct {
	Method          string
	Path            string
	ProtocolVersion string
	Headers         map[string]string
	Body            []byte
}

func Parse(bs []byte) Request {
	requestLine := bytes.Split(bs, []byte("\r\n"))[0]
	method, path, protocolVersion := parseRequestLine(requestLine)

	headers := parseHeaders(bs)

	body := extractBody(bs)
	return Request{Method: method, Path: path,
		ProtocolVersion: protocolVersion, Headers: headers, Body: body}
}

func parseRequestLine(requestLine []byte) (method, path, protocolVersion string) {
	spaceSplitted := bytes.Split(requestLine, []byte(" "))
	method = string(spaceSplitted[0])
	path = string(spaceSplitted[1])
	protocolVersion = string(spaceSplitted[2])
	return
}

func parseHeaders(rawReq []byte) (headers map[string]string) {
	headers = make(map[string]string)
	for _, rawHeader := range bytes.Split(rawReq, []byte("\r\n"))[1:] {
		if len(rawHeader) == 0 {
			break
		}
		name, val := parseHeader(rawHeader)
		headers[name] = val
	}
	return
}

func parseHeader(rawHeader []byte) (name, val string) {
	colonSplitted := bytes.Split(rawHeader, []byte(":"))
	name = string(colonSplitted[0])
	val = string(colonSplitted[1])
	val = strings.TrimSpace(val)
	return
}

func extractBody(rawReq []byte) []byte {
	twoRns := []byte("\r\n\r\n")
	bodyIndex := bytes.Index(rawReq, twoRns) + len(twoRns)
	return rawReq[bodyIndex:]
}
