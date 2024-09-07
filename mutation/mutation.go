package mutation

import "github.com/kamil-s-solecki/haze/http"

type Mutation struct {
	name  string
	apply func(http.Request, Mutable) []http.Request
}

func singleQuotes(rq http.Request, mutable Mutable) []http.Request {
	return suffixMutation(rq, mutable, "'")
}

var SingleQuotes = Mutation{"SingleQuotes", singleQuotes}

func doubleQuotes(rq http.Request, mutable Mutable) []http.Request {
	return suffixMutation(rq, mutable, "\"")
}

var DoubleQuotes = Mutation{"DoubleQuotes", doubleQuotes}

func sstiFuzz(rq http.Request, mutable Mutable) []http.Request {
	return suffixMutation(rq, mutable, "${{<%[%'\"}}%\\.")
}

var SstiFuzz = Mutation{"SstiFuzz", sstiFuzz}

func negative(rq http.Request, mutable Mutable) []http.Request {
	return prefixMutation(rq, mutable, "-")
}

var Negative = Mutation{"Negative", negative}

func minusOne(rq http.Request, mutable Mutable) []http.Request {
	return suffixMutation(rq, mutable, "-1")
}

var MinusOne = Mutation{"MinusOne", minusOne}

func timesSeven(rq http.Request, mutable Mutable) []http.Request {
	return suffixMutation(rq, mutable, "*7")
}

var TimesSeven = Mutation{"TimesSeven", timesSeven}

func brackets(rq http.Request, mutable Mutable) []http.Request {
	return suffixMutation(rq, mutable, ")]}>")
}

var Brackets = Mutation{"Brackets", brackets}

func backtick(rq http.Request, mutable Mutable) []http.Request {
	return suffixMutation(rq, mutable, "`")
}

var Backtick = Mutation{"Backtick", backtick}

func comma(rq http.Request, mutable Mutable) []http.Request {
	return suffixMutation(rq, mutable, ",")
}

var Comma = Mutation{"Comma", comma}

var Arraize = Mutation{"Arraize", arraize}

func arraize(rq http.Request, mutable Mutable) []http.Request {
	return suffixMutation(rq, mutable, "[]")
}

func suffixMutation(rq http.Request, mutable Mutable, suffix string) []http.Request {
	trans := func(val string) string {
		return val + suffix
	}
	return mutable.apply(rq, trans)
}

func prefixMutation(rq http.Request, mutable Mutable, prefix string) []http.Request {
	trans := func(val string) string {
		return prefix + val
	}
	return mutable.apply(rq, trans)
}

func cannotApply(mutation Mutation, mutable Mutable) bool {
	if mutation.name == Arraize.name && mutable.name != ParameterName.name {
		return true
	}
	return false
}

func AllMutations() []Mutation {
	return []Mutation{SingleQuotes, DoubleQuotes, SstiFuzz, Negative, MinusOne,
		TimesSeven, Brackets, Backtick, Comma}
}

func Mutate(rq http.Request, mutations []Mutation, mutables []Mutable) []http.Request {
	result := []http.Request{}
	for _, mutation := range mutations {
		for _, mutable := range mutables {
			if cannotApply(mutation, mutable) {
				continue
			}
			mrq := mutation.apply(rq, mutable)
			result = append(result, mrq...)
		}
	}
	return result
}
