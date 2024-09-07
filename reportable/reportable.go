package reportable

import (
	"github.com/kamil-s-solecki/haze/http"
	"strconv"
	"strings"
)

type Matcher func(http.Response) bool

type Range struct{ From, To int }

func MatchCodes(codes string) Matcher {
	ranges := parseRanges(codes)
	return func(res http.Response) bool {
		return isValueInRanges(ranges, res.Code)
	}
}

func MatchLengths(codes string) Matcher {
	ranges := parseRanges(codes)
	return func(res http.Response) bool {
		return isValueInRanges(ranges, int(res.Length))
	}
}

func isValueInRanges(ranges []Range, val int) bool {
	for _, ran := range ranges {
		if val >= ran.From && val <= ran.To {
			return true
		}
	}
	return false
}

func parseRanges(val string) []Range {
	ranges := []Range{}
	for _, ran := range strings.Split(val, ",") {
		ranges = append(ranges, parseRange(ran))
	}
	return ranges
}

func parseRange(val string) Range {
	ran := Range{}
	splitted := strings.Split(val, "-")

	ran.From, _ = strconv.Atoi(splitted[0])
	if len(splitted) == 2 {
		ran.To, _ = strconv.Atoi(splitted[1])
	} else {
		ran.To = ran.From
	}
	return ran
}

func IsReportable(res http.Response, matchers []Matcher) bool {
	for _, matcher := range matchers {
		if matcher(res) {
			return true
		}
	}
	return false
}
