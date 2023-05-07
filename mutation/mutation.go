package mutation

import "github.com/kamil-s-solecki/haze/http"

type Mutation func(http.Request) http.Request

func SingleQuotes(rq http.Request) http.Request {
	return rq.WithPath(rq.Path + "'")
}

func DoubleQuotes(rq http.Request) http.Request {
	return rq.WithPath(rq.Path + "\"")
}

func AllMutations() []Mutation {
	return []Mutation{SingleQuotes, DoubleQuotes}
}

func Mutate(rq http.Request, mutations []Mutation) []http.Request {
	result := []http.Request{}
	for _, mut := range mutations {
		result = append(result, mut(rq))
	}
	return result
}
