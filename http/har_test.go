package http

import (
	"github.com/kamil-s-solecki/haze/testutils"
	"os"
	"testing"
)

func readHar(path string) []byte {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return content
}

func TestReturnEmptyArrayForNoRequestsInHar(t *testing.T) {
	har := readHar("../var/hars/empty.har")

	got := ParseHar(har)

	testutils.AssertLen(t, got, 0)
}

func TestParseGetRequestFromHar(t *testing.T) {
	har := readHar("../var/hars/get.har")

	got := ParseHar(har)

	testutils.AssertLen(t, got, 1)
	testutils.AssertEquals(t, got[0].Method, "GET")
	testutils.AssertEquals(t, got[0].RequestUri, "/?foo=bar")
	testutils.AssertEquals(t, got[0].Path, "/")
	testutils.AssertEquals(t, got[0].Query, "foo=bar")
	testutils.AssertMapEquals(t, got[0].Cookies, map[string]string{"foocookie": "barcookie"})
}
