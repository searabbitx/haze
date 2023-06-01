package reportable

import (
	"bytes"
	"github.com/kamil-s-solecki/haze/cliargs"
	"github.com/kamil-s-solecki/haze/http"
	"strconv"
	"strings"
)

type Matcher func(http.Response) bool

type Filter func(http.Response) bool

type Range struct{ From, To int }

func MatchCodes(codes string) Matcher {
	ranges := parseRanges(codes)
	return func(res http.Response) bool {
		return isValueInRanges(ranges, res.Code)
	}
}

func MatchLengths(lens string) Matcher {
	ranges := parseRanges(lens)
	return func(res http.Response) bool {
		return isValueInRanges(ranges, int(res.Length))
	}
}

func MatchString(str string) Matcher {
	return func(res http.Response) bool {
		return bytes.Contains(res.Raw, []byte(str))
	}
}

func FilterCodes(codes string) Filter {
	ranges := parseRanges(codes)
	return func(res http.Response) bool {
		return !isValueInRanges(ranges, res.Code)
	}
}

func FilterLengths(lens string) Filter {
	ranges := parseRanges(lens)
	return func(res http.Response) bool {
		return !isValueInRanges(ranges, int(res.Length))
	}
}

func FilterString(str string) Filter {
	return func(res http.Response) bool {
		return bytes.Contains(res.Raw, []byte(str))
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

func FromArgs(args cliargs.Args) ([]Matcher, []Filter) {
	matchers := []Matcher{}
	if args.MatchLengths != "" {
		matchers = append(matchers, MatchLengths(args.MatchLengths))
	}
	if args.MatchString != "" {
		matchers = append(matchers, MatchString(args.MatchString))
	}
	if !(len(matchers) > 0 && args.MatchCodes == "500-599") {
		matchers = append(matchers, MatchCodes(args.MatchCodes))
	}

	filters := []Filter{}
	if args.FilterCodes != "" {
		filters = append(filters, FilterCodes(args.FilterCodes))
	}
	if args.FilterLengths != "" {
		filters = append(filters, FilterLengths(args.FilterLengths))
	}
	if args.FilterString != "" {
		filters = append(filters, FilterString(args.FilterString))
	}
	return matchers, filters
}

func IsReportable(res http.Response, matchers []Matcher, filters []Filter) bool {
	matched := false
	filtered := true

	for _, matcher := range matchers {
		if matcher(res) {
			matched = true
			break
		}
	}
	for _, filter := range filters {
		if !filter(res) {
			filtered = false
			break
		}
	}
	return filtered && (matched || len(matchers) == 0)
}
