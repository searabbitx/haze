package mutation

import "github.com/kamil-s-solecki/haze/http"

type Mutation string

const (
	SingleQuotes Mutation = "SingleQuotes"
)

func AllMutations() []Mutation {
	return []Mutation{SingleQuotes}
}

func Mutate(rq http.Request, mutations []Mutation) []http.Request {
	result := []http.Request{}
	for i := 0; i < len(mutations); i++ {
		result = append(result, rq.WithPath(rq.Path+"'"))
	}
	return result
}
