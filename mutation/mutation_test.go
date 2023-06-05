package mutation

import (
	"github.com/kamil-s-solecki/haze/http"
	"github.com/kamil-s-solecki/haze/mutable"
	"github.com/kamil-s-solecki/haze/testutils"
	"testing"
)

func TestEmpty(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{}, []mutable.Mutable{})

	testutils.AssertEmpty(t, got)
}

func TestApplySingleQuotesMutationToPath(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{SingleQuotes}, []mutable.Mutable{mutable.Path})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Path, "/somepath'")
	testutils.AssertEquals(t, got[0].RequestUri, "/somepath'")
}

func TestApplyDoubleQuotesMutationToPath(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{DoubleQuotes}, []mutable.Mutable{mutable.Path})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Path, "/somepath%22")
	testutils.AssertEquals(t, got[0].RequestUri, "/somepath%22")
}

func TestApplySingleQuotesMutationToParameter(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath?foo=bar HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))
	got := Mutate(rq, []Mutation{SingleQuotes}, []mutable.Mutable{mutable.Parameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Query, "foo=bar'")
	testutils.AssertEquals(t, got[0].RequestUri, "/somepath?foo=bar'")
}

func TestApplyDoubleQuotesMutationToParameter(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath?foo=bar HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{DoubleQuotes}, []mutable.Mutable{mutable.Parameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Query, "foo=bar%22")
	testutils.AssertEquals(t, got[0].RequestUri, "/somepath?foo=bar%22")
}

func TestApplyDoubleQuotesMutationToBothParameters(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath?foo=bar&baz=quix HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{DoubleQuotes}, []mutable.Mutable{mutable.Parameter})

	testutils.AssertLen(t, got, 2)
	testutils.AssertEquals(t, got[0].Query, "foo=bar%22&baz=quix")
	testutils.AssertEquals(t, got[0].RequestUri, "/somepath?foo=bar%22&baz=quix")
	testutils.AssertEquals(t, got[1].Query, "foo=bar&baz=quix%22")
	testutils.AssertEquals(t, got[1].RequestUri, "/somepath?foo=bar&baz=quix%22")
}

func TestDoNothingForEmptyQuery(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{DoubleQuotes}, []mutable.Mutable{mutable.Parameter})

	testutils.AssertLen(t, got, 0)
}

func TestDoNothingForEmptyBody(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{DoubleQuotes}, []mutable.Mutable{mutable.BodyParameter})

	testutils.AssertLen(t, got, 0)
}

func TestDoNothingForNonFormUrlencodedBody(t *testing.T) {
	rq := http.Parse([]byte("POST /auth HTTP/1.1\r\nContent-Type: application/json\r\nContent-Length: 13\r\n\r\n\"bar\""))

	got := Mutate(rq, []Mutation{DoubleQuotes}, []mutable.Mutable{mutable.BodyParameter})

	testutils.AssertLen(t, got, 0)
}

