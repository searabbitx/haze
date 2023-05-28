package reportable

import (
	"github.com/kamil-s-solecki/haze/cliargs"
	"github.com/kamil-s-solecki/haze/http"
	"github.com/kamil-s-solecki/haze/testutils"
	"testing"
)

func TestShouldNotReport200(t *testing.T) {
	res := http.Response{Code: 200}

	got := IsReportable(res, []Matcher{MatchCodes("500")}, []Filter{})

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

		got := IsReportable(res, []Matcher{MatchCodes(c.scode)}, []Filter{})

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

		got := IsReportable(res, []Matcher{MatchLengths(c.slen)}, []Filter{})

		testutils.AssertTrue(t, got)
	}
}

func TestShouldConstructFromArgsWithCodesOnly(t *testing.T) {
	args := cliargs.Args{MatchCodes: "500,501-502"}

	got := FromArgs(args)

	testutils.AssertLen(t, got, 1)
	testutils.AssertTrue(t, IsReportable(http.Response{Code: 500}, got, []Filter{}))
}

func TestShouldConstructFromArgsWithCodesAndLens(t *testing.T) {
	args := cliargs.Args{MatchCodes: "500,501-502", MatchLengths: "100-200"}

	got := FromArgs(args)

	testutils.AssertLen(t, got, 2)
	testutils.AssertTrue(t, IsReportable(http.Response{Code: 500}, got, []Filter{}))
	testutils.AssertTrue(t, IsReportable(http.Response{Code: 200, Length: 150}, got, []Filter{}))
}

func TestShouldReport500When200IsFiltered(t *testing.T) {
	res := http.Response{Code: 500}

	got := IsReportable(res, []Matcher{}, []Filter{FilterCodes("200")})

	testutils.AssertTrue(t, got)
}

func TestShouldNotReport500When500IsFiltered(t *testing.T) {
	res := http.Response{Code: 500}

	got := IsReportable(res, []Matcher{}, []Filter{FilterCodes("500")})

	testutils.AssertFalse(t, got)
}

func TestShouldNotReportMatched500When500IsFiltered(t *testing.T) {
	res := http.Response{Code: 500}

	got := IsReportable(res, []Matcher{MatchCodes("500")}, []Filter{FilterCodes("500")})

	testutils.AssertFalse(t, got)
}

func TestShouldNotReportMatched500WhenLenIsFiltered(t *testing.T) {
	res := http.Response{Code: 500, Length: 1500}

	got := IsReportable(res, []Matcher{MatchCodes("500")}, []Filter{FilterLengths("1500")})

	testutils.AssertFalse(t, got)
}

func TestShouldNotReportWhenOneFilterFiltered(t *testing.T) {
	res := http.Response{Code: 500, Length: 2000}

	got := IsReportable(res, []Matcher{MatchCodes("500-599")}, []Filter{FilterLengths("1500"), FilterCodes("500")})

	testutils.AssertFalse(t, got)
}
