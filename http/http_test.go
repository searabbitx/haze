package http

import (
	"github.com/kamil-s-solecki/haze/testutils"
	"testing"
)

func TestMethod(t *testing.T) {
	cases := []struct {
		req    []byte
		method string
	}{
		{[]byte("GET /somepath HTTP/1.1\r\nHost:www.example.com\r\n\r\n"), "GET"},
		{[]byte("POST /somepath HTTP/1.1\r\nHost:www.example.com\r\n\r\n"), "POST"},
	}

	for _, c := range cases {
		got := Parse(c.req).Method
		want := c.method

		testutils.AssertEquals(t, got, want)
	}
}

func TestRequestUri(t *testing.T) {
	cases := []struct {
		req        []byte
		requestUri string
	}{
		{[]byte("GET /somepath HTTP/1.1\r\nHost:www.example.com\r\n\r\n"), "/somepath"},
		{[]byte("GET /otherpath?foo=bar HTTP/1.1\r\nHost:www.example.com\r\n\r\n"), "/otherpath?foo=bar"},
	}

	for _, c := range cases {
		got := Parse(c.req).RequestUri
		want := c.requestUri

		testutils.AssertEquals(t, got, want)
	}
}

func TestProtocolVersion(t *testing.T) {
	cases := []struct {
		req []byte
		pv  string
	}{
		{[]byte("GET /somepath HTTP/1.1\r\nHost:www.example.com\r\n\r\n"), "HTTP/1.1"},
		{[]byte("GET /somepath HTTP/2.0\r\nHost:www.example.com\r\n\r\n"), "HTTP/2.0"},
	}

	for _, c := range cases {
		got := Parse(c.req).ProtocolVersion
		want := c.pv

		testutils.AssertEquals(t, got, want)
	}
}

func TestHeaders(t *testing.T) {
	cases := []struct {
		req []byte
		hs  map[string]string
	}{
		{[]byte("GET /somepath HTTP/1.1\r\nHost:www.example.com\r\n\r\n"), map[string]string{"Host": "www.example.com"}},
		{[]byte("GET /somepath HTTP/1.1\r\nHost:example.com\r\n\r\n"), map[string]string{"Host": "example.com"}},
		{[]byte("GET /somepath HTTP/1.1\r\nHost: example.com\r\n\r\n"), map[string]string{"Host": "example.com"}},
		{[]byte("GET /somepath HTTP/1.1\r\nHost: example.com\r\nConnection: close\r\n\r\n"),
			map[string]string{"Host": "example.com", "Connection": "close"}},
	}

	for _, c := range cases {
		got := Parse(c.req).Headers
		want := c.hs

		testutils.AssertMapEquals(t, got, want)
	}
}

func TestBody(t *testing.T) {
	req := []byte("POST /auth HTTP/1.1\r\nContent-Type: application/x-www-form-urlencoded\r\nContent-Length: 7\r\n\r\nfoo=bar")

	got := Parse(req).Body
	want := []byte("foo=bar")

	testutils.AssertByteEquals(t, got, want)
}

func TestPath(t *testing.T) {
	cases := []struct {
		req  []byte
		path string
	}{
		{[]byte("GET /somepath HTTP/1.1\r\nHost:www.example.com\r\n\r\n"), "/somepath"},
		{[]byte("GET /somepath?foo=bar HTTP/1.1\r\nHost:www.example.com\r\n\r\n"), "/somepath"},
	}

	for _, c := range cases {
		got := Parse(c.req).Path
		want := c.path

		testutils.AssertEquals(t, got, want)
	}
}

func TestQuery(t *testing.T) {
	cases := []struct {
		req   []byte
		query string
	}{
		{[]byte("GET /somepath HTTP/1.1\r\nHost:www.example.com\r\n\r\n"), ""},
		{[]byte("GET /somepath?foo=bar HTTP/1.1\r\nHost:www.example.com\r\n\r\n"), "foo=bar"},
	}

	for _, c := range cases {
		got := Parse(c.req).Query
		want := c.query

		testutils.AssertEquals(t, got, want)
	}
}

func TestClone(t *testing.T) {
	want := Request{
		Method:          "Foo1",
		RequestUri:      "Foo2",
		Path:            "Foo2",
		Query:           "Foo3",
		ProtocolVersion: "Foo4",
		Headers:         map[string]string{"Bar": "Baz"},
		Body:            []byte("Foo5"),
	}

	got := want.Clone()

	testutils.AssertEquals(t, want.Method, got.Method)
	testutils.AssertEquals(t, want.RequestUri, got.RequestUri)
	testutils.AssertEquals(t, want.Path, got.Path)
	testutils.AssertEquals(t, want.Query, got.Query)
	testutils.AssertEquals(t, want.ProtocolVersion, got.ProtocolVersion)
	testutils.AssertMapEquals(t, want.Headers, got.Headers)
	testutils.AssertByteEquals(t, want.Body, got.Body)
}

func TestCloneHeaders(t *testing.T) {
	orig := Request{
		Method:          "Foo1",
		RequestUri:      "Foo2",
		Path:            "Foo2",
		Query:           "Foo3",
		ProtocolVersion: "Foo4",
		Headers:         map[string]string{"Bar": "Baz"},
		Body:            []byte("Foo5"),
	}
	clone := orig.Clone()
	clone.Headers["Bar"] = "Edit"

	testutils.AssertEquals(t, orig.Headers["Bar"], "Baz")
}

func TestCookies(t *testing.T) {
	cases := []struct {
		req []byte
		cks map[string]string
	}{
		{[]byte("GET /somepath HTTP/1.1\r\nHost:www.example.com\r\nCookie:foo=bar\r\n\r\n"), map[string]string{"foo": "bar"}},
		{[]byte("GET /somepath HTTP/1.1\r\nHost:www.example.com\r\nCookie:foo=bar \r\n\r\n"), map[string]string{"foo": "bar"}},
		{[]byte("GET /somepath HTTP/1.1\r\nHost:www.example.com\r\nCookie:foo=bar; baz=quix\r\n\r\n"),
			map[string]string{"foo": "bar", "baz": "quix"}},
	}

	for _, c := range cases {
		r := Parse(c.req)
		got := r.Cookies
		want := c.cks

		testutils.AssertMapEquals(t, got, want)
		testutils.AssertMapHasNoKey(t, r.Headers, "Cookie")
	}
}

func TestResponseStringer(t *testing.T) {
	cases := []struct {
		res Response
		str string
	}{
		{Response{Code: 200, Length: 1234}, "[Code: 200, Len: 1234]"},
		{Response{Code: 400, Length: 4321}, "[Code: 400, Len: 4321]"},
	}

	for _, c := range cases {
		res := c.res

		got := res.String()

		testutils.AssertEquals(t, got, c.str)
	}
}
