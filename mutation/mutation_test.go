package mutation

import (
	"github.com/kamil-s-solecki/haze/http"
	"github.com/kamil-s-solecki/haze/testutils"
	"testing"
)

func TestEmpty(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{}, []Mutable{})

	testutils.AssertEmpty(t, got)
}

func TestApplySingleQuotesMutationToPath(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{SingleQuotes}, []Mutable{Path})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Path, "/somepath'")
	testutils.AssertEquals(t, got[0].RequestUri, "/somepath'")
}

func TestApplyDoubleQuotesMutationToPath(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{DoubleQuotes}, []Mutable{Path})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Path, "/somepath\"")
	testutils.AssertEquals(t, got[0].RequestUri, "/somepath\"")
}

func TestApplySingleQuotesMutationToParameter(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath?foo=bar HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{SingleQuotes}, []Mutable{Parameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Query, "foo=bar'")
	testutils.AssertEquals(t, got[0].RequestUri, "/somepath?foo=bar'")
}

func TestApplyDoubleQuotesMutationToParameter(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath?foo=bar HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{DoubleQuotes}, []Mutable{Parameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Query, "foo=bar\"")
	testutils.AssertEquals(t, got[0].RequestUri, "/somepath?foo=bar\"")
}

func TestApplyDoubleQuotesMutationToBothParameters(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath?foo=bar&baz=quix HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{DoubleQuotes}, []Mutable{Parameter})

	testutils.AssertLen(t, got, 2)
	testutils.AssertEquals(t, got[0].Query, "foo=bar\"&baz=quix")
	testutils.AssertEquals(t, got[0].RequestUri, "/somepath?foo=bar\"&baz=quix")
	testutils.AssertEquals(t, got[1].Query, "foo=bar&baz=quix\"")
	testutils.AssertEquals(t, got[1].RequestUri, "/somepath?foo=bar&baz=quix\"")
}

func TestDoNothingForEmptyQuery(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{DoubleQuotes}, []Mutable{Parameter})

	testutils.AssertLen(t, got, 0)
}

func TestDoNothingForEmptyBody(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{DoubleQuotes}, []Mutable{BodyParameter})

	testutils.AssertLen(t, got, 0)
}

func TestApplyDoubleQuotesMutationToBodyParameter(t *testing.T) {
	rq := http.Parse([]byte("POST /auth HTTP/1.1\r\nContent-Type: application/x-www-form-urlencoded\r\nContent-Length: 7\r\n\r\nfoo=bar"))

	got := Mutate(rq, []Mutation{DoubleQuotes}, []Mutable{BodyParameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertByteEquals(t, got[0].Body, []byte("foo=bar\""))
}

func TestApplySingleQuotesMutationToBodyParameter(t *testing.T) {
	rq := http.Parse([]byte("POST /auth HTTP/1.1\r\nContent-Type: application/x-www-form-urlencoded\r\nContent-Length: 7\r\n\r\nfoo=bar"))

	got := Mutate(rq, []Mutation{SingleQuotes}, []Mutable{BodyParameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertByteEquals(t, got[0].Body, []byte("foo=bar'"))
}

func TestApplySstiFuzzMutationToParameter(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath?foo=bar HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{SstiFuzz}, []Mutable{Parameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Query, "foo=bar${{<%25[%25'%22}}%25%5c.")
	testutils.AssertEquals(t, got[0].RequestUri, "/somepath?foo=bar${{<%25[%25'%22}}%25%5c.")
}
