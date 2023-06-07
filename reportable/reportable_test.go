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

func TestShouldConstructFromArgsWithCodesOnly(t *testing.T) {
	args := cliargs.Args{MatchCodes: "500,501-502"}

	ms, fs := FromArgs(args)

	testutils.AssertLen(t, ms, 1)
	testutils.AssertLen(t, fs, 0)
	testutils.AssertTrue(t, IsReportable(http.Response{Code: 500}, ms, fs))
}

func TestShouldConstructFromArgsWithCodesAndLens(t *testing.T) {
	args := cliargs.Args{MatchCodes: "500,501-502", MatchLengths: "100-200"}

	ms, fs := FromArgs(args)

	testutils.AssertLen(t, ms, 2)
	testutils.AssertLen(t, fs, 0)
	testutils.AssertTrue(t, IsReportable(http.Response{Code: 500}, ms, fs))
	testutils.AssertTrue(t, IsReportable(http.Response{Code: 200, Length: 150}, ms, fs))
}

func TestShouldConstructFromArgsWithFilters(t *testing.T) {
	args := cliargs.Args{MatchCodes: "500", FilterCodes: "510", FilterLengths: "100-200"}

	ms, fs := FromArgs(args)

	testutils.AssertLen(t, ms, 1)
	testutils.AssertLen(t, fs, 2)
	testutils.AssertTrue(t, IsReportable(http.Response{Code: 500}, ms, fs))
	testutils.AssertFalse(t, IsReportable(http.Response{Code: 510}, ms, fs))
	testutils.AssertFalse(t, IsReportable(http.Response{Code: 500, Length: 150}, ms, fs))
}

func TestShouldReportWhenStringMatches(t *testing.T) {
	res := http.Response{Raw: []byte("foo bar baz")}

	got := IsReportable(res, []Matcher{MatchString("bar")}, []Filter{})

	testutils.AssertTrue(t, got)
}

func TestShouldNotReportWhenStringDoesNotMatch(t *testing.T) {
	res := http.Response{Raw: []byte("foo bad baz")}

	got := IsReportable(res, []Matcher{MatchString("bar")}, []Filter{})

	testutils.AssertFalse(t, got)
}

func TestShouldNotReportWhenStringFilters(t *testing.T) {
	res := http.Response{Code: 500, Raw: []byte("foo bar baz")}

	got := IsReportable(res, []Matcher{MatchCodes("500")}, []Filter{FilterString("bar")})

	testutils.AssertFalse(t, got)
}

func TestShouldReportWhenStringDoesNotFilter(t *testing.T) {
	res := http.Response{Code: 400, Raw: []byte("foo ban baz")}

	got := IsReportable(res, []Matcher{MatchCodes("400-599")}, []Filter{FilterString("bar"), FilterCodes("404")})

	testutils.AssertTrue(t, got)
}
