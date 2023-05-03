package http

import (
	"bytes"
	"net/http"
	"strings"
)

type Request struct {
	Method          string
	RequestUri      string
	Path            string
	Query           string
	ProtocolVersion string
	Headers         map[string]string
	Body            []byte
}

func Parse(bs []byte) Request {
	requestLine := bytes.Split(bs, []byte("\r\n"))[0]
	method, requestUri, protocolVersion := parseRequestLine(requestLine)
	path, query := parseRequestUri(requestUri)

	headers := parseHeaders(bs)

	body := extractBody(bs)
	return Request{Method: method, RequestUri: requestUri, Path: path, Query: query,
		ProtocolVersion: protocolVersion, Headers: headers, Body: body}
}

func parseRequestLine(requestLine []byte) (method, requestUri, protocolVersion string) {
	spaceSplitted := bytes.Split(requestLine, []byte(" "))
	method = string(spaceSplitted[0])
	requestUri = string(spaceSplitted[1])
	protocolVersion = string(spaceSplitted[2])
	return
}

func parseRequestUri(requestUri string) (path, query string) {
	if i := strings.Index(requestUri, "?"); i > 0 {
		path = requestUri[:i]
		query = requestUri[i+1:]
	} else {
		path = requestUri
	}
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

func (r Request) Send(host string) {
	url := host + r.RequestUri
	req, err := http.NewRequest(r.Method, url, nil)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}
}
