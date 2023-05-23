package mutation

import "github.com/kamil-s-solecki/haze/http"

type Mutation func(http.Request, Mutable) []http.Request

func SingleQuotes(rq http.Request, mutable Mutable) []http.Request {
	trans := func(val string) string {
		return val + "'"
	}
	return mutable(rq, trans)
}

func DoubleQuotes(rq http.Request, mutable Mutable) []http.Request {
	trans := func(val string) string {
		return val + "\""
	}
	return mutable(rq, trans)
}

func SstiFuzz(rq http.Request, mutable Mutable) []http.Request {
	trans := func(val string) string {
		return val + "${{<%[%'\"}}%\\."
	}
	return mutable(rq, trans)
}

func Negative(rq http.Request, mutable Mutable) []http.Request {
	trans := func(val string) string {
		return "-" + val
	}
	return mutable(rq, trans)
}

func MinusOne(rq http.Request, mutable Mutable) []http.Request {
	trans := func(val string) string {
		return val + "-1"
	}
	return mutable(rq, trans)
}

func TimesSeven(rq http.Request, mutable Mutable) []http.Request {
	trans := func(val string) string {
		return val + "*7"
	}
	return mutable(rq, trans)
}

func AllMutations() []Mutation {
	return []Mutation{SingleQuotes, DoubleQuotes, SstiFuzz, Negative, MinusOne, TimesSeven}
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
