package http

import (
	"bytes"
	"reflect"
	"testing"
)

func assertEquals[T comparable](t *testing.T, got T, want T) {
	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func assertMapEquals[T comparable](t *testing.T, got map[T]T, want map[T]T) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func assertByteEquals(t *testing.T, got []byte, want []byte) {
	if !bytes.Equal(got, want) {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

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

		assertEquals(t, got, want)
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

		assertEquals(t, got, want)
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

		assertEquals(t, got, want)
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

		assertMapEquals(t, got, want)
	}
}

func TestBody(t *testing.T) {
	req := []byte("POST /auth HTTP/1.1\r\nContent-Type: application/x-www-form-urlencoded\r\nContent-Length: 7\r\n\r\nfoo=bar")

	got := Parse(req).Body
	want := []byte("foo=bar")

	assertByteEquals(t, got, want)
}
