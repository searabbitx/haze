package mutation

import (
	"github.com/kamil-s-solecki/haze/http"
	"github.com/kamil-s-solecki/haze/mutable"
	"strings"
)

type Mutation struct {
	name  string
	apply func(http.Request, mutable.Mutable) []http.Request
}

var SingleQuotes = Mutation{"SingleQuotes", singleQuotes}

func singleQuotes(rq http.Request, mutable mutable.Mutable) []http.Request {
	return suffixMutation(rq, mutable, "'")
}

var DoubleQuotes = Mutation{"DoubleQuotes", doubleQuotes}

func doubleQuotes(rq http.Request, mutable mutable.Mutable) []http.Request {
	return suffixMutation(rq, mutable, "\"")
}

var SstiFuzz = Mutation{"SstiFuzz", sstiFuzz}

func sstiFuzz(rq http.Request, mutable mutable.Mutable) []http.Request {
	return suffixMutation(rq, mutable, "${{<%[%'\"}}%\\.")
}

var Negative = Mutation{"Negative", negative}

func negative(rq http.Request, mutable mutable.Mutable) []http.Request {
	return prefixMutation(rq, mutable, "-")
}

var MinusOne = Mutation{"MinusOne", minusOne}

func minusOne(rq http.Request, mutable mutable.Mutable) []http.Request {
	return suffixMutation(rq, mutable, "-1")
}

var TimesSeven = Mutation{"TimesSeven", timesSeven}

func timesSeven(rq http.Request, mutable mutable.Mutable) []http.Request {
	return suffixMutation(rq, mutable, "*7")
}

var Brackets = Mutation{"Brackets", brackets}

func brackets(rq http.Request, mutable mutable.Mutable) []http.Request {
	return suffixMutation(rq, mutable, ")]}>")
}

var Backtick = Mutation{"Backtick", backtick}

func backtick(rq http.Request, mutable mutable.Mutable) []http.Request {
	return suffixMutation(rq, mutable, "`")
}

var Comma = Mutation{"Comma", comma}

func comma(rq http.Request, mutable mutable.Mutable) []http.Request {
	return suffixMutation(rq, mutable, ",")
}

var Arraize = Mutation{"Arraize", arraize}

func arraize(rq http.Request, mutable mutable.Mutable) []http.Request {
	return suffixMutation(rq, mutable, "[]")
}

var TwentyTimes = Mutation{"TwentyTimes", twentyTimes}

func twentyTimes(rq http.Request, mutable mutable.Mutable) []http.Request {
	trans := func(val string) string {
		return strings.Repeat(val, 20)
	}
	return mutable.Apply(rq, trans)
}

var Nullbyte = Mutation{"Nullbyte", nullbyte}

func nullbyte(rq http.Request, mutable mutable.Mutable) []http.Request {
	return prefixMutation(rq, mutable, "\x00")
}

var DotDotSlash = Mutation{"DotDotSlash", dotDotSlash}

func dotDotSlash(rq http.Request, mutable mutable.Mutable) []http.Request {
	return suffixMutation(rq, mutable, "/../../idontexist.txt")
}

var XmlEscape = Mutation{"XmlEscape", xmlEscape}

func xmlEscape(rq http.Request, mutable mutable.Mutable) []http.Request {
	return suffixMutation(rq, mutable, `"><foons:Foo "`)
}

var Whitespaces = Mutation{"Whitespaces", whitespaces}

func whitespaces(rq http.Request, mutable mutable.Mutable) []http.Request {
	return prefixMutation(rq, mutable, " \t\f\r\n")
}

var SemicolonCsv = Mutation{"SemicolonCsv", semicolonCsv}

func semicolonCsv(rq http.Request, mutable mutable.Mutable) []http.Request {
	trans := func(val string) string {
		return val + ";" + val
	}
	return mutable.Apply(rq, trans)
}

var Colon = Mutation{"Colon", colon}

func colon(rq http.Request, mutable mutable.Mutable) []http.Request {
	trans := func(val string) string {
		return val + ":" + val
	}
	return mutable.Apply(rq, trans)
}

var NeNosqli = Mutation{"NeNosqli", neNosqli}

func neNosqli(rq http.Request, mutable mutable.Mutable) []http.Request {
	return suffixMutation(rq, mutable, "[$ne]")
}

func suffixMutation(rq http.Request, mutable mutable.Mutable, suffix string) []http.Request {
	trans := func(val string) string {
		return val + suffix
	}
	return mutable.Apply(rq, trans)
}

func prefixMutation(rq http.Request, mutable mutable.Mutable, prefix string) []http.Request {
	trans := func(val string) string {
		return prefix + val
	}
	return mutable.Apply(rq, trans)
}

func canApply(mutation Mutation, mtbl mutable.Mutable) bool {
	switch mutation.name {
	case Arraize.name, NeNosqli.name:
		switch mtbl.Name {
		case mutable.ParameterName.Name, mutable.BodyParameterName.Name:
			return true
		default:
			return false
		}
	case Nullbyte.name:
		switch mtbl.Name {
		case mutable.Header.Name:
			return false
		default:
			return true
		}
	case Whitespaces.name:
		switch mtbl.Name {
		case mutable.Header.Name:
			return false
		default:
			return true
		}
	default:
		return true
	}
}

func Mutate(rq http.Request, mutations []Mutation, mutables []mutable.Mutable) []http.Request {
	result := []http.Request{}
	for _, mutation := range mutations {
		for _, mutable := range mutables {
			if !canApply(mutation, mutable) {
				continue
			}
			mrq := mutation.apply(rq, mutable)
			result = append(result, mrq...)
		}
	}
	return result
}

func AllMutations() []Mutation {
	return []Mutation{SingleQuotes, DoubleQuotes, SstiFuzz, Negative, MinusOne,
		TimesSeven, Brackets, Backtick, Comma, Arraize, TwentyTimes, Nullbyte,
		DotDotSlash, XmlEscape, Whitespaces, SemicolonCsv, Colon, NeNosqli}
}
