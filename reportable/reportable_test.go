package reportable

import (
	"github.com/kamil-s-solecki/haze/http"
	"github.com/kamil-s-solecki/haze/testutils"
	"testing"
)

func TestShouldNotReport200(t *testing.T) {
	res := http.Response{Code: 200}

	got := IsReportable(res, []Matcher{MatchCodes("500")})

	testutils.AssertFalse(t, got)
}

func TestShouldReportCodes(t *testing.T) {
	cases := []struct {
		icode int
		scode string
	}{
		{500, "500"},
		{510, "500,510"},
		{510, "500,505-520"},
	}

	for _, c := range cases {
		res := http.Response{Code: c.icode}

		got := IsReportable(res, []Matcher{MatchCodes(c.scode)})

		testutils.AssertTrue(t, got)
	}
}

func TestShouldReportLengths(t *testing.T) {
	cases := []struct {
		ilen int64
		slen string
	}{
		{250, "100-300"},
	}

	for _, c := range cases {
		res := http.Response{Length: c.ilen}

		got := IsReportable(res, []Matcher{MatchLengths(c.slen)})

		testutils.AssertTrue(t, got)
	}
}
