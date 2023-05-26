package mutation

import "github.com/kamil-s-solecki/haze/http"

type Mutation func(http.Request, Mutable) []http.Request

func SingleQuotes(rq http.Request, mutable Mutable) []http.Request {
	return suffixMutation(rq, mutable, "'")
}

func DoubleQuotes(rq http.Request, mutable Mutable) []http.Request {
	return suffixMutation(rq, mutable, "\"")
}

func SstiFuzz(rq http.Request, mutable Mutable) []http.Request {
	return suffixMutation(rq, mutable, "${{<%[%'\"}}%\\.")
}

func Negative(rq http.Request, mutable Mutable) []http.Request {
	return prefixMutation(rq, mutable, "-")
}

func MinusOne(rq http.Request, mutable Mutable) []http.Request {
	return suffixMutation(rq, mutable, "-1")
}

func TimesSeven(rq http.Request, mutable Mutable) []http.Request {
	return suffixMutation(rq, mutable, "*7")
}

func Brackets(rq http.Request, mutable Mutable) []http.Request {
	return suffixMutation(rq, mutable, ")]}>")
}

func Backtick(rq http.Request, mutable Mutable) []http.Request {
	return suffixMutation(rq, mutable, "`")
}

func Comma(rq http.Request, mutable Mutable) []http.Request {
	return suffixMutation(rq, mutable, ",")
}

func suffixMutation(rq http.Request, mutable Mutable, suffix string) []http.Request {
	trans := func(val string) string {
		return val + suffix
	}
	return mutable(rq, trans)
}

func prefixMutation(rq http.Request, mutable Mutable, prefix string) []http.Request {
	trans := func(val string) string {
		return prefix + val
	}
	return mutable(rq, trans)
}

func AllMutations() []Mutation {
	return []Mutation{SingleQuotes, DoubleQuotes, SstiFuzz, Negative, MinusOne,
		TimesSeven, Brackets, Backtick, Comma}
}

func Mutate(rq http.Request, mutations []Mutation, mutables []Mutable) []http.Request {
	result := []http.Request{}
	for _, mutation := range mutations {
		for _, mutable := range mutables {
			mrq := mutation(rq, mutable)
			result = append(result, mrq...)
		}
	}
	return result
}