func TestApplyDoubleQuotesMutationToBodyParameter(t *testing.T) {
	rq := http.Parse([]byte("POST /auth HTTP/1.1\r\nContent-Type: application/x-www-form-urlencoded\r\nContent-Length: 7\r\n\r\nfoo=bar"))

	got := Mutate(rq, []Mutation{DoubleQuotes}, []mutable.Mutable{mutable.BodyParameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertByteEquals(t, got[0].Body, []byte("foo=bar%22"))
}

func TestApplySingleQuotesMutationToBodyParameter(t *testing.T) {
	rq := http.Parse([]byte("POST /auth HTTP/1.1\r\nContent-Type: application/x-www-form-urlencoded\r\nContent-Length: 7\r\n\r\nfoo=bar"))

	got := Mutate(rq, []Mutation{SingleQuotes}, []mutable.Mutable{mutable.BodyParameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertByteEquals(t, got[0].Body, []byte("foo=bar'"))
}

func TestApplySstiFuzzMutationToParameter(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath?foo=bar HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{SstiFuzz}, []mutable.Mutable{mutable.Parameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Query, "foo=bar${{<%25[%25'%22}}%25%5c.")
	testutils.AssertEquals(t, got[0].RequestUri, "/somepath?foo=bar${{<%25[%25'%22}}%25%5c.")
}

func TestApplyDoubleQuotesMutationToHeader(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath HTTP/1.1\r\nFoo: bar\r\n\r\n"))

	got := Mutate(rq, []Mutation{DoubleQuotes}, []mutable.Mutable{mutable.Header})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Headers["Foo"], "bar\"")
}

func TestApplySstiFuzzMutationToHeader(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath?foo=bar HTTP/1.1\r\nFoo:bar\r\n\r\n"))

	got := Mutate(rq, []Mutation{SstiFuzz}, []mutable.Mutable{mutable.Header})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Headers["Foo"], "bar${{<%[%'\"}}%\\.")
}

func TestApplyDoubleQuotesMutationToCookie(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath HTTP/1.1\r\nCookie: foo=bar\r\n\r\n"))

	got := Mutate(rq, []Mutation{DoubleQuotes}, []mutable.Mutable{mutable.Cookie})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Cookies["foo"], "bar%22")
}

func TestSkipCertainHeaders(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath HTTP/1.1\r\nHost: bar\r\nConnection: bar\r\nContent-Type: bar\r\nContent-Length: bar\r\nAccept-Encoding: bar\r\n\r\n"))

	got := Mutate(rq, []Mutation{DoubleQuotes}, []mutable.Mutable{mutable.Header})

	testutils.AssertLen(t, got, 0)
}

func TestApplyNegativeMutationToParameter(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath?foo=123 HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{Negative}, []mutable.Mutable{mutable.Parameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Query, "foo=-123")
	testutils.AssertEquals(t, got[0].RequestUri, "/somepath?foo=-123")
}

func TestApplyNegativeMutationToPath(t *testing.T) {
	rq := http.Parse([]byte("GET /123 HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{Negative}, []mutable.Mutable{mutable.Path})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Path, "/-123")
	testutils.AssertEquals(t, got[0].RequestUri, "/-123")
}

func TestApplyMinusOneMutationToParameter(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath?foo=123 HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{MinusOne}, []mutable.Mutable{mutable.Parameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Query, "foo=123-1")
	testutils.AssertEquals(t, got[0].RequestUri, "/somepath?foo=123-1")
}

func TestApplyTimesSevenMutationToParameter(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath?foo=123 HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{TimesSeven}, []mutable.Mutable{mutable.Parameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Query, "foo=123*7")
	testutils.AssertEquals(t, got[0].RequestUri, "/somepath?foo=123*7")
}

func TestDoNothingWithNonJsonBody(t *testing.T) {
	rq := http.Parse([]byte("POST /auth HTTP/1.1\r\nContent-Type: application/x-www-form-urlencoded\r\nContent-Length: 7\r\n\r\nfoo=bar"))

	got := Mutate(rq, []Mutation{SingleQuotes}, []mutable.Mutable{mutable.JsonParameter})

	testutils.AssertLen(t, got, 0)
}

func TestApplySingleQuotesMutationToJsonParameter(t *testing.T) {
	rq := http.Parse([]byte("POST /auth HTTP/1.1\r\nContent-Type: application/json\r\nContent-Length: 13\r\n\r\n{\"foo\":\"bar\"}"))

	got := Mutate(rq, []Mutation{SingleQuotes}, []mutable.Mutable{mutable.JsonParameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertByteEquals(t, got[0].Body, []byte("{\"foo\":\"bar'\"}"))
}

func TestApplySingleQuotesMutationToNumericJsonParameter(t *testing.T) {
	rq := http.Parse([]byte("POST /auth HTTP/1.1\r\nContent-Type: application/json\r\nContent-Length: 13\r\n\r\n{\"foo\": 3}"))

	got := Mutate(rq, []Mutation{SingleQuotes}, []mutable.Mutable{mutable.JsonParameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertByteEquals(t, got[0].Body, []byte("{\"foo\":\"3'\"}"))
}

func TestApplySingleQuotesMutationToANestedJsonParameter(t *testing.T) {
	rq := http.Parse([]byte("POST /auth HTTP/1.1\r\nContent-Type: application/json\r\nContent-Length: 13\r\n\r\n{\"foo\":{\"bar\":\"baz\"}}"))

	got := Mutate(rq, []Mutation{SingleQuotes}, []mutable.Mutable{mutable.JsonParameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertByteEquals(t, got[0].Body, []byte(`{"foo":{"bar":"baz'"}}`))
}

func TestApplySingleQuotesMutationToANestedArrayJsonParameter(t *testing.T) {
	rq := http.Parse([]byte("POST /auth HTTP/1.1\r\nContent-Type: application/json\r\nContent-Length: 13\r\n\r\n{\"foo\":[\"bar\"]}"))

	got := Mutate(rq, []Mutation{SingleQuotes}, []mutable.Mutable{mutable.JsonParameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertByteEquals(t, got[0].Body, []byte(`{"foo":["bar'"]}`))
}

func TestApplySingleQuotesMutationToAllValuesInANestedArrayJsonParameter(t *testing.T) {
	rq := http.Parse([]byte("POST /auth HTTP/1.1\r\nContent-Type: application/json\r\nContent-Length: 13\r\n\r\n{\"foo\":[\"bar\",\"baz\"]}"))

	got := Mutate(rq, []Mutation{SingleQuotes}, []mutable.Mutable{mutable.JsonParameter})

	testutils.AssertLen(t, got, 2)
	testutils.AssertByteEquals(t, got[0].Body, []byte(`{"foo":["bar'","baz"]}`))
	testutils.AssertByteEquals(t, got[1].Body, []byte(`{"foo":["bar","baz'"]}`))
}

func TestApplySingleQuotesMutationToANestedJsonInNestedArrayJsonParameter(t *testing.T) {
	rq := http.Parse([]byte("POST /auth HTTP/1.1\r\nContent-Type: application/json\r\nContent-Length: 13\r\n\r\n{\"foo\":[{\"bar\":\"baz\"}]}"))

	got := Mutate(rq, []Mutation{SingleQuotes}, []mutable.Mutable{mutable.JsonParameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertByteEquals(t, got[0].Body, []byte(`{"foo":[{"bar":"baz'"}]}`))
}

func TestApplySingleQuotesMutationToANestedArrayJsonNumericParameter(t *testing.T) {
	rq := http.Parse([]byte("POST /auth HTTP/1.1\r\nContent-Type: application/json\r\nContent-Length: 13\r\n\r\n{\"foo\":[123]}"))

	got := Mutate(rq, []Mutation{SingleQuotes}, []mutable.Mutable{mutable.JsonParameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertByteEquals(t, got[0].Body, []byte(`{"foo":["123'"]}`))
}

func TestApplySingleQuotesMutationToAnArrayJson(t *testing.T) {
	rq := http.Parse([]byte("POST /auth HTTP/1.1\r\nContent-Type: application/json\r\nContent-Length: 13\r\n\r\n[\"bar\"]"))

	got := Mutate(rq, []Mutation{SingleQuotes}, []mutable.Mutable{mutable.JsonParameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertByteEquals(t, got[0].Body, []byte(`["bar'"]`))
}

func TestApplySingleQuotesMutationToAStringJson(t *testing.T) {
	rq := http.Parse([]byte("POST /auth HTTP/1.1\r\nContent-Type: application/json\r\nContent-Length: 13\r\n\r\n\"bar\""))

	got := Mutate(rq, []Mutation{SingleQuotes}, []mutable.Mutable{mutable.JsonParameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertByteEquals(t, got[0].Body, []byte(`"bar'"`))
}

func TestApplyBracketsMutationToHeader(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath HTTP/1.1\r\nFoo: bar\r\n\r\n"))

	got := Mutate(rq, []Mutation{Brackets}, []mutable.Mutable{mutable.Header})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Headers["Foo"], "bar)]}>")
}

func TestApplyBacktickMutationToHeader(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath HTTP/1.1\r\nFoo: bar\r\n\r\n"))

	got := Mutate(rq, []Mutation{Backtick}, []mutable.Mutable{mutable.Header})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Headers["Foo"], "bar`")
}

func TestApplyCommaMutationToHeader(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath HTTP/1.1\r\nFoo: bar\r\n\r\n"))

	got := Mutate(rq, []Mutation{Comma}, []mutable.Mutable{mutable.Header})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Headers["Foo"], "bar,")
}

func TestApplyArraizeToQueryParameter(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath?foo=bar HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{Arraize}, []mutable.Mutable{mutable.ParameterName})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Query, "foo[]=bar")
	testutils.AssertEquals(t, got[0].RequestUri, "/somepath?foo[]=bar")
}

func TestShouldNotArraizeQueryParamVal(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath?foo=bar HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))

	got := Mutate(rq, []Mutation{Arraize}, []mutable.Mutable{mutable.Parameter})

	testutils.AssertLen(t, got, 0)
}

func TestApplyArraizeToBodyParameterName(t *testing.T) {
	rq := http.Parse([]byte("POST /auth HTTP/1.1\r\nContent-Type: application/x-www-form-urlencoded\r\nContent-Length: 7\r\n\r\nfoo=bar"))

	got := Mutate(rq, []Mutation{Arraize}, []mutable.Mutable{mutable.BodyParameterName})

	testutils.AssertLen(t, got, 1)
	testutils.AssertByteEquals(t, got[0].Body, []byte("foo[]=bar"))
}

func TestApplySingleQuotesMutationToMultipartFormParameter(t *testing.T) {
	head := []byte("POST /multi HTTP/1.1\r\nContent-Type: multipart/form-data; boundary=----WebKitFormBoundaryQdBweljBPtRAAu9f\r\nContent-Length: 144\r\n\r\n")
	body := []byte("------WebKitFormBoundaryQdBweljBPtRAAu9f\r\nContent-Disposition: form-data; name=\"foo\"\r\n\r\nbar\r\n------WebKitFormBoundaryQdBweljBPtRAAu9f--\r\n")
	rq := http.Parse(append(head, body...))

	got := Mutate(rq, []Mutation{SingleQuotes}, []mutable.Mutable{mutable.MultipartFormParameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertByteEquals(t, got[0].Body, []byte("------WebKitFormBoundaryQdBweljBPtRAAu9f\r\nContent-Disposition: form-data; name=\"foo\"\r\n\r\nbar'\r\n------WebKitFormBoundaryQdBweljBPtRAAu9f--\r\n"))
}

func TestApplySingleQuotesMutationToAllMultipartFormParameters(t *testing.T) {
	head := []byte("POST /multi HTTP/1.1\r\nContent-Type: multipart/form-data; boundary=----WebKitFormBoundaryQdBweljBPtRAAu9f\r\nContent-Length: 144\r\n\r\n")
	body := []byte("------WebKitFormBoundaryQdBweljBPtRAAu9f\r\nContent-Disposition: form-data; name=\"foo\"\r\n\r\nbar\r\n------WebKitFormBoundaryQdBweljBPtRAAu9f\r\nContent-Disposition: form-data; name=\"baz\"\r\n\r\nquix\r\n------WebKitFormBoundaryQdBweljBPtRAAu9f--\r\n")
	rq := http.Parse(append(head, body...))

	got := Mutate(rq, []Mutation{SingleQuotes}, []mutable.Mutable{mutable.MultipartFormParameter})

	testutils.AssertLen(t, got, 2)
	testutils.AssertByteEquals(t, got[0].Body, []byte("------WebKitFormBoundaryQdBweljBPtRAAu9f\r\nContent-Disposition: form-data; name=\"foo\"\r\n\r\nbar'\r\n------WebKitFormBoundaryQdBweljBPtRAAu9f\r\nContent-Disposition: form-data; name=\"baz\"\r\n\r\nquix\r\n------WebKitFormBoundaryQdBweljBPtRAAu9f--\r\n"))
	testutils.AssertByteEquals(t, got[1].Body, []byte("------WebKitFormBoundaryQdBweljBPtRAAu9f\r\nContent-Disposition: form-data; name=\"foo\"\r\n\r\nbar\r\n------WebKitFormBoundaryQdBweljBPtRAAu9f\r\nContent-Disposition: form-data; name=\"baz\"\r\n\r\nquix'\r\n------WebKitFormBoundaryQdBweljBPtRAAu9f--\r\n"))
}

func TestDoNothingForNonMultipartBody(t *testing.T) {
	rq := http.Parse([]byte("POST /auth HTTP/1.1\r\nContent-Type: application/json\r\nContent-Length: 13\r\n\r\n\"bar\""))

	got := Mutate(rq, []Mutation{DoubleQuotes}, []mutable.Mutable{mutable.MultipartFormParameter})

	testutils.AssertLen(t, got, 0)
}

func TestMultiplyValue(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath?foo=b HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))
	got := Mutate(rq, []Mutation{TwentyTimes}, []mutable.Mutable{mutable.Parameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Query, "foo=bbbbbbbbbbbbbbbbbbbb")
	testutils.AssertEquals(t, got[0].RequestUri, "/somepath?foo=bbbbbbbbbbbbbbbbbbbb")
}

func TestNullbyte(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath?foo=bar HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))
	got := Mutate(rq, []Mutation{Nullbyte}, []mutable.Mutable{mutable.Parameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Query, "foo=%00bar")
	testutils.AssertEquals(t, got[0].RequestUri, "/somepath?foo=%00bar")
}

func TestNotApplyNullbyteToHeaders(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath?foo=bar HTTP/1.1\r\nHost:www.example.com\r\nUser-Agent:foo\r\n\r\n"))
	got := Mutate(rq, []Mutation{Nullbyte}, []mutable.Mutable{mutable.Header})

	testutils.AssertLen(t, got, 0)
}

func TestDotDotSlash(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath?foo=bar HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))
	got := Mutate(rq, []Mutation{DotDotSlash}, []mutable.Mutable{mutable.Parameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Query, "foo=bar/../../idontexist.txt")
}

func TestXmlEscape(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath?foo=bar HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))
	got := Mutate(rq, []Mutation{XmlEscape}, []mutable.Mutable{mutable.Parameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Query, "foo=bar%22><foons:Foo%20%22")
}

func TestWhitespaces(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath?foo=bar HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))
	got := Mutate(rq, []Mutation{Whitespaces}, []mutable.Mutable{mutable.Parameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Query, "foo=%20%09%0c%0d%0abar")
}

func TestNotApplyWhitespacesToHeaders(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath?foo=bar HTTP/1.1\r\nHost:www.example.com\r\nUser-Agent:foo\r\n\r\n"))
	got := Mutate(rq, []Mutation{Whitespaces}, []mutable.Mutable{mutable.Header})

	testutils.AssertLen(t, got, 0)
}

func TestSemicolonCsv(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath?foo=bar HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))
	got := Mutate(rq, []Mutation{SemicolonCsv}, []mutable.Mutable{mutable.Parameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Query, "foo=bar%3bbar")
}

func TestColon(t *testing.T) {
	rq := http.Parse([]byte("GET /somepath?foo=bar HTTP/1.1\r\nHost:www.example.com\r\n\r\n"))
	got := Mutate(rq, []Mutation{Colon}, []mutable.Mutable{mutable.Parameter})

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Query, "foo=bar:bar")
}
