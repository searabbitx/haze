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
	testutils.AssertEquals(t, got[0].Path, "/somepath%22")
	testutils.AssertEquals(t, got[0].RequestUri, "/somepath%22")
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
	testutils.AssertEquals(t, got[0].Query, "foo=bar%22")
	testutils.AssertEquals(t, got[0].RequestUri, "/somepath?foo=bar%22")
}

func TestApplyDoubleQuotesMutationToBothParameters(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath?foo=bar&baz=quix HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{DoubleQuotes}, []Mutable{Parameter})

	testutils.AssertLen(t, got, 2)
	testutils.AssertEquals(t, got[0].Query, "foo=bar%22&baz=quix")
	testutils.AssertEquals(t, got[0].RequestUri, "/somepath?foo=bar%22&baz=quix")
	testutils.AssertEquals(t, got[1].Query, "foo=bar&baz=quix%22")
	testutils.AssertEquals(t, got[1].RequestUri, "/somepath?foo=bar&baz=quix%22")
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
	testutils.AssertByteEquals(t, got[0].Body, []byte("foo=bar%22"))
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

func TestApplyDoubleQuotesMutationToHeader(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath HTTP/1.1\r\nFoo: bar\r\n\r\n"))

	got := Mutate(rq, []Mutation{DoubleQuotes}, []Mutable{Header})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Headers["Foo"], "bar\"")
}

func TestApplySstiFuzzMutationToHeader(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath?foo=bar HTTP/1.1\r\nFoo:bar\r\n\r\n"))

	got := Mutate(rq, []Mutation{SstiFuzz}, []Mutable{Header})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Headers["Foo"], "bar${{<%[%'\"}}%\\.")
}

func TestApplyDoubleQuotesMutationToCookie(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath HTTP/1.1\r\nCookie: foo=bar\r\n\r\n"))

	got := Mutate(rq, []Mutation{DoubleQuotes}, []Mutable{Cookie})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Cookies["foo"], "bar%22")
}

func TestSkipCertainHeaders(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath HTTP/1.1\r\nHost: bar\r\nConnection: bar\r\nContent-Type: bar\r\nContent-Length: bar\r\nAccept-Encoding: bar\r\n\r\n"))

	got := Mutate(rq, []Mutation{DoubleQuotes}, []Mutable{Header})

	testutils.AssertLen(t, got, 0)
}

func TestApplyNegativeMutationToParameter(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath?foo=123 HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{Negative}, []Mutable{Parameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Query, "foo=-123")
	testutils.AssertEquals(t, got[0].RequestUri, "/somepath?foo=-123")
}

func TestApplyMinusOneMutationToParameter(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath?foo=123 HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{MinusOne}, []Mutable{Parameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Query, "foo=123-1")
	testutils.AssertEquals(t, got[0].RequestUri, "/somepath?foo=123-1")
}

func TestApplyTimesSevenMutationToParameter(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath?foo=123 HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{TimesSeven}, []Mutable{Parameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Query, "foo=123*7")
	testutils.AssertEquals(t, got[0].RequestUri, "/somepath?foo=123*7")
}

func TestDoNothingWithNonJsonBody(t *testing.T) {
	rq := http.Parse([]byte("POST /auth HTTP/1.1\r\nContent-Type: application/x-www-form-urlencoded\r\nContent-Length: 7\r\n\r\nfoo=bar"))

	got := Mutate(rq, []Mutation{SingleQuotes}, []Mutable{JsonParameter})

	testutils.AssertLen(t, got, 0)
}

func TestApplySingleQuotesMutationToJsonParameter(t *testing.T) {
	rq := http.Parse([]byte("POST /auth HTTP/1.1\r\nContent-Type: application/json\r\nContent-Length: 13\r\n\r\n{\"foo\":\"bar\"}"))

	got := Mutate(rq, []Mutation{SingleQuotes}, []Mutable{JsonParameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertByteEquals(t, got[0].Body, []byte("{\"foo\":\"bar'\"}"))
}

func TestApplySingleQuotesMutationToNumericJsonParameter(t *testing.T) {
	rq := http.Parse([]byte("POST /auth HTTP/1.1\r\nContent-Type: application/json\r\nContent-Length: 13\r\n\r\n{\"foo\": 3}"))

	got := Mutate(rq, []Mutation{SingleQuotes}, []Mutable{JsonParameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertByteEquals(t, got[0].Body, []byte("{\"foo\":\"3'\"}"))
}

func TestApplySingleQuotesMutationToANestedJsonParameter(t *testing.T) {
	rq := http.Parse([]byte("POST /auth HTTP/1.1\r\nContent-Type: application/json\r\nContent-Length: 13\r\n\r\n{\"foo\":{\"bar\":\"baz\"}}"))

	got := Mutate(rq, []Mutation{SingleQuotes}, []Mutable{JsonParameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertByteEquals(t, got[0].Body, []byte(`{"foo":{"bar":"baz'"}}`))
}
