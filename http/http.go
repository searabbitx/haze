package http

import (
	"bytes"
	"net/http"
	"net/http/httputil"
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

type Response struct {
	Code int
	Raw  []byte
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

func (r Request) asHttpReq(host string) *http.Request {
	url := host + r.RequestUri
	req, err := http.NewRequest(r.Method, url, nil)
	if err != nil {
		panic(err)
	}

	for key, val := range r.Headers {
		req.Header.Set(key, val)
	}
	return req
}

func (r Request) Send(host string) Response {
	req := r.asHttpReq(host)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	raw, err := httputil.DumpResponse(res, true)
	return Response{res.StatusCode, raw}
}

func (r Request) Raw(host string) []byte {
	bs, _ := httputil.DumpRequestOut(r.asHttpReq(host), true)
	return bs
}

func (r Request) WithPath(path string) Request {
	result := r.Clone()
	result.RequestUri = strings.Replace(r.RequestUri, r.Path, path, 1)
	result.Path = path
	return result
}

func (r Request) WithQuery(query string) Request {
	result := r.Clone()
	result.RequestUri = strings.Replace(r.RequestUri, r.Query, query, 1)
	result.Query = query
	return result
}

func (r Request) WithBody(body []byte) Request {
	result := r.Clone()
	result.Body = body
	return result
}

func (r Request) WithHeader(key, val string) Request {
	result := r.Clone()
	result.Headers[key] = val
	return result
}

func (r Request) Clone() Request {
	return Request{Method: r.Method, RequestUri: r.RequestUri, Path: r.Path, Query: r.Query,
		ProtocolVersion: r.ProtocolVersion, Headers: copyHeaders(r.Headers), Body: r.Body}
}

func copyHeaders(hs map[string]string) map[string]string {
	res := make(map[string]string)
	for k, v := range hs {
		res[k] = v
	}
	return res
}
