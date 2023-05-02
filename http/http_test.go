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

func TestPath(t *testing.T) {
	cases := []struct {
		req  []byte
		path string
	}{
		{[]byte("GET /somepath HTTP/1.1\r\nHost:www.example.com\r\n\r\n"), "/somepath"},
		{[]byte("GET /otherpath HTTP/1.1\r\nHost:www.example.com\r\n\r\n"), "/otherpath"},
	}

	for _, c := range cases {
		got := Parse(c.req).Path
		want := c.path

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
