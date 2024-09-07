package mutation

import (
	"github.com/kamil-s-solecki/haze/http"
	"github.com/kamil-s-solecki/haze/testutils"
	"testing"
)

func TestEmpty(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{})

	testutils.AssertEmpty(t, got)
}

func TestApplySingleQuotesMutationToPath(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{SingleQuotes})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Path, "/somepath'")
	testutils.AssertEquals(t, got[0].RequestUri, "/somepath'")
}

func TestApplyDoubleQuotesMutationToPath(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{DoubleQuotes})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Path, "/somepath\"")
	testutils.AssertEquals(t, got[0].RequestUri, "/somepath\"")
}
